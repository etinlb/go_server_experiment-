package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Message struct {
	Event       string          `json:"event"`
	Data        json.RawMessage `json:"data"` // how data is parsed depends on the event
	SequenceNum int             `json:"seq_num"`
}

type Events struct {
	Events []Message `json:"events"`
}

// keeps track of data from a client
type ClientData struct {
	Client               *websocket.Conn
	GameObjects          map[string]GameObject
	ClientId             int
	CurrentSequnceNumber int
}

type ClientMessage struct {
	Id   string `json:"id"`
	Data json.RawMessage
}

// send all game objects that are currently in the game object map to the
// client connected
func SyncClient(client *websocket.Conn) {
	// TODO: Assess whether or not this is going to be to slow
	syncData := make([]AddMessage, 0)
		client.RemoteAddr(), client.LocalAddr())

	for _, obj := range gameObjects {
		syncData = append(syncData, obj.BuildAddMessage())
	}

	var broadcastMessages []Message
	syncDataBytes, _ := json.Marshal(syncData)
	// message := Message{Event: "df", Data: syncDataBytes}
	broadcastMessages = appendEventMessage("blahblah", syncDataBytes, broadcastMessages)

	broadcastEventsMessages := makeEventsMessage(broadcastMessages)
	broadcastBytes, _ := json.Marshal(broadcastEventsMessages)
	broadcastBytes2, _ := json.Marshal(broadcastMessages)

	message2 := makeSingleEventMessage("blahblah", syncDataBytes)

	messageJson, _ := json.Marshal(message2)

	clientBackend.SendToClient(messageJson, client)
}
