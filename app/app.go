package app

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/genetic"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

type App struct {
	bestScore float64
}

func New() *App {
	return &App{}
}

func (a *App) MaxScore() float64 {
	return a.bestScore
}

func (a *App) NewGeneration(bots []*world.Bot, goal world.Object, top int, chance float64) {
	sort.Slice(bots, func(i, j int) bool {
		return bots[i].Score() > bots[j].Score()
	})

	a.bestScore = bots[0].Score()
	allNegative := a.bestScore < 0
	fmt.Println(a.bestScore, bots[0].Distance(goal))

	best := make([][]float64, 0, top)
	for i := 0; i < top; i++ {
		if !allNegative && bots[i].Score() < 0 {
			break
		}

		w := bots[i].Brain().Weights()
		cp := make([]float64, len(w))
		copy(cp, w)
		best = append(best, cp)
	}

	top = len(best)

	if !allNegative {
		goal.SetPos(float64(300+rand.Intn(500)), float64(100+rand.Intn(200)))
	}

	for i, b := range bots {
		b.Reset()
		// b.SetPos(float64(100+rand.Intn(1100)), float64(500+rand.Intn(150)))
		b.SetPos(600, 600)
		if i < top && !allNegative {
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
			if !allNegative || (allNegative && rand.Float64() > 0.8) {
				continue
			}
		}

		b.Brain().RandomWeights()

	}
}

func (a *App) Run(sleep time.Duration) (*world.World, <-chan struct{}) {
	wait := make(chan struct{})
	rand.Seed(time.Now().UnixNano())

	height := 800.0
	width := 1200.0
	w := world.New(world.Rect(0, 0, width, height))
	maxDist := w.MaxDistance()

	/// w.AddObject(world.NewWall(world.Rect(230, 230, 20, 20)))
	/// w.AddObject(world.NewWall(world.Rect(500, 350, 20, 40)))
	/// w.AddObject(world.NewWall(world.Rect(800, 200, 20, 120)))
	goal := world.NewGoal(world.Rect(500, 200, 5, 5))
	w.AddObject(goal)

	w.AddObject(world.NewWall(world.Rect(-20, -20, 20, height+40)))
	w.AddObject(world.NewWall(world.Rect(-20, -20, width+40, 20)))
	w.AddObject(world.NewWall(world.Rect(width, -20, 20, height+40)))
	w.AddObject(world.NewWall(world.Rect(-20, height, width+40, 20)))

	nbots := 300
	bots := make([]*world.Bot, 0, nbots)
	for i := 0; i < nbots; i++ {
		bots = append(bots, world.NewBot(600, 600, 3))
	}

	for _, bot := range bots {
		w.AddBot(bot)
	}

	newGenCount := 1000

	go func() {
		var count int
		var dist float64
		//var dir float64
		var b *world.Bot
		var i int
		var score float64

		// dirs := make([]float64, nbots)
		dists := make([]float64, nbots)
		for {
			for i, b = range bots {
				dists[i] = b.Distance(goal)
				//dirs[i] = b.Direction(goal)
			}
			w.Tick()
			count++

			if count%newGenCount == 0 {
				a.NewGeneration(bots, goal, 4, 0.02)
			}

			for i, b = range bots {
				dist = b.Distance(goal)
				// dir = b.Direction(goal)

				score = -0.1
				// score = -math.Sqrt(dist / maxDist)
				if dist > 1 && dist < dists[i] {
					score += maxDist / (dist * dist)
				}

				// score -= 0.5
				// if dist > dists[i] {
				// 	score -= 5
				// }

				if dist <= dists[i] {
					if dist <= 5 {
						score += 10000
					} else if dist < maxDist/100 {
						score += 1000
					} else if dist < maxDist/25 {
						score += 100
					} else if dist < maxDist/15 {
						score += 10
					}
				}

				// if dist > 200 {
				// 	score -= 3
				// 	if dist < dists[i] {
				// 		score += 6
				// 	}
				// }

				// if dist < 2 {
				// 	score += 200
				// } else if dist < 10 {
				// 	score += 30
				// } else if dist < 50 {
				// 	score += 15
				// } else if dist < 100 {
				// 	score += 7
				// } else if dist < 200 {
				// 	score += 3
				// }

				// if dist > 100 {
				//score -= 3
				// if dir > -0.15 && dir < 0.15 {
				// 	//score += 5
				// } else if math.Abs(dir) < math.Abs(dirs[i]) {
				// 	//score += 5
				// }
				/// }
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

	return w, wait
}
