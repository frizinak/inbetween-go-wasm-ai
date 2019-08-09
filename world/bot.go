package world

import (
	"math"

	"github.com/frizinak/inbetween-go-wasm-ai/neural"
)

const Pi2 = 2 * math.Pi

type Bot struct {
	*Wall

	brain *neural.Network

	score float64

	speed     float64
	direction float64

	maxSpeed float64

	ClosestType string
}

func NewBot(x, y float64, maxSpeed float64) *Bot {
	b := neural.New(neural.Sigmoid, []int{4, 16, 2})
	// b := neural.New(neural.Sigmoid, []int{4, 18, 32, 18, 2})
	b.RandomWeights()
	return &Bot{
		Wall:      NewWall(Rect(x, y, 8, 8)),
		brain:     b,
		direction: 0,
		maxSpeed:  maxSpeed,
	}
}

func (b *Bot) Reward(i float64) {
	b.score += i
}

func (b *Bot) Score() float64 {
	return b.score
}

func (b *Bot) Brain() *neural.Network {
	return b.brain
}

func (b *Bot) Reset() {
	b.score = 0
}

func (b *Bot) Direction(o Object) float64 {
	obj := b.Wall.Direction(o)
	own := b.direction * Pi2
	// for own < 0 {
	// 	own += Pi2
	// }
	// for own > Pi2 {
	// 	own -= Pi2
	// }

	diff := obj - own
	for diff < -math.Pi {
		diff += Pi2
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

func (b *Bot) Tick(o Object, maxDistance float64) {
	input := make([]float64, 4)
	input[0] = b.speed / b.maxSpeed
	///input[1] = b.direction
	input[1] = b.Distance(o) / maxDistance
	input[2] = b.Direction(o) / math.Pi

	b.ClosestType = "wall"
	switch o.(type) {
	case *Goal:
		b.ClosestType = "goal"
		input[3] = 1
	}

	output := b.brain.Input(input)
	//fmt.Println(output)
	if output[0] > 0.9 {
		b.direction += 0.04 * (output[0] - 0.5)
	} else if output[0] < 0.1 {
		b.direction -= 0.04 * ((1 - output[0]) - 0.5)
	}

	if b.direction < 0 {
		b.direction += 1
	} else if b.direction > 1 {
		b.direction -= 1
	}

	b.speed -= 0.1
	if output[1] > 0.5 {
		b.speed += (output[1] - 0.5) - 0.25
	}

	// if output[1] > 0.9 {
	// 	b.speed += b.maxSpeed * 0.1
	// } else if output[1] < 0.1 {
	// 	b.speed -= b.maxSpeed * 0.1
	// }

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
