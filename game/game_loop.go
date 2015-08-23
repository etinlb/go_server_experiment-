package main

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"time"
)

// Various channel structs that are used for communicating with the game and physics loop
// TODO: Why is this the same as move message in message.go? Maybe make those the same
type MoveRequest struct {
	Xvel     int
	Yvel     int
	PlayerId string
}

// Adds an object with the id the game object map
type AddRequest struct {
	X        int
	Y        int
	PlayerId string
}

// const time
// func GameLoopGoRoutine()

// TODO: Figure out some sort of communication from the message thread to the
// game and physics thread
func RequestChange(id string) {

}

// Spawns the game loop and returns the channels to comminucate with the game
// TODO: Currently that is just the move channels, maybe return the ticker channel?
// TODO: TODO: Make it return channel of channels
func GameLoop() (chan MoveRequest, chan AddRequest) {
	// about 16 milliseconds for 60 fps a second
	tickChan := time.NewTicker(time.Millisecond * 1000)
	moveChannel := make(chan MoveRequest)
	addChannel := make(chan AddRequest)

	go func() {
		lastTime := time.Now()
		// Run the game loop forever.
		for t := range tickChan.C {
			fmt.Printf("%+v: ", t)
			fmt.Printf("%+v\n", addChannel)

			select {
			case msg := <-moveChannel:
				fmt.Printf("%+v", msg)
			case msg := <-addChannel:
				fmt.Printf("Added!!!!!! %+v\n", msg)
			default:
				// Move on to other things

			}
			// update time vars
			timeElapsed := time.Since(lastTime).Seconds()
			lastTime = t
			// fmt.Printf("%+v T is\n", time.Since(lastTime).Seconds())
			DoPhysics(timeElapsed)
			lastTime = t
		}
	}()
	fmt.Println("Started Game Loop Go Routine atlkj ")
	return moveChannel, addChannel

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
