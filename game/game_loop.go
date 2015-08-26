package main

import (
	"fmt"
	"time"
)

// Spawns the game loop and returns the channels to comminucate with the game
// TODO: Currently that is just the move channels, maybe return the ticker channel?
// TODO: TODO: Make it return channel of channels
func StartGameLoop() (chan *MoveRequest, chan *AddRequest) {
	// about 16 milliseconds for 60 fps a second
	gameTick := time.NewTicker(time.Millisecond * 1000)

	// Physics runs at 50 fps
	physicsTick := time.NewTicker(time.Millisecond * 2)
	timeStep := (time.Millisecond * 2).Seconds()

	moveChannel := make(chan *MoveRequest)
	addChannel := make(chan *AddRequest)

	// actual Game Loop. TODO: Should this be a function call?
	go func() {
		// Run the game loop forever.
		for range gameTick.C {

			// NOTE TO FUTURE SELF: if multiple channels are ready, select will
			// pick one randomly and move on!! There are a few solutions I can see
			// to help this. First, have a select for each channel or read the
			// channels outside of the game loop.
			select {
			case msg := <-addChannel:
				fmt.Printf("Added!!!!!! %+v\n", msg)
			default:
				// Move on to other things
			}
		}
	}()

	// Start phyics loop, give it the movement channel and it's ticker
	go PhysicsLoop(physicsTick, moveChannel, timeStep)

	fmt.Println("Started Game Loop Go Routine atlkj ")

	return moveChannel, addChannel
}

// Physics loops listens from move requests and
func PhysicsLoop(physicsTick *time.Ticker, moveChannel chan *MoveRequest, timeStep float64) {
	frameSimulated := 0
	for range physicsTick.C {
		// Read any movement updates
		select {
		case msg := <-moveChannel:
			fmt.Printf("Physics doing movement%+v\n", msg)
		default:
			// Move on to other things
		}

		TickPhysics(timeStep)
		// TODO: Send this to a channel after reading an event so we can listen
		// in and know exactly which tick the event was registered
		frameSimulated++
	}
}

// Ticks the physics engine once by time elapsed
func TickPhysics(timeElapsed float64) {
	for _, gameObj := range physicsComponents {
		// Basic movement for now
		gameObj.X += gameObj.XVel * timeElapsed
		gameObj.Y += gameObj.YVel * timeElapsed
	}
}
