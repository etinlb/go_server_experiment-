// Objects for coordinating the communications of channesl
package main

import (
	"encoding/json"
	"fmt"
)

// Global struct I think?
type ComunicationChannels struct {
	moveChannel chan MoveRequest
	addChannel  chan AddRequest
}

// Takes the event bytes from handle client event and processes them
func (c *ComunicationChannels) ProcessEvents(event string, data []byte) {
	// var message Message
	// json.Unmarshal(event, &message)

	fmt.Printf("%s, event with this data %s", event, string(data))
	fmt.Printf("%+v, writing to this channel ", c.addChannel)
	if event == "createPlayer" {
		addReq := ReadCreatePlayerEvent(data)
		fmt.Println("adding to channel")
		c.addChannel <- addReq
		fmt.Println("add to channel")
	}
}

func ReadCreatePlayerEvent(data json.RawMessage) AddRequest {
	var dataMessage AddRequest
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}
