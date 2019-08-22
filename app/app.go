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
	avgScore float64
	top      []float64
	gen      int
}

func New(nbots int) *App {
	app := &App{}

	height := 600.0
	width := 800.0

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

	// MAZE
	th := 0.003
	sp := 0.07
	pt := 0.04
	app.w.AddObject(world.NewWall(app.w.Relative(sp, pt, th, 1-pt)))
	app.w.AddObject(world.NewWall(app.w.Relative(sp+th, 0.2, sp-pt, th)))
	app.w.AddObject(world.NewWall(app.w.Relative(2*sp-(sp-pt), 0.3, sp-pt, th)))

	app.w.AddObject(world.NewWall(app.w.Relative(2*sp, 0.0, th, 1-pt)))

	// app.w.AddObject(world.NewWall(world.Rect(100, 100, 10, 700)))

	// app.w.AddObject(world.NewWall(world.Rect(110, 300, 70, 10)))

	// app.w.AddObject(world.NewWall(world.Rect(200, 0, 10, 600)))

	app.goal = world.NewGoal(app.w.Relative(0.98, 0.98, 0.01, 0.01))
	// END MAZE

	// // TRAINING SET 1
	// // |--   |
	// app.w.AddObject(world.NewWall(world.Rect(530, 180, 70, 20)))
	// app.w.AddObject(world.NewWall(world.Rect(600, 180, 20, 70)))
	// // |   --|
	// app.w.AddObject(world.NewWall(world.Rect(600, 300, 70, 20)))
	// // |--   |
	// app.w.AddObject(world.NewWall(world.Rect(530, 400, 70, 20)))
	// // |--  -|
	// app.w.AddObject(world.NewWall(world.Rect(530, 500, 70, 20)))
	// app.w.AddObject(world.NewWall(world.Rect(630, 500, 40, 20)))

	// app.w.AddObject(world.NewWall(world.Rect(400, 100, 20, 80)))

	// app.w.AddObject(world.NewWall(world.Rect(650, 0, 20, 800)))
	// app.w.AddObject(world.NewWall(world.Rect(530, 180, 20, 620)))

	// app.w.AddObject(world.NewWall(world.Rect(320, 180, 220, 20)))
	// app.w.AddObject(world.NewWall(world.Rect(200, 50, 350, 20)))

	// app.w.AddObject(world.NewWall(world.Rect(200, 300, 120, 20)))
	// app.w.AddObject(world.NewWall(world.Rect(320, 180, 20, 140)))

	// app.w.AddObject(world.NewWall(world.Rect(200, 50, 20, 250)))

	// app.w.AddObject(world.NewWall(world.Rect(530, 0, 20, 50)))

	// app.goal = world.NewGoal(world.Rect(250, 230, 50, 50))
	// // END TRAINING SET 1

	app.w.AddObject(app.goal)

	app.w.AddObject(world.NewWall(world.Rect(-20, -20, 20, height+40)))
	app.w.AddObject(world.NewWall(world.Rect(-20, -20, width+40, 20)))
	app.w.AddObject(world.NewWall(world.Rect(width, -20, 20, height+40)))
	app.w.AddObject(world.NewWall(world.Rect(-20, height, width+40, 20)))

	return app
}

func (app *App) botPos(bot *world.Bot) {
	bot.SetPos(
		app.w.Dx()*(0.01+rand.Float64()/50),
		app.w.Dy()*(0.8+rand.Float64()/8),
	)
}

func (app *App) World() *world.World {
	return app.w
}

func (app *App) MaxScore() float64 {
	return app.maxScore
}

func (app *App) AvgScore() float64 {
	return app.avgScore
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

func (app *App) NewGeneration(top int, chance float64, tick chan struct{}) {
	sort.Slice(app.w.Bots, func(i, j int) bool {
		return app.w.Bots[i].Score() > app.w.Bots[j].Score()
	})

	app.top = app.w.Bots[0].Brain().Weights()
	app.maxScore = app.w.Bots[0].Score()
	var tot float64
	for _, b := range app.w.Bots {
		s := b.Score()
		if s > 0 {
			tot += s
		}
	}
	app.avgScore = tot / float64(len(app.w.Bots))
	app.gen++

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
		//550 - 650
		app.botPos(b)
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

func (app *App) Run(sleep time.Duration, n int) (<-chan struct{}, <-chan struct{}, chan<- struct{}) {
	wait := make(chan struct{})
	tick := make(chan struct{})
	stop := make(chan struct{})
	newGenCount := 2000

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

			if count%newGenCount == 0 {
				count = 0
				app.NewGeneration(2, 0.005, tick)
				//app.NewGeneration(3, 0.01, tick)
				if times {
					n--
					if n <= 0 {
						break
					}
				}
			}

			for i, b = range app.w.Bots {
				dist = b.Distance(app.goal)
				score = 0
				if dist < dists[i] {
					score += 0.2 //* b.Speed()
				}
				score += app.maxDist / (10 * dist)
				if dist < 5 {
					score = 1000
				} else if b.Speed() < 0.2 {
					score = -2
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
