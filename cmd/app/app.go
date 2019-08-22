// +build js,wasm

package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/app"
	"github.com/frizinak/inbetween-go-wasm-ai/bound"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

func render(ctx js.Value, app *app.App) {
	world := app.World()
	ctx.Set("fillStyle", "rgba(204, 204, 204, 0.30)")
	ctx.Call("fillRect", 0, 0, world.Dx(), world.Dy())

	for _, o := range world.Objects {
		draw(ctx, o, app.MaxScore())
	}

	for _, o := range world.Bots {
		draw(ctx, o, app.MaxScore())
	}
}

func draw(ctx js.Value, o world.Object, maxScore float64) {
	if maxScore < 0 {
		maxScore = 0
	}

	b := o.Bounds()
	clr := "rgba(0,0,0,0.5)"
	switch v := o.(type) {
	case *world.Bot:
		// if v.ClosestType == "wall" {
		// 	clr = "black"
		// } else if v.ClosestType == "goal" {
		// 	clr = "red"
		// }

		s := v.Score()
		if s > 0 {
			clr = fmt.Sprintf("rgb(%d, 30, 30)", int(s/maxScore*255))
		}

		cx, cy := v.Center()
		a := v.AbsDirection() + math.Pi/4
		beak := 0.5
		ctx.Set("fillStyle", clr)
		ctx.Call("beginPath")
		ctx.Call("arc", cx, cy, b.Dx()/2, a-beak, a-beak+math.Pi*2-beak)
		ctx.Call("lineTo", cx, cy)
		ctx.Call("fill")
		return
	case *world.Goal:
		clr = "red"
	}

	ctx.Set("fillStyle", clr)
	ctx.Call("fillRect", b.Min.X, b.Min.Y, b.Dx(), b.Dy())
}

func main() {
	window := js.Global()
	document := window.Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	btnFast := document.Call("getElementById", "fast")
	btnImport := document.Call("getElementById", "import")
	btnExport := document.Call("getElementById", "export")
	textarea := document.Call("getElementById", "brain")
	stats := document.Call("getElementById", "stats")

	textarea.Set("value", string(bound.MustAsset("weights1.txt")))

	done := false
	a := app.New(80)

	var tick <-chan struct{}
	var wait <-chan struct{}
	var stop chan<- struct{}
	interval := time.Microsecond * 250

	run := func(iv time.Duration, n int) {
		fmt.Println("run with", iv, n)
		tick, wait, stop = a.Run(iv, n)
		for range tick {
			stats.Set(
				"innerText",
				fmt.Sprintf(
					"Generation: %d\nMax score: %5.2f\nAvg Score: %5.2f\n",
					a.Generation(),
					a.MaxScore(),
					a.AvgScore(),
				),
			)
		}
	}
	go run(interval, -1)

	var anim js.Func
	anim = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			render(ctx, a)
			if done {
				anim.Release()
				return
			}
			window.Call("requestAnimationFrame", anim)
		}()
		return nil
	})

	window.Call("requestAnimationFrame", anim)

	btnExport.Call(
		"addEventListener",
		"click",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			textarea.Set("value", a.Export())
			return nil
		}),
	)
	btnImport.Call(
		"addEventListener",
		"click",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			err := a.Import(textarea.Get("value").String())
			if err != nil {
				window.Call("alert", err.Error())
			}
			return nil
		}),
	)
	btnFast.Call(
		"addEventListener",
		"click",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go func() {
				stop <- struct{}{}
				<-wait
				now := time.Now()
				run(0, 5)
				fmt.Println(time.Now().Sub(now))
				go run(interval, -1)
			}()
			return nil
		}),
	)

	c := make(chan struct{})
	<-c
}
