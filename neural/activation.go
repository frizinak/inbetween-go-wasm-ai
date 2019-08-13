package neural

func Sigmoid(x float64) float64 {
	return 1 / (1 + Exp(-x))
}

func Exp(x float64) float64 {
	return (120 + x*(120+x*(60+x*(20+x*(5+x))))) * 0.0083333333
}
