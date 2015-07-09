package main

import (
	"encoding/json"
	// "fmt"
	// "log"

	// "github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

type Message struct {
	Event string
	Data  json.RawMessage // how data is parsed depends on the event
}

// keeps track of data from a client
type ClientData struct {
	Client      *websocket.Conn
	GameObjects map[string]*GameObject
}

type CreateMessage struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Id string `json:"id"`
}

type SyncMessage struct {
	Event   string       `json:"event"` // client works with lowercase
	Objects []GameObject `json:"data"`
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
	X int `json:"x"`
	Y int `json:"y"`
}

type GameObject struct {
	Rect Rect
	Id   string `json:"id"`
}

func MakeCreateMessage(obj GameObject) CreateMessage {
	message := CreateMessage{X: obj.Rect.X, Y: obj.Rect.Y, Id: obj.Id}

	return message
}

// send all game objects that are currently in the game object map to the
// client connected
func SyncClient(client *websocket.Conn) {
	syncData := make([]GameObject, 0) // TODO: Assess whether or not this is going to be to slow

	for conn, connData := range clients {
		if conn == client {
			continue
		}

		for _, obj := range connData.GameObjects {
			syncData = append(syncData, *obj)
		}
	}

	syncMessage := SyncMessage{Event: "sync", Objects: syncData}
	syncMessageAsJson, _ := json.Marshal(syncMessage)

	clientBackend.SendToClient(syncMessageAsJson, client)
}

func BuildObjectPackage(event string, gameObj *GameObject) []byte {
	updateMessage := ObjectMessage{Event: event, Packet: *gameObj}
	message, _ := json.Marshal(updateMessage)

	return message
}

func ReadCreateMessage(data json.RawMessage) CreateMessage {
	var dataMessage CreateMessage
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}

func MakeObjectFromJson(data json.RawMessage) GameObject {
	dataMessage := ReadCreateMessage(data)
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
