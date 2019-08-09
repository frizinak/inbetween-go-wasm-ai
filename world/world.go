package world

import (
	"math"

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
}

func New(r geometry.Rectangle) *World {
	return &World{
		r,
		make([]*Bot, 0),
		make([]Object, 0),
		math.Sqrt(r.Dx()*r.Dx() + r.Dy()*r.Dy()),
	}
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
	var score float64
	var pos geometry.Rectangle
	var b *Bot
	var o Object
	var closest Object
	var dist float64
	var closestDist float64

	for _, b = range w.Bots {
		pos = b.Bounds()
		closest = nil
		closestDist = math.MaxFloat64
		for _, o = range w.Objects {
			dist = b.Distance(o)
			if dist < closestDist {
				closest = o
				closestDist = dist
			}
		}

		b.Tick(closest, w.maxDistance)

		score = 0
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

		//if score != 0 {
		//	fmt.Println(score)
		//}
		b.Reward(score)
	}
}
