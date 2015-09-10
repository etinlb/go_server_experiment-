package main

import (
	"time"
)

// Spawns the game loop and returns the channels to comminucate with the game
// TODO: Currently that is just the move channels, maybe return the ticker channel?
// TODO: TODO: Make it return channel of channels
func StartGameLoop() (chan *MoveRequest, chan *AddRequest, chan *AddRequest) {
	// about 16 milliseconds for 60 fps a second
	gameTick := time.NewTicker(time.Millisecond * 10)

	// Physics runs at 50 fps
	physicsTick := time.NewTicker(time.Millisecond * 20)
	timeStep := (time.Millisecond * 2).Seconds()

	// TODO: Figure out buffering properly
	moveChannel := make(chan *MoveRequest, 10)
	addChannel := make(chan *AddRequest, 10)
	broadcastAddChannel := make(chan *AddRequest, 10)

	// actual Game Loop. TODO: Should this be a function call?
	go func() {
		// Run the game loop forever.
		for range gameTick.C {

			// NOTE TO FUTURE SELF: if multiple channels are ready, select will
			// pick one randomly and move on!! There are a few solutions I can see
			// to help this. First, have a select for each channel or read the
			// channels outside of the game loop.
			// TODO: this could be a function probably
			for i := 0; i < 10; i++ {
				// Arbitraily read up to ten add requests in a single frame
				select {
				case msg := <-addChannel:
					Trace.Printf("Adding with  %+v\n", msg)
					player := NewPlayer(msg.X, msg.Y, msg.Id)
					AddPlayerObjectToWorld(player)
					// TODO: Have proper error checking and only add to broadcast channel if
					// successful
					broadcastAddChannel <- msg
				default:
					// Move on to other things
				}
			}
			// TODO: Have this done with a channel I think...
			broadCastGameObjects()
		}
	}()

	// Start phyics loop, give it the movement channel and it's ticker
	go PhysicsLoop(physicsTick, moveChannel, timeStep)

	Info.Println("Started Game Loop")

	return moveChannel, addChannel, broadcastAddChannel
}

func AddPhysicsComp(comp *PhysicsComponent, id string) {
	physicsComponents[id] = comp
}

func AddPlayerObjectToWorld(player Player) {
	playerObjects[player.Id] = &player
	gameObjects[player.Id] = &player
	AddPhysicsComp(player.PhysicsComp, player.Id)
}

// Physics loops listens from move requests and
func PhysicsLoop(physicsTick *time.Ticker, moveChannel chan *MoveRequest, timeStep float64) {
	frameSimulated := 0
	for range physicsTick.C {
		// Read any movement updates
		select {
		// Right now, a move request only comes in through player movement
		case msg := <-moveChannel:
			id := msg.Id
			if physicsComp, ok := physicsComponents[id]; ok {
				//do something here
				physicsComp.Move(msg.Xvel, msg.Yvel)
			}
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
