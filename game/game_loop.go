package main

import (
	"fmt"
	"time"
)

func GameLoop() {
	// frameRate := time.Millisecond * (1 / 60 * 1000) // number of milliseconds for the frame
	// lastUpdate = time.Now()
	tickChan := time.NewTicker(time.Millisecond * 1000) // about 16 milliseconds for 60 fps a second

	go func() {
		for range tickChan.C {
			for _, gameObj := range gameObjects {
				gameObj.Update()
			}
		}
	}()

	fmt.Println("Started Game Loop Go Routine at")
}
