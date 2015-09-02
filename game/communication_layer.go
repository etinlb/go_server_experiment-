// Objects for coordinating the communications of channesl
package main

import (
	"encoding/json"
	"log"
)

// Various channel structs that are used for communicating with the game and physics loop
type MoveRequest struct {
	Xvel float64
	Yvel float64
	Id   string
}

// Adds an object with the id the game object map. Same as update message, just
// named more clearly
type AddRequest struct {
	// UpdateMessage
	X  float64
	Y  float64
	Id string
}

// Global struct I think?
type ComunicationChannels struct {
	moveChannel         chan *MoveRequest // client move request
	addChannel          chan *AddRequest  // client add request
	broadcastAddChannel chan *AddRequest  // add request that gets broadcasts to the other clients
}

// Takes the event bytes from handle client event and processes them
func (c *ComunicationChannels) ProcessEvents(event string, data []byte) {
	log.Printf("%s event with this data: %s\n", event, string(data))

	if event == "createPlayer" {
		addReq := ReadCreatePlayerEvent(data)
		c.addChannel <- addReq
	} else if event == "move" {
		moveReq := ReadMoveEvent(data)
		c.moveChannel <- moveReq
	}
}

func ReadCreatePlayerEvent(data json.RawMessage) *AddRequest {
	var dataMessage *AddRequest
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}

func ReadMoveEvent(data json.RawMessage) *MoveRequest {
	var dataMessage *MoveRequest
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}

func broadCastGameObjects() {
	updateData := make([]UpdateMessage, 0)
	for _, gameObj := range gameObjects {
		jsonData := gameObj.BuildUpdateMessage()
		updateData = append(updateData, jsonData)
	}

	// TODO: Have muliple events being sent back to the client
	updateEvent := UpdateEvent{"update", updateData}

	updateBytes, _ := json.Marshal(updateEvent)

	clientBackend.BroadCastPackets(updateBytes, nil)

	// Broadcast any added game object

}
