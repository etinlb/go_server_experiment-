// Objects for coordinating the communications of channesl
package main

import (
	"encoding/json"
	"fmt"
)

// Various channel structs that are used for communicating with the game and physics loop
type MoveRequest struct {
	Xvel     int
	Yvel     int
	PlayerId string
}

// Adds an object with the id the game object map
type AddRequest struct {
	X        int
	Y        int
	PlayerId string
}

// Global struct I think?
type ComunicationChannels struct {
	moveChannel chan *MoveRequest
	addChannel  chan *AddRequest
}

// Takes the event bytes from handle client event and processes them
func (c *ComunicationChannels) ProcessEvents(event string, data []byte) {
	fmt.Printf("%s, event with this data %s", event, string(data))
	fmt.Printf("%+v, writing to this channel ", c.addChannel)

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
