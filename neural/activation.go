package neural

func Sigmoid(x float64) float64 {
	//return 1 / (1 + math.Exp(-x))
	return 1 / (1 + Exp(-x))

	if x >= 4 {
		return 1
	}
	tmp := 1 - 0.25*x
	tmp *= tmp
	tmp *= tmp
	tmp *= tmp
	tmp *= tmp
	return 1 / (1 + tmp)
}

func Exp(x float64) float64 {
	return (120 + x*(120+x*(60+x*(20+x*(5+x))))) * 0.0083333333
}
