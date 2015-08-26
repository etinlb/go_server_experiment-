package main

// TODO: What is this file and why does it exist?

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	Event string
	Data  json.RawMessage // how data is parsed depends on the event
}

// keeps track of data from a client
type ClientData struct {
	Client      *websocket.Conn
	GameObjects map[string]GameObject
}

type ClientMessage struct {
	Id   string `json:"id"`
	Data json.RawMessage
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

// send all game objects that are currently in the game object map to the
// client connected
func SyncClient(client *websocket.Conn) {
	// TODO: Assess whether or not this is going to be to slow
	syncData := make([]GameObject, 0)

	for conn, connData := range clients {
		if conn == client {
			continue
		}

		for _, obj := range connData.GameObjects {
			syncData = append(syncData, obj)
		}
	}

	syncMessage := SyncMessage{Event: "sync", Objects: syncData}
	syncMessageAsJson, _ := json.Marshal(syncMessage)

	clientBackend.SendToClient(syncMessageAsJson, client)
}

func BuildObjectPackage(event string, gameObj GameObject) []byte {
	updateMessage := ObjectMessage{Event: event, Packet: gameObj}
	message, _ := json.Marshal(updateMessage)

	return message
}
