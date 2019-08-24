package app

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/genetic"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

type App struct {
	w    *world.World
	goal world.Object

	maxDist float64

	maxScore float64
	scores   []float64
	gen      int

	top []float64

	ticks           int
	parents         int
	parentsPerChild int
	recombChance    float64
	mutChance       float64
}

func New(nbots int) *App {
	app := &App{
		ticks:           1500,
		recombChance:    0.2,
		mutChance:       0.005,
		parentsPerChild: 3,
		parents:         nbots / 20,
	}

	app.parents = nbots / 5
	if app.parents < 2 {
		app.parents = 2
	} else if app.parents > app.parentsPerChild*6 {
		app.parents = app.parentsPerChild * 6
	}

	height := 800.0
	width := 1200.0

	app.w = world.New(world.Rect(0, 0, width, height))
	app.maxDist = app.w.MaxDistance()

	bots := make([]*world.Bot, nbots)
	for i := 0; i < nbots; i++ {
		bots[i] = world.NewBot(600, 600, 2)
		app.botPos(bots[i])
	}

	for _, bot := range bots {
		app.w.AddBot(bot)
	}
	// ONE app.w.AddObject(world.NewWall(world.Rect(0, 380, 550, 20)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(650, 380, 550, 20)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 380)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 200)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 800)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 620)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 180, 200, 20)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 50, 200, 20)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 50, 20, 130)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(530, 0, 20, 50)))
	// ONE app.goal = world.NewGoal(world.Rect(400, 100, 50, 50))

	// |--   |
	app.w.AddObject(world.NewWall(world.Rect(530, 180, 70, 20)))
	app.w.AddObject(world.NewWall(world.Rect(600, 180, 20, 70)))
	// |   --|
	app.w.AddObject(world.NewWall(world.Rect(600, 300, 70, 20)))
	// |--   |
	app.w.AddObject(world.NewWall(world.Rect(530, 400, 70, 20)))
	// |--  -|
	app.w.AddObject(world.NewWall(world.Rect(530, 500, 70, 20)))
	app.w.AddObject(world.NewWall(world.Rect(630, 500, 40, 20)))

	app.w.AddObject(world.NewWall(world.Rect(400, 100, 20, 80)))

	app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 800)))
	app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 620)))

	app.w.AddObject(world.NewWall(world.Rect(320, 180, 220, 20)))
	app.w.AddObject(world.NewWall(world.Rect(200, 50, 350, 20)))

	app.w.AddObject(world.NewWall(world.Rect(200, 300, 120, 20)))
	app.w.AddObject(world.NewWall(world.Rect(320, 180, 20, 140)))

	app.w.AddObject(world.NewWall(world.Rect(200, 50, 20, 250)))

	app.w.AddObject(world.NewWall(world.Rect(530, 0, 20, 50)))

	app.goal = world.NewGoal(world.Rect(250, 230, 50, 50))

	app.w.AddObject(app.goal)

	app.w.AddObject(world.NewWall(world.Rect(-20, -20, 20, height+40)))
	app.w.AddObject(world.NewWall(world.Rect(-20, -20, width+40, 20)))
	app.w.AddObject(world.NewWall(world.Rect(width, -20, 20, height+40)))
	app.w.AddObject(world.NewWall(world.Rect(-20, height, width+40, 20)))

	return app
}

func (app *App) botPos(bot *world.Bot) {
	//if app.gen > 30 {
	bot.SetPos(float64(570+rand.Intn(60)), float64(650+rand.Intn(100)))
	bot.Reset(rand.Float64())
	return
	//}

	bot.Reset(0)
	bot.SetPos(600, 750)
}

func (app *App) World() *world.World {
	return app.w
}

func (app *App) MaxScore() float64 {
	return app.maxScore
}

func (app *App) MedianScore() float64 {
	if len(app.scores) == 0 {
		return 0
	}
	return app.scores[len(app.scores)/2]
}

func (app *App) Generation() int {
	return app.gen
}

func (app *App) Export() string {
	if app.top == nil {
		return ""
	}

	out := make([]string, len(app.top))
	for i := range app.top {
		out[i] = strconv.FormatFloat(app.top[i], 'g', -1, 64)
	}

	return strings.Join(out, "\n")
}

func (app *App) Import(s string) error {
	lines := strings.Split(s, "\n")
	weights := make([]float64, 0, len(lines))
	for _, l := range lines {
		l = strings.Trim(l, "\n\r ")
		if l == "" {
			continue
		}

		f, err := strconv.ParseFloat(l, 64)
		if err != nil {
			return err
		}
		weights = append(weights, f)
	}

	for _, b := range app.w.Bots {
		if err := b.Brain().SetWeights(weights); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) NewGeneration(tick chan struct{}) {
	sort.Slice(app.w.Bots, func(i, j int) bool {
		return app.w.Bots[i].Fitness() > app.w.Bots[j].Fitness()
	})

	if len(app.scores) < len(app.w.Bots) {
		app.scores = make([]float64, len(app.w.Bots))
	}
	for i := range app.w.Bots {
		app.scores[i] = app.w.Bots[i].Fitness()
	}
	sort.Float64s(app.scores)

	app.maxScore = app.w.Bots[0].Fitness()
	w := app.w.Bots[0].Brain().Weights()
	app.top = make([]float64, len(w))
	copy(app.top, w)
	app.gen++

	tick <- struct{}{}

	e := make([]genetic.Entity, len(app.w.Bots))
	for i := range app.w.Bots {
		e[i] = app.w.Bots[i]
	}

	parents := app.parents
	best := make([][]float64, 0, parents)
	for i := 0; i < app.parents; i++ {
		sel := genetic.Select(e).(*world.Bot)
		w := sel.Brain().Weights()
		cp := make([]float64, len(w))
		copy(cp, w)
		best = append(best, cp)
	}
	parents = len(best)
	perChild := int(math.Ceil(float64(parents) / float64(app.parentsPerChild)))
	if perChild == 0 {
		perChild = 1
	}

	for i := parents; i < len(app.w.Bots); i++ {
		n := i % perChild
		min, max := n*app.parentsPerChild, (n+1)*app.parentsPerChild
		if max > parents {
			max = parents
		}

		app.w.Bots[i].Brain().SetWeights(
			genetic.Reproduce(
				best[min:max],
				app.recombChance,
				app.mutChance,
			),
		)
	}

	for _, b := range app.w.Bots {
		app.botPos(b)
	}
}

func (app *App) Run(sleep time.Duration, n int) (<-chan struct{}, <-chan struct{}, chan<- struct{}) {
	wait := make(chan struct{})
	tick := make(chan struct{})
	stop := make(chan struct{})

	go func() {
		var count int
		var dist float64
		var b *world.Bot
		var i int
		var score float64

		dists := make([]float64, len(app.w.Bots))
		times := n > 0

	outer:
		for {
			select {
			case <-stop:
				break outer
			default:
			}

			for i, b = range app.w.Bots {
				dists[i] = b.Distance(app.goal)
			}
			app.w.Tick()
			count++

			if count%app.ticks == 0 {
				count = 0
				app.NewGeneration(tick)
				if times {
					n--
					if n <= 0 {
						break
					}
				}
			}

			for _, b = range app.w.Bots {
				dist = b.Distance(app.goal)
				score = app.maxDist / (10 * dist)
				if dist < 5 {
					score = 1000
				} else if b.Speed() < 0.2 {
					score -= 10
				}

				b.Reward(score)
			}

			if sleep != 0 {
				time.Sleep(sleep)
			}
		}

		close(tick)
		wait <- struct{}{}
	}()

	return tick, wait, stop
}
