package world

import (
	"math"
	"runtime"
	"sync"

	"github.com/frizinak/inbetween-go-wasm-ai/geometry"
)

func Rect(x, y float64, width, height float64) geometry.Rectangle {
	r := geometry.Rect(x, y, x+width, y+height)
	return r
}

type World struct {
	geometry.Rectangle
	Bots        []*Bot
	Objects     []Object
	maxDistance float64
	tick        chan botRange
	wg          sync.WaitGroup
}

type botRange struct {
	min int
	max int
}

func New(r geometry.Rectangle) *World {
	workers := runtime.NumCPU()
	w := &World{
		r,
		make([]*Bot, 0),
		make([]Object, 0),
		math.Sqrt(r.Dx()*r.Dx() + r.Dy()*r.Dy()),
		make(chan botRange, workers),
		sync.WaitGroup{},
	}

	for i := 0; i < workers; i++ {
		go func() {
			for t := range w.tick {
				for _, b := range w.Bots[t.min:t.max] {
					botTick(w, b)
				}
				w.wg.Done()
			}
		}()
	}

	return w
}

func (w *World) MaxDistance() float64 {
	return w.maxDistance
}

func (w *World) AddObject(o Object) {
	w.Objects = append(w.Objects, o)
}

func (w *World) AddBot(b *Bot) {
	w.Bots = append(w.Bots, b)
}

func (w *World) Collision(o1, o2 Object) geometry.Rectangle {
	if !o1.Solid() || !o2.Solid() {
		return geometry.ZR
	}

	r := o1.Bounds().Intersect(o2.Bounds())
	return r
}

func (w *World) Tick() {
	w.wg = sync.WaitGroup{}
	total := len(w.Bots)
	workers := cap(w.tick)
	size := total / workers
	if size < 1 {
		size = 1
	}

	for i := 0; i < total; i += size {
		w.wg.Add(1)
		d := botRange{i, i + size}
		if d.max > total {
			d.max = total
		}
		w.tick <- d
	}

	w.wg.Wait()
	return
}

type job struct {
	b *Bot
}

func botTick(w *World, b *Bot) {
	var closest Object
	var o Object
	var dist float64

	pos := b.Bounds()
	closestDist := math.MaxFloat64
	dir := b.AbsDirection()
	x1, y1, x2, y2 := b.Sides()
	x1a, y1a := geometry.Line(x1, y1, dir, w.maxDistance)
	x2a, y2a := geometry.Line(x2, y2, dir, w.maxDistance)
	for _, o = range w.Objects {
		if !o.Intersected(x1, y1, x1a, y1a) && !o.Intersected(x2, y2, x2a, y2a) {
			continue
		}

		dist = b.Distance(o)
		if dist < closestDist {
			closest = o
			closestDist = dist
		}
	}

	b.Tick(closest, closestDist, w.maxDistance)

	score := 0.0
	for _, o = range w.Objects {
		d := o.Distance(b)
		if d == 0 {
			c := w.Collision(b, o)
			if c.Empty() {
				continue
			}

			score -= 20
			b.Translate(pos.Min.X-b.Min.X, pos.Min.Y-b.Min.Y)
		}
	}

	b.Reward(score)
}
