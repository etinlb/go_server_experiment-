package main

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"time"
)

// const time
// func GameLoopGoRoutine()

func GameLoop() {
	tickChan := time.NewTicker(time.Millisecond * 16) // about 16 milliseconds for 60 fps a second

	go func() {
		lastTime := time.Now()
		for t := range tickChan.C {
			// update time vars
			timeElapsed := time.Since(lastTime).Seconds()
			lastTime = t
			// fmt.Printf("%+v T is\n", time.Since(lastTime).Seconds())
			DoPhysics(timeElapsed)
			lastTime = t
		}
	}()
	fmt.Println("Started Game Loop Go Routine at")
}

// TODO: Pass in the game objects as to be simulated rather than dangerously reading
// from the gameObjects global map
func DoPhysics(timePassed float64) {
	// TODO: figure out the implication of multi threading this yeah, this is
	for _, gameObj := range physicsComponents {
		// Basic movement for now
		gameObj.X += gameObj.XVel * timePassed
		gameObj.Y += gameObj.YVel * timePassed
		// TODO: This is dumb
		// packet := BuildObjectPackage("update", *gameObj)
		// clientBackend.BroadCastPackets(packet, make(map[*websocket.Conn]bool))
	}

}
