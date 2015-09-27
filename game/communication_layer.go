// Objects for coordinating the communications of channesl
package main

import (
	"encoding/json"
)

// Various channel structs that are used for communicating with the game and physics loop
type MoveRequest struct {
	Xvel float64
	Yvel float64
	BaseGameObjData
}

// Adds an object with the id the game object map. Same as update message, just
// named more clearly
type AddRequest struct {
	// UpdateMessage
	X float64
	Y float64
	BaseGameObjData
}

// Global struct I think?
// moveChannel         - client move request
// addChannel          - client add request
// broadcastAddChannel - add request that gets broadcasts to the other clients
// serverAddChannel    - Channel used to communicate to the data structure keeping track of which objects belong to which source
type ComunicationChannels struct {
	moveChannel         chan *MoveRequest
	addChannel          chan *AddRequest
	broadcastAddChannel chan *AddRequest
	serverAddChannel    chan *BaseGameObjData
}

// Takes the event bytes from handle client event and processes them
func (c *ComunicationChannels) ProcessEvents(event string, data []byte, eventSourceId int) {
	Trace.Printf("Received %s event with this data: %s\n", event, string(data))

	if event == "createPlayer" {
		addReq := ReadCreatePlayerEvent(data)
		// TODO: Add the channel id to the
		addReq.sourceId = eventSourceId
		c.addChannel <- addReq

	} else if event == "move" {
		moveReq := ReadMoveEvent(data)
		moveReq.sourceId = eventSourceId
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

	// TODO: Have multiple events being sent back to the client
	updateEvent := UpdateEvent{"update", updateData}

	updateBytes, _ := json.Marshal(updateEvent)

	// Broadcast any added game object
	clientBackend.BroadCastPackets(updateBytes, nil)

	// TODO: This should be batched with with the update message
	syncEvent := readBroadCastEvents(channelCoordinator.broadcastAddChannel)
	if syncEvent.Objects != nil {
		syncBytes, _ := json.Marshal(syncEvent)
		Trace.Printf("Broadcasting some shit %s", string(syncBytes))
		clientBackend.BroadCastPackets(syncBytes, nil)
	}

}

func readBroadCastEvents(broadCastAddChannel chan *AddRequest) SyncEvent {
	syncEvent := SyncEvent{Event: "createObject"}
	for i := 0; i < 10; i++ {
		// Arbitraily read up to ten add requests in a single frame
		select {
		case msg := <-broadCastAddChannel:
			// TODO: add some sort of types
			syncMessage := SyncMessage{ObjType: "blah",
				Id: msg.Id, X: msg.X, Y: msg.Y}

			syncEvent.Objects = append(syncEvent.Objects, syncMessage)
			Trace.Printf("Adding with  %+v\n", syncEvent)
		default:
			// Move on to other things
		}
	}
	return syncEvent
}

// associates the object to the sourceId
func AddObjectToConnectionData(object GameObject, objectId string, sourceId int) {
	Trace.Println("Adding to client data")
	clientData := clientIdMap[sourceId]
	clientData.GameObjects[objectId] = object
}
