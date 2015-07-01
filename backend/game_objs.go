package main

import (
	"encoding/json"
	// "fmt"
	"log"
)

type Message struct {
	Event string
	Data  json.RawMessage // how data is parsed depends on the event
}

type CreateMessage struct {
	X, Y int
	Id   string
}

// messages to send back to client...Can't be raw json?
// TODO: Figure out the struct stuff in go.
type ObjectMessage struct {
	Event  string     `json:"event"` // client works with lowercase
	Packet GameObject `json:"data"`
}

// TODO: Learn go better so these and the messages structs could be combined
// Might have to structure the json data begin sent differently
type Rect struct {
	X, Y int
}

type GameObject struct {
	Rect Rect
	Id   string
}

func HandleEvent(event []byte) {
	// TODO: Fix this so it's not just the general interface object
	var message Message
	log.Println(string(event))
	json.Unmarshal(event, &message)

	if message.Event == "createUnit" {
		newGameObj := MakeObjectFromJson(message.Data)
		gameObjects[newGameObj.Id] = &newGameObj
		log.Println(newGameObj.Rect)
	} else if message.Event == "update" {
		updateData := ReadCreateMessage(message.Data)
		gameObjects[updateData.Id].Rect.Y = updateData.Y
		gameObjects[updateData.Id].Rect.X = updateData.X
		BuildUpdatePackage(gameObjects[updateData.Id])
	}
}

func BuildUpdatePackage(gameObj *GameObject) {
	updateMessage := ObjectMessage{Event: "update", Packet: *gameObj}
	message, err := json.Marshal(updateMessage)
	log.Println(string(message))
	log.Println(err)
	SendUpdatePackage(message)
}

//TODO: Do some logic to not send data to the client that sent the update message?
func SendUpdatePackage(message []byte) {
	sendAll(message)
}

func ReadCreateMessage(data json.RawMessage) CreateMessage {
	var dataMessage CreateMessage
	json.Unmarshal(data, &dataMessage)
	return dataMessage
}

func MakeObjectFromJson(data json.RawMessage) GameObject {
	dataMessage := ReadCreateMessage(data)

	log.Println(string(data))
	log.Println(dataMessage)

	gameObject := NewGameObject(dataMessage.X, dataMessage.Y, dataMessage.Id)
	return gameObject
}

func NewGameObject(x, y int, id string) GameObject {
	rect := MakeRect(x, y)

	gameObject := GameObject{Rect: rect, Id: id}
	return gameObject
}

func MakeRect(x, y int) Rect {
	rect := Rect{X: x, Y: y}
	return rect
}
