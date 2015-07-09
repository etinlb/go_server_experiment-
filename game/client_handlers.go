// Holds all the client handler functiosn
package main

import (
	"encoding/json"
	// "fmt"
	"log"

	// "github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

// Client events is data sent from the client to the server
func HandleClientEvent(event []byte, client *websocket.Conn) {
	// TODO: Fix this so it's not just the general interface object
	var message Message
	json.Unmarshal(event, &message)

	if message.Event == "createUnit" {
		log.Println("Creating a new object from " + string(message.Data))
		newGameObj := MakeObjectFromJson(message.Data)

		gameObjects[newGameObj.Id] = &newGameObj
		clients[client].GameObjects[newGameObj.Id] = &newGameObj

		packet := BuildObjectPackage("createUnit", &newGameObj)
		clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	} else if message.Event == "update" {
		updateData := ReadCreateMessage(message.Data)

		gameObjects[updateData.Id].Rect.Y = updateData.Y
		gameObjects[updateData.Id].Rect.X = updateData.X

		packet := BuildObjectPackage("update", gameObjects[updateData.Id])
		clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	} else if message.Event == "coordinationEvent" {

		log.Println("coordination event")
		log.Println(message.Data)
	}
}

func initializeClientData(conn *websocket.Conn) {
	// initialize the connection
	connections[conn] = true
	clients[conn] = ClientData{Client: conn, GameObjects: make(map[string]*GameObject)}
	SyncClient(conn)
}

func ExcludeClient(client *websocket.Conn) map[*websocket.Conn]bool {
	// makes a map with only this one client to pass to sendPackets
	connections := make(map[*websocket.Conn]bool)
	connections[client] = true

	return connections
}
