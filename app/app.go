package app

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/genetic"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

type App struct {
	w        *world.World
	goal     world.Object
	maxDist  float64
	maxScore float64
	top      []float64
}

func New() *App {
	app := &App{}

	nbots := 100
	height := 800.0
	width := 1200.0

	app.w = world.New(world.Rect(0, 0, width, height))
	app.maxDist = app.w.MaxDistance()

	bots := make([]*world.Bot, 0, nbots)
	for i := 0; i < nbots; i++ {
		bots = append(bots, world.NewBot(600, 600, 2.5))
	}

	for _, bot := range bots {
		app.w.AddBot(bot)
	}
	// ONE app.w.AddObject(world.NewWall(world.Rect(0, 380, 550, 20)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(650, 380, 550, 20)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 380)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 200)))
	// ONE // app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 800)))
	// ONE // app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 620)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 180, 200, 20)))
	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 50, 200, 20)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(350, 50, 20, 130)))

	// ONE app.w.AddObject(world.NewWall(world.Rect(530, 0, 20, 50)))
	// ONE app.goal = world.NewGoal(world.Rect(400, 100, 50, 50))

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

func (app *App) MaxScore() float64 {
	return app.maxScore
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

func (app *App) NewGeneration(top int, chance float64, tick chan struct{}) {
	sort.Slice(app.w.Bots, func(i, j int) bool {
		return app.w.Bots[i].Score() > app.w.Bots[j].Score()
	})

	app.top = app.w.Bots[0].Brain().Weights()
	app.maxScore = app.w.Bots[0].Score()

	select {
	case tick <- struct{}{}:
	default:
	}

	best := make([][]float64, 0, top)
	for i := 0; i < top; i++ {
		w := app.w.Bots[i].Brain().Weights()
		cp := make([]float64, len(w))
		copy(cp, w)
		best = append(best, cp)
	}

	top = len(best)

	// app.goal.SetPos(float64(300+rand.Intn(500)), float64(100+rand.Intn(200)))
	for i, b := range app.w.Bots {
		b.Reset()
		//b.SetPos(float64(400+rand.Intn(400)), float64(600+rand.Intn(150)))
		b.SetPos(float64(540+rand.Intn(100)), float64(600+rand.Intn(150)))
		if i < top {
			continue
		}

		if top >= 1 {
			ix1 := 0
			ix2 := 0
			if top > 2 {
				ix1 = rand.Intn(top)
				ix2 = rand.Intn(top)
			} else if top > 1 {
				ix2 = 1
			}

			p1 := best[ix1]
			p2 := best[ix2]
			b.Brain().SetWeights(genetic.Reproduce(p1, p2, chance))
			continue
		}

		b.Brain().RandomWeights()
	}
}

func (app *App) Run(sleep time.Duration) (*world.World, <-chan struct{}, <-chan struct{}) {
	wait := make(chan struct{})
	tick := make(chan struct{})
	newGenCount := 2000

	go func() {
		var count int
		var dist float64
		var b *world.Bot
		var i int
		var score float64

		dists := make([]float64, len(app.w.Bots))
		for {
			for i, b = range app.w.Bots {
				dists[i] = b.Distance(app.goal)
			}
			app.w.Tick()
			count++

			if count%newGenCount == 0 {
				app.NewGeneration(3, 0.005, tick)
			}

			for _, b = range app.w.Bots {
				dist = b.Distance(app.goal)
				score = app.maxDist / (dist * dist)
				if dist < 2 {
					score = 1000
				}

				b.Reward(score)
			}

			if count == 1e6 {
				count = 0
			}

			if sleep != 0 {
				time.Sleep(sleep)
			}
		}
		wait <- struct{}{}
	}()

	return app.w, tick, wait
}
