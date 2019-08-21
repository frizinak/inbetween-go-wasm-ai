package world

import (
	"math"
	"runtime"

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
	jobs        chan *job
	jobsDone    chan struct{}
}

func New(r geometry.Rectangle) *World {
	workers := runtime.NumCPU()
	w := &World{
		r,
		make([]*Bot, 0),
		make([]Object, 0),
		math.Sqrt(r.Dx()*r.Dx() + r.Dy()*r.Dy()),
		make(chan *job, workers),
		make(chan struct{}, workers),
	}

	for i := 0; i < workers; i++ {
		go func() {
			for job := range w.jobs {
				job.Tick(w, job.b)
				w.jobsDone <- struct{}{}
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
	n := len(w.Bots)
	done := make(chan struct{})
	go func() {
		for i := 0; i < n; i++ {
			<-w.jobsDone
		}
		done <- struct{}{}
	}()
	for _, b := range w.Bots {
		w.jobs <- &job{b}
	}
	<-done
}

type job struct {
	b *Bot
}

func (job *job) Tick(w *World, b *Bot) {
	var closest Object
	var o Object
	var dist float64

	pos := b.Bounds()
	closestDist := math.MaxFloat64
	dir := b.AbsDirection()
	x, y := b.Center()
	for _, o = range w.Objects {
		if !o.Intersected(x, y, dir) {
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

			score -= 20000
			b.Translate(pos.Min.X-b.Min.X, pos.Min.Y-b.Min.Y)
		}
	}

	b.Reward(score)
}
