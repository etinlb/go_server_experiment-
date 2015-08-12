// Holds all the client handler functiosn
package main

import (
	"encoding/json"
	"fmt"
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
		log.Println(newGameObj.X)

		gameObjects[newGameObj.Id] = &newGameObj
		clients[client].GameObjects[newGameObj.Id] = &newGameObj

		packet := BuildObjectPackage("createUnit", &newGameObj)
		clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	} else if message.Event == "update" {
		updateData := ReadCreateMessage(message.Data)
		log.Println("In update loop")
		// gameObjects[updateData.Id].Update()

		packet := BuildObjectPackage("update", gameObjects[updateData.Id])
		clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	} else if message.Event == "move" {
		// The client requested moving
		log.Println("moving with this packet")
		log.Println(string(message.Data))

		updateData := ReadMoveMessage(message.Data)
		fmt.Printf("%+v Reed this message\n", updateData)

		if mover, ok := gameObjects[updateData.Id].(Mover); ok {

			mover.Move(updateData.XVel, updateData.YVel)
		}
	}
}

func initializeClientData(conn *websocket.Conn) {
	// initialize the connection
	connections[conn] = true
	clients[conn] = ClientData{Client: conn, GameObjects: make(map[string]GameObject)}
	SyncClient(conn)
}

func ExcludeClient(client *websocket.Conn) map[*websocket.Conn]bool {
	// makes a map with only this one client to pass to sendPackets
	connections := make(map[*websocket.Conn]bool)
	connections[client] = true

	return connections
}
