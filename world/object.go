package world

import (
	"github.com/frizinak/inbetween-go-wasm-ai/geometry"
)

type Object interface {
	Bounds() geometry.Rectangle
	Solid() bool
	Distance(Object) float64
	Direction(Object) float64
	Intersected(x1, y1, x2, y2 float64) bool
	Translate(dx, dy float64)
	SetPos(x, y float64)
}

type Wall struct {
	geometry.Rectangle
}

func NewWall(r geometry.Rectangle) *Wall {
	return &Wall{r}
}

func (w *Wall) Bounds() geometry.Rectangle {
	return w.Rectangle
}

func (w *Wall) Solid() bool {
	return true
}

func (w *Wall) Translate(dx, dy float64) {
	w.Rectangle = w.Rectangle.Translate(dx, dy)
}

func (w *Wall) SetPos(x, y float64) {
	w.Translate(x-w.Min.X, y-w.Min.Y)
}

func (w *Wall) Distance(o Object) float64 {
	return w.Rectangle.Distance(o.Bounds())
}

func (w *Wall) Direction(o Object) float64 {
	return w.Rectangle.Direction(o.Bounds())
}

type Goal struct {
	Wall
}

func (g *Goal) Solid() bool {
	return false
}

func NewGoal(r geometry.Rectangle) *Goal {
	return &Goal{Wall{r}}
}
