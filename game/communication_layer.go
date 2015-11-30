// Objects for coordinating the communications of channesl
package main

import (
	"encoding/json"
)

// Base Data needed for update, adds and move requests.
// Meant to be embedded in the various requests
type BaseRectMessage struct {
	X float64
	Y float64
	BaseGameObjData
}

func buildBaseRectData(x, y float64, id string) BaseRectMessage {
	objectRect := BaseRectMessage{X: x, Y: y, BaseGameObjData: BaseGameObjData{Id: id}}
	return objectRect
}

// Various channel structs that are used for communicating with the game and physics loop
type MoveMessage struct {
	BaseRectMessage
}

// Adds an object with the id the game object map. Same as update message, just
// named more clearly
type AddMessage struct {
	BaseRectMessage
	ObjType string `json: "type"`
}

func buildAddMessage(base BaseRectMessage, objType string) AddMessage {
	return AddMessage{BaseRectMessage: base, ObjType: objType}
}

type UpdateMessage struct {
	BaseRectMessage
}

// Global struct I think?
// moveChannel         - client move request
// addChannel          - client add request
// broadcastAddChannel - add request that gets broadcasts to the other clients
// serverAddChannel    - Channel used to communicate to the data structure keeping track of which
// 						 objects belong to which source
type ComunicationChannels struct {
	moveChannel         chan *MoveMessage
	addChannel          chan *AddMessage
	broadcastAddChannel chan *AddMessage
	serverAddChannel    chan *BaseGameObjData
}

// Takes the event bytes from handle client event and processes them
func (c *ComunicationChannels) ProcessEvents(event string, data []byte, eventSourceId int) {
	Trace.Printf("Received %s event with this data: %s\n", event, string(data))

	if event == "createPlayer" {
		addReq := ReadCreatePlayerEvent(data)
		addReq.sourceId = eventSourceId
		c.addChannel <- addReq

	} else if event == "move" {
		moveReq := ReadMoveEvent(data)
		moveReq.sourceId = eventSourceId
		c.moveChannel <- moveReq
	}
}

func ReadCreatePlayerEvent(data json.RawMessage) *AddMessage {
	var dataMessage *AddMessage
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}

func ReadMoveEvent(data json.RawMessage) *MoveMessage {
	var dataMessage *MoveMessage
	json.Unmarshal(data, &dataMessage)

	return dataMessage
}

func broadCastGameObjects() {
	// run the sync event first to ensure objects are created by the client first
	// TODO: This should be batched with with the update message
	var broadcastMessages []Message

	addMessages := readAddRequests(channelCoordinator.broadcastAddChannel)

	if addMessages != nil {
		addBytes, _ := json.Marshal(addMessages)
		broadcastMessages = appendEventMessage("add", addBytes, broadcastMessages)
	}

	// if len(gameObjects) == 0 {

	// }
	updateMessage := buildUpdateEvent(gameObjects)
	broadcastMessages = appendEventMessage("update", updateMessage, broadcastMessages)

	broadcastEventsMessages := makeEventsMessage(broadcastMessages)
	broadcastBytes, _ := json.Marshal(broadcastEventsMessages)

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
func readAddRequests(broadCastAddChannel chan *AddMessage) []AddMessage {
	var addMessages []AddMessage
	for i := 0; i < 10; i++ {
		// Arbitraily read up to ten add requests in a single frame
		select {
		case msg := <-broadCastAddChannel:
			// TODO: add some sort of types
			Trace.Printf("Adding with  %+v\n", msg)
			addMessages = append(addMessages, *msg)
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
