package neural

func Sigmoid(x float64) float64 {
	return 1 / (1 + exp(-x))
}

func exp(x float64) float64 {
	//return (720 + x*(720+x*(360+x*(120+x*(30+x*(6+x)))))) * 0.0013888888
	return (120 + x*(120+x*(60+x*(20+x*(5+x))))) * 0.0083333333
}
