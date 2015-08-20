package main

// import "fmt"
import (
	"encoding/json"
)

var gameId = 1
var playerMovementXVel = 100.0
var playerMovementYVel = 100.0

// Game Object is a struct with various components, components themselves
// aren't game objects
type GameObject interface {
	Update()
	ReadMessage()    // process data it gets from the client
	PackageMessage() // process data it gets from the client
}

// TODO: Change to physics object
// type PhysicalObject interface {
// 	AddImpulse(xVel, yVel float64)
// 	CurrentVelocity() (int, int)
// }

// I'm not sure what I want for this yet
// type PlayerControlledObject interface {
// 	Move(xVel, yVel float64)
// }

type RectComponent struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type BaseGameObjData struct {
	Id string `json:"id"`
}

type PhysicsComponent struct {
	RectComponent
	XVel   float64 `json:"xVel"`
	YVel   float64 `json:"yVel"`
	XForce float64
	YForce float64
}

type Player struct {
	PhysicsComp *PhysicsComponent //
	Id          string            // the identifier of the client controlling this object
}

func (m *PhysicsComponent) Move(xAxis, yAxis float64) {
	m.XVel += playerMovementXVel * xAxis
	m.YVel += playerMovementYVel * yAxis
}

func (m *Player) Update() {

}

func (m *Player) ReadMessage() {

}

func (m *Player) PackageMessage() {

}

func MakePlayerObjectFromJson(data json.RawMessage) Player {
	dataMessage := ReadCreateMessage(data)
	newPlayer := NewPlayer(dataMessage.X, dataMessage.Y, dataMessage.Id)
	return newPlayer
}

func (m *PhysicsComponent) AddImpulse(xForce, yForce float64) {
	m.XForce += xForce
	m.YForce += yForce
	// m.XVel = xVel
	// m.YVel = yVel
	// fmt.Printf("%+v After Moving\n", m)
}

func NewPhysicsComponent(x, y float64) PhysicsComponent {
	rect := NewRectComponent(x, y)
	gameObject := PhysicsComponent{
		RectComponent: rect,
		XVel:          0,
		YVel:          0,
	}

	return gameObject
}

func NewPlayer(x, y float64, id string) Player {
	physicsComponenet := NewPhysicsComponent(x, y)
	playerObject := Player{
		PhysicsComp: &physicsComponenet,
		Id:          id,
	}
	return playerObject
}

func AddPlayerObject() {

}

func NewRectComponent(x, y float64) RectComponent {
	rect := RectComponent{X: x, Y: y}

	return rect
}
