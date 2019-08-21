package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/frizinak/inbetween-go-wasm-ai/app"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

func main() {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	rand.Seed(time.Now().UnixNano())

	r1 := world.NewWall(world.Rect(200, 200, 50, 100))
	bot := world.NewBot(198, 0, 1)
	x, y := bot.Center()

	l := 100.0
	angle := 0.0
	x1, y1 := x+math.Cos(angle)*l, y+math.Sin(angle)*l
	fmt.Println(x, y, x1, y1)

	angle = math.Pi / 4
	for i := 0.0; i < math.Pi*4; i += angle {
		lx, ly, rx, ry := bot.Sides()
		//fmt.Println("l, r", int(bot.AbsDirection()*180/math.Pi), lx, ly, rx, ry)
		fmt.Printf(
			"%t %t %4d %4d %4d %4d %5.2f %5.2f\n",
			r1.Intersected(lx, ly, bot.AbsDirection()),
			r1.Intersected(rx, ry, bot.AbsDirection()),
			int(i*180/math.Pi),
			int(bot.AbsDirection()*180/math.Pi),
			int(bot.Wall.Direction(r1)*180/math.Pi),
			int(bot.Direction(r1)*180/math.Pi),
			bot.Direction(r1),
			bot.Direction(r1)/(math.Pi),
		)
		bot.Rotate(angle)
	}

	a := app.New()
	//a.Import(string(bound.MustAsset("weights1.txt")))
	tick, wait, _ := a.Run(time.Duration(0), -1)

	last := time.Now()
	for range tick {
		fmt.Printf("top: %5.1f avg: %5.1f\n", a.MaxScore(), a.AvgScore())

		now := time.Now()
		if now.Sub(last) > time.Second*10 {
			last = now
			fmt.Println("written /tmp/brain")
			f, err := os.Create("/tmp/brain.tmp")
			if err != nil {
				panic(err)
			}
			f.WriteString(a.Export())
			f.Close()
			os.Rename("/tmp/brain.tmp", "/tmp/brain")
		}
	}
	<-wait
}
