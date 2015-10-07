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

type Test struct {
	Event string
	Data  json.RawMessage
}

type Tester struct {
	Blah []Test
}

func broadCastGameObjects() {
	// run the sync event first to ensure objects are created by the client first
	// TODO: This should be batched with with the update message
	var broadcastMessages []Message

	addMessages := readAddRequests(channelCoordinator.broadcastAddChannel)

	if addMessages != nil {
		// broadcastMessages = append(broadcastMessages, syncEvent.Objects)
		addBytes, _ := json.Marshal(addMessages)
		broadcastMessages = appendEventMessage("add", addBytes, broadcastMessages)
		// Trace.Printf("Broadcasting some shit %s", string(syncBytes))
		// clientBackend.BroadCastPackets(syncBytes, nil)

		// // ==================Test=================
		// // blahBlah := json.RawMessage(syncBytes)
		// test := Test{Event: "test1", Data: syncBytes}
		// test2 := Test{Event: "test2", Data: syncBytes}
		// testSlice := make([]Test, 2)
		// testSlice[0] = test
		// testSlice[1] = test2
		// tester := Tester{Blah: testSlice}
		// testEvent, _ := json.Marshal(tester)
		// Trace.Printf("Broadcasting some shit %s", string(testEvent))
		// clientBackend.BroadCastPackets(testEvent, nil)

	}

	updateMessage := buildUpdateEvent(gameObjects)
	broadcastMessages = appendEventMessage("update", updateMessage, broadcastMessages)
	// broadcastMessages = append(broadcastMessages, updateMessage)

	// TODO: Have multiple events being sent back to the client

	// updateBytes, _ := json.Marshal(updateEvent)
	broadcastEventsMessages := makeEventsMessage(broadcastMessages)
	broadcastBytes, _ := json.Marshal(broadcastEventsMessages)
	// Broadcast any added game object
	Trace.Println("broadcasting this data\n " + string(broadcastBytes))
	clientBackend.BroadCastPackets(broadcastBytes, nil)
}

// Build the update event from the game objects
// TODO: Channel?
// TODO: I believes this copies the entire game object map. That may be
// desirable as it means it just a quick snap shot of the game objects
func buildUpdateEvent(gameObjects map[string]GameObject) []byte {
	updateData := make([]UpdateMessage, 0)

	for _, gameObj := range gameObjects {
		jsonData := gameObj.BuildUpdateMessage()
		updateData = append(updateData, jsonData)
	}

	updateBytes, _ := json.Marshal(updateData)

	// updateMessage := makeSingleEventMessage("update", updateBytes)

	return updateBytes
}

func makeSingleEventMessage(eventName string, eventData []byte) Message {
	message := Message{Event: eventName, Data: eventData}
	return message
}

func makeEventsMessage(events []Message) Events {
	eventsMessage := Events{Events: events}
	return eventsMessage
}

func appendEventMessage(event string, eventData []byte, currentMessageArr []Message) []Message {
	eventMessage := makeSingleEventMessage(event, eventData)
	appendedSlice := append(currentMessageArr, eventMessage)
	return appendedSlice
}

// Read the add request sent to the broadcast channel
func readAddRequests(broadCastAddChannel chan *AddRequest) []SyncMessage {
	var addMessages []SyncMessage
	for i := 0; i < 10; i++ {
		// Arbitraily read up to ten add requests in a single frame
		select {
		case msg := <-broadCastAddChannel:
			// TODO: add some sort of types
			syncMessage := SyncMessage{ObjType: "blah",
				Id: msg.Id, X: msg.X, Y: msg.Y}

			addMessages = append(addMessages, syncMessage)
			Trace.Printf("Adding with  %+v\n", syncMessage)
		default:
			// Move on to other things
			break
		}
	}
	return addMessages
}

// associates the object to the sourceId
func AddObjectToConnectionData(object GameObject, objectId string, sourceId int) {
	Trace.Println("Adding to client data")
	clientData := clientIdMap[sourceId]
	clientData.GameObjects[objectId] = object
}
