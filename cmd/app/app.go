package main

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"

	"github.com/frizinak/inbetween-go-wasm-ai/app"
	"github.com/frizinak/inbetween-go-wasm-ai/world"
)

func render(ctx js.Value, world *world.World, app *app.App) <-chan struct{} {
	wait := make(chan struct{})

	//ctx.Call("clearRect", 0, 0, world.Dx(), world.Dy())
	ctx.Set("fillStyle", "rgba(204, 204, 204, 0.10)")
	ctx.Call("fillRect", 0, 0, world.Dx(), world.Dy())

	for _, o := range world.Objects {
		draw(ctx, o, app.MaxScore())
	}

	for _, o := range world.Bots {
		draw(ctx, o, app.MaxScore())
	}

	go func() { wait <- struct{}{} }()
	return wait
}

func draw(ctx js.Value, o world.Object, maxScore float64) {
	if maxScore < 0 {
		maxScore = 0
	}
	clr := "black"
	switch v := o.(type) {
	case *world.Bot:
		s := v.Score()
		clr = "blue"
		// if v.ClosestType == "wall" {
		// 	clr = "black"
		// } else if v.ClosestType == "goal" {
		// 	clr = "red"
		// }
		if s > 0 {
			clr = fmt.Sprintf("rgb(%d, 30, 30)", int(s/maxScore*255))
		}
	case *world.Goal:
		clr = "red"
	}

	b := o.Bounds()
	ctx.Set("fillStyle", clr)
	ctx.Call("fillRect", b.Min.X, b.Min.Y, b.Dx(), b.Dy())
}

func main() {
	window := js.Global()
	document := window.Get("document")
	//body := document.Get("body")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	var rw sync.Mutex
	done := false
	a := app.New()
	world, _ := a.Run(time.Microsecond * 150)

	var anim js.Func
	anim = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			rw.Lock()
			wait := render(ctx, world, a)
			rw.Unlock()
			<-wait
			if done {
				anim.Release()
				return
			}
			window.Call("requestAnimationFrame", anim)
		}()
		return nil
	})

	go func() {
		window.Call("requestAnimationFrame", anim)
	}()

	c := make(chan struct{})
	<-c
}
