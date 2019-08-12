package app

import (
	"math/rand"
	"sort"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/genetic"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

type App struct {
	w        *world.World
	goal     world.Object
	maxDist  float64
	maxScore float64
}

func New() *App {

	app := &App{}

	nbots := 300
	height := 800.0
	width := 1200.0

	app.w = world.New(world.Rect(0, 0, width, height))
	app.maxDist = app.w.MaxDistance()

	bots := make([]*world.Bot, 0, nbots)
	for i := 0; i < nbots; i++ {
		bots = append(bots, world.NewBot(600, 600, 3))
	}

	for _, bot := range bots {
		app.w.AddBot(bot)
	}

	app.w.AddObject(world.NewWall(world.Rect(400, 350, 400, 20)))
	app.goal = world.NewGoal(world.Rect(500, 200, 50, 50))
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

func (app *App) NewGeneration(top int, chance float64, tick chan struct{}) {
	sort.Slice(app.w.Bots, func(i, j int) bool {
		return app.w.Bots[i].Score() > app.w.Bots[j].Score()
	})

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

	app.goal.SetPos(float64(300+rand.Intn(500)), float64(100+rand.Intn(200)))
	for i, b := range app.w.Bots {
		b.Reset()
		b.SetPos(float64(400+rand.Intn(400)), float64(600+rand.Intn(150)))
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
	newGenCount := 1000

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
				app.NewGeneration(6, 0.08, tick)
			}

			for _, b = range app.w.Bots {
				dist = b.Distance(app.goal)
				score = app.maxDist / (dist * dist)
				if dist < 2 {
					score = 1000
				}

				if dist > 100 && b.Speed() < 0.2 {
					score -= 1
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
