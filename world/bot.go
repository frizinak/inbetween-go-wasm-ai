package world

import (
	"math"
	"math/rand"

	"github.com/frizinak/inbetween-go-wasm-ai/neural"
)

const Pi2 = 2 * math.Pi

type Bot struct {
	*Wall

	brain *neural.Network

	fitness float64

	speed     float64
	direction float64

	maxSpeed float64

	ClosestType string
}

func NewBot(x, y float64, maxSpeed float64) *Bot {
	b := neural.New(neural.Sigmoid, []int{3, 14, 10, 2})
	b.RandomWeights()
	return &Bot{
		Wall:      NewWall(Rect(x, y, 12, 12)),
		brain:     b,
		direction: rand.Float64(),
		maxSpeed:  maxSpeed,
	}
}

func (b *Bot) Reward(i float64) {
	b.fitness += i
}

func (b *Bot) Fitness() float64 {
	return b.fitness
}

func (b *Bot) Brain() *neural.Network {
	return b.brain
}

func (b *Bot) Reset(direction float64) {
	b.speed = 0
	b.fitness = 0
	b.direction = direction //rand.Float64()
}

func (b *Bot) Speed() float64 {
	return b.speed
}

func (b *Bot) setSpeed(s float64) {
	b.speed = s
}

func (b *Bot) Sides() (x1, y1, x2, y2 float64) {
	r := b.Dx()/2 + 4
	x, y := b.Center()
	angle := b.AbsDirection() - math.Pi/2
	cos, sin := math.Cos(angle)*r, math.Sin(angle)*r

	x1 = x - cos
	y1 = y - sin
	x2 = x + cos
	y2 = y + sin

	return
}

func (b *Bot) Center() (x, y float64) {
	x, y = b.Min.X+b.Dx()/2, b.Min.Y+b.Dy()/2
	return
}

func (b *Bot) AbsDirection() float64 {
	return b.direction * Pi2
}

func (b *Bot) Direction(o Object) float64 {
	obj := b.Wall.Direction(o)
	own := b.direction * Pi2
	diff := obj - own
	for diff < -math.Pi {
		diff += Pi2
	}
	for diff > math.Pi {
		diff -= Pi2
	}
	return diff
}

func (b *Bot) Rotate(rad float64) {
	dir := b.direction*Pi2 + rad
	for dir > Pi2 {
		dir -= Pi2
	}
	for dir < 0 {
		dir += Pi2
	}

	b.direction = dir / Pi2
}

func (b *Bot) Tick(o Object, dist, maxDistance float64) {
	input := []float64{b.speed / b.maxSpeed, 1, 0}
	if o != nil && dist < b.maxSpeed*5 {
		input[1] = dist / maxDistance
		// input[2] = dir/Pi2 + 0.5
		b.ClosestType = "wall"
		switch o.(type) {
		case *Goal:
			b.ClosestType = "goal"
			input[2] = 1
		}
	}

	output := b.brain.Input(input)
	if output[0] > 0.9 {
		b.direction += 0.03
	} else if output[0] > 0.8 {
		b.direction -= 0.03
	}

	if b.direction < -1 {
		b.direction += 1
	} else if b.direction > 1 {
		b.direction -= 1
	}

	if output[1] > 0.9 {
		b.speed += 0.02
	} else if output[1] > 0.8 {
		b.speed -= 0.5
	}

	if b.speed > b.maxSpeed {
		b.speed = b.maxSpeed
	} else if b.speed < 0 {
		b.speed = 0
	}

	angle := 2 * math.Pi * b.direction
	dx := math.Cos(angle) * b.speed
	dy := math.Sin(angle) * b.speed
	b.Translate(dx, dy)
}

func (b *Bot) Solid() bool {
	return true
}
