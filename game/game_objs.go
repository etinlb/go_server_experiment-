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

// Component for the bare minimum representation of a game object
// :id     - unique object id
// :source - source id of the client or server the object belongs to
type BaseGameObjData struct {
	Id       string `json:"id"`
	sourceId int    `json:"sourceId"`
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
	Location Vector2D `json:"location"`
	Velocity Vector2D `json:"velocity"`
	Force    Vector2D
}

type Player struct {
	PhysicsComp *PhysicsComponent //
	Id          string            // the identifier of the client controlling this object
}

func (m *PhysicsComponent) Move(xAxis, yAxis float64) {
	m.Velocity.X += playerMovementXVel * xAxis
	m.Velocity.Y += playerMovementYVel * yAxis
}

func (m *Player) BuildSyncMessage() SyncMessage {
	message := SyncMessage{"player", m.Id, m.PhysicsComp.Location.X, m.PhysicsComp.Location.Y}
	return message
}

func (m *Player) Update() {

}

func (m *Player) ReadMessage() {

}

// Packages the player Physics state into a json byte array
func (m *Player) BuildUpdateMessage() UpdateMessage {
	updateMessage := UpdateMessage{m.PhysicsComp.Location.X, m.PhysicsComp.Location.Y, m.Id}

	return updateMessage
}

// func (m *PhysicsComponent) AddImpulse(xForce, yForce float64) {
// 	m.XForce += xForce
// 	m.YForce += yForce
// }

func NewPhysicsComponent(x, y float64) PhysicsComponent {
	locationVector := NewVector(x, y)
	gameObject := PhysicsComponent{Location: locationVector}

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
