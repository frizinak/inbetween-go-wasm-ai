package neural

import (
	"errors"
	"math/rand"
)

type Activation func(x float64) float64

type Network struct {
	activation Activation
	layout     []int
	weights    []float64
	output     []float64
	bias       []float64
	buf        [][]float64
}

func New(activation Activation, layout []int) *Network {
	total := 0
	for i := 0; i < len(layout)-1; i++ {
		total += layout[i] * layout[i+1]
	}

	bias := make([]float64, len(layout)-1)
	for i := 0; i < len(layout)/2+1; i++ {
		bias[i] = 1
	}

	buf := make([][]float64, len(layout))
	for layer := 0; layer < len(layout)-1; layer++ {
		buf[layer] = make([]float64, layout[layer+1])
	}

	return &Network{
		activation,
		layout,
		make([]float64, total),
		make([]float64, layout[len(layout)-1]),
		bias,
		buf,
	}
}

func (n *Network) RandomWeights() {
	for i := range n.weights {
		n.weights[i] = rand.Float64() + rand.Float64() - 1.0
	}
}

func (n *Network) Weights() []float64 {
	return n.weights
}

func (n *Network) SetWeights(weights []float64) error {
	if len(weights) != len(n.weights) {
		return errors.New("invalid length")
	}

	n.weights = weights
	return nil
}

func (n *Network) Input(inputs []float64) []float64 {
	var total float64
	var offset int

	var inp, outp int
	var in, out int
	var layer int

	for layer = 0; layer < len(n.layout)-1; layer++ {
		in = n.layout[layer]
		out = n.layout[layer+1]

		for outp = 0; outp < out; outp++ {
			total = 0
			for inp = 0; inp < in; inp++ {
				total += inputs[inp] * n.weights[offset]
				offset += 1
			}

			total = n.activation(total + n.bias[layer])
			n.buf[layer][outp] = total
		}

		inputs = n.buf[layer]
	}

	return inputs
}
