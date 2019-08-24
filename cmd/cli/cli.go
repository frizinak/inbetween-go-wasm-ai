package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/frizinak/inbetween-go-wasm-ai/app"
)

func main() {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	rand.Seed(time.Now().UnixNano())

	a := app.New(300)
	//a.Import(string(bound.MustAsset("weights1.txt")))
	tick, wait, _ := a.Run(time.Duration(0), -1)

	last := time.Now()
	for range tick {
		fmt.Printf("top: %5.1f median: %5.1f\n", a.MaxScore(), a.MedianScore())

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
