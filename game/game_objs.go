package main

// import (
// 	"encoding/json"
// 	"fmt"
// )

var gameId = 1
var playerMovementXVel = 1000.0
var playerMovementYVel = 1000.0

// Game Object is a struct with various components, components themselves
// aren't game objects
type GameObject interface {
	Update()
	ReadMessage() // process data it gets from the client
	BuildSyncMessage() SyncMessage
	BuildUpdateMessage() UpdateMessage // process data it gets from the client
}

type RectComponent struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type BaseGameObjData struct {
	Id string `json:"id"`
}

type UpdateMessage struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Id string  `json:"id"`
}

type UpdateEvent struct {
	Event   string          `json:"event"`
	Objects []UpdateMessage `json:"data"`
}

type PhysicsComponent struct {
	RectComponent
	XVel   float64 `json:"xVel"`
	YVel   float64 `json:"yVel"`
	XForce float64
	YForce float64
	Location Vector2D `json:"location"`
	Velocity Vector2D `json:"velocity"`
	Force    Vector2D
}

type Player struct {
	PhysicsComp *PhysicsComponent //
	Id          string            // the identifier of the client controlling this object
}

func (m *PhysicsComponent) Move(xAxis, yAxis float64) {
	m.XVel += playerMovementXVel * xAxis
	m.YVel += playerMovementYVel * yAxis
	m.Velocity.X += playerMovementXVel * xAxis
	m.Velocity.Y += playerMovementYVel * yAxis
}

func (m *Player) BuildSyncMessage() SyncMessage {
	message := SyncMessage{"player", m.Id}
	message := SyncMessage{"player", m.Id, m.PhysicsComp.Location.X, m.PhysicsComp.Location.Y}
	return message
}

func (m *Player) Update() {

}

func (m *Player) ReadMessage() {

}

// Packages the player Physics state into a json byte array
func (m *Player) BuildUpdateMessage() UpdateMessage {
	updateMessage := UpdateMessage{m.PhysicsComp.X, m.PhysicsComp.Y, m.Id}

	return updateMessage
}

func (m *PhysicsComponent) AddImpulse(xForce, yForce float64) {
	m.XForce += xForce
	m.YForce += yForce
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

func NewRectComponent(x, y float64) RectComponent {
	rect := RectComponent{X: x, Y: y}

	return rect
}
