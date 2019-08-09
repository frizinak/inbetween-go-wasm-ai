package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/frizinak/inbetween-go-wasm-ai/app"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

func main() {
	go func() {
		panic(http.ListenAndServe(":8080", nil))
	}()
	rand.Seed(time.Now().UnixNano())

	// 	n1 := neural.New(neural.Sigmoid, []int{2, 8, 3})
	// 	n1.RandomWeights()
	//
	// 	n2 := neural.New(neural.Sigmoid, []int{2, 8, 3})
	// 	n2.RandomWeights()
	//
	// 	fmt.Println(n1.Input([]float64{0.3, 1.2}))
	// 	fmt.Println(n2.Input([]float64{0.3, 1.2}))
	//
	// 	w := genetic.Reproduce(n1.Weights(), n2.Weights(), 0.01)
	// 	child := neural.New(neural.Sigmoid, []int{2, 8, 3})
	// 	if err := child.SetWeights(w); err != nil {
	// 		panic(err)
	// 	}
	//
	// 	fmt.Println(child.Input([]float64{0.3, 1.2}))
	// 	fmt.Println(child.Input([]float64{0.3, 1.2}))
	// 	fmt.Println(child.Input([]float64{0.3, 1.2}))
	// 	fmt.Println(child.Input([]float64{0.3, 1.2}))
	//
	r1 := world.NewWall(world.Rect(200, 200, 5, 5))
	// r2 := world.NewWall(world.Rect(10, 450, 20, 30))
	//fmt.Println(r1.Distance(r2), r2.Distance(r1))
	// rad := r1.Direction(r2)
	//fmt.Println(rad, (rad * 180 / math.Pi))

	bot := world.NewBot(200, 0, 1)
	for i := 0.0; i < math.Pi*4; i += math.Pi / 4 {
		fmt.Printf(
			"%4d %4d %4d %5.2f %5.2f\n",
			int(i*180/math.Pi),
			int(bot.Wall.Direction(r1)*180/math.Pi),
			int(bot.Direction(r1)*180/math.Pi),
			bot.Direction(r1),
			bot.Direction(r1)/(math.Pi),
		)
		bot.Rotate(math.Pi / 4)
	}
	// fmt.Println(bot.Direction(r1) * 180 / math.Pi)

	//bot.Tick(r1, 500)

	fmt.Println(math.Pi * 2 * -1)
	fmt.Println(math.Pi * 2 * -0.5)
	fmt.Println(math.Pi * 2 * 0.5)
	fmt.Println(math.Pi * 2 * 1)
	//os.Exit(0)

	a := app.New()
	_, wait := a.Run(time.Duration(0))
	<-wait
}
