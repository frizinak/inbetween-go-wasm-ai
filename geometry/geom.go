package geometry

import (
	"math"
)

var ZR Rectangle

type Point struct {
	X float64
	Y float64
}

type Rectangle struct {
	Min Point
	Max Point
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Min: Point{x0, y0}, Max: Point{x1, y1}}
}

func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Translate(dx, dy float64) Rectangle {
	return Rectangle{
		Min: Point{r.Min.X + dx, r.Min.Y + dy},
		Max: Point{r.Max.X + dx, r.Max.Y + dy},
	}
}

func (r Rectangle) IntersectedDistance(x, y, angle float64) (bool, float64) {
	l := 1000000.0
	x2, y2 := x+math.Cos(angle)*l, y+math.Sin(angle)*l
	var dist float64 = math.MaxFloat64
	var d float64
	var i, b bool
	var sx, sy float64
	if i, sx, sy = LineIntersection(x, y, x2, y2, r.Min.X, r.Min.Y, r.Min.X, r.Max.Y); i {
		b = true
		if d = Distance(x, y, sx, sy); d < dist {
			dist = d
		}
	}

	if i, sx, sy = LineIntersection(x, y, x2, y2, r.Max.X, r.Min.Y, r.Max.X, r.Max.Y); i {
		b = true
		if d = Distance(x, y, sx, sy); d < dist {
			dist = d
		}
	}
	if i, sx, sy = LineIntersection(x, y, x2, y2, r.Min.X, r.Min.Y, r.Max.X, r.Min.Y); i {
		b = true
		if d = Distance(x, y, sx, sy); d < dist {
			dist = d
		}
	}
	if i, sx, sy = LineIntersection(x, y, x2, y2, r.Min.X, r.Max.Y, r.Max.X, r.Max.Y); i {
		b = true
		if d = Distance(x, y, sx, sy); d < dist {
			dist = d
		}
	}

	return b, dist

}

func LineIntersection(x1, y1, x2, y2, x3, y3, x4, y4 float64) (bool, float64, float64) {
	a := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))
	if a < 0 || a > 1 {
		return false, 0, 0
	}
	b := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))
	if b < 0 || b > 1 {
		return false, 0, 0
	}

	return true, x3 + (a * (x4 - x3)), y3 + (b * (y4 - y3))
}

func (r Rectangle) Intersect(o Rectangle) Rectangle {
	if r.Min.X < o.Min.X {
		r.Min.X = o.Min.X
	}
	if r.Min.Y < o.Min.Y {
		r.Min.Y = o.Min.Y
	}
	if r.Max.X > o.Max.X {
		r.Max.X = o.Max.X
	}
	if r.Max.Y > o.Max.Y {
		r.Max.Y = o.Max.Y
	}

	return r
}

func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle) Distance(o Rectangle) float64 {
	w := r.Max.X < o.Min.X
	e := o.Max.X < r.Min.X
	n := r.Max.Y < o.Min.Y
	s := o.Max.Y < r.Min.Y
	switch {
	case n && w:
		return Distance(r.Max.X, r.Max.Y, o.Min.X, o.Min.Y)
	case n && e:
		return Distance(r.Min.X, r.Max.Y, o.Max.X, o.Min.Y)
	case s && w:
		return Distance(r.Max.X, r.Min.Y, o.Min.X, o.Max.Y)
	case s && e:
		return Distance(r.Min.X, r.Min.Y, o.Max.X, o.Max.Y)
	case w:
		return o.Min.X - r.Max.X
	case e:
		return r.Min.X - o.Max.X
	case n:
		return o.Min.Y - r.Max.Y
	case s:
		return r.Min.Y - o.Max.Y
	}

	return 0
}

func (r Rectangle) Direction(o Rectangle) float64 {
	x1 := r.Min.X + r.Dx()/2
	y1 := r.Min.Y + r.Dy()/2
	x2 := o.Min.X + o.Dx()/2
	y2 := o.Min.Y + o.Dy()/2

	return math.Atan2(y1-y2, x2-x1)
}

func Distance(x1, y1, x2, y2 float64) float64 {
	x := x1 - x2
	y := y1 - y2
	return math.Sqrt(x*x + y*y)
}
