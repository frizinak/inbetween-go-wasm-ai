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

	r1 := world.NewWall(world.Rect(200, 200, 50, 100))
	bot := world.NewBot(0, 250, 1)
	x, y := bot.Center()

	l := 100.0
	angle := 0.0
	x1, y1 := x+math.Cos(angle)*l, y+math.Sin(angle)*l
	fmt.Println(x1, y1)

	angle = math.Pi / 4
	for i := 0.0; i < math.Pi*4; i += angle {
		fmt.Printf(
			"%t %4d %4d %4d %4d %5.2f %5.2f\n",
			r1.Intersected(x, y, bot.AbsDirection()),
			int(i*180/math.Pi),
			int(bot.AbsDirection()*180/math.Pi),
			int(bot.Wall.Direction(r1)*180/math.Pi),
			int(bot.Direction(r1)*180/math.Pi),
			bot.Direction(r1),
			bot.Direction(r1)/(math.Pi),
		)
		bot.Rotate(angle)
	}
	//os.Exit(0)

	a := app.New()
	_, tick, wait := a.Run(time.Duration(0))

	for range tick {
		fmt.Printf("top: %5.1f\n", a.MaxScore())
	}
	<-wait
}
