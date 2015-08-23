// Holds all the client handler functiosn
package main

import (
	"encoding/json"
	"fmt"
	// "log"

	// "github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

// Client events is data sent from the client to the server
func HandleClientEvent(event []byte, client *websocket.Conn) {
	// TODO: Fix this so it's not just the general interface object
	var message Message
	json.Unmarshal(event, &message)

	fmt.Println("Handling client data")

	channelCoordinator.ProcessEvents(message.Event, message.Data)
	// if message.Event == "createPlayer" {
	// 	log.Println("Creating a new object from " + string(message.Data))
	// 	newPlayer := MakePlayerObjectFromJson(message.Data)
	// 	// log.Println(newGameObj.X)

	// 	// TODO: Why are we keeping track of this twice?
	// 	gameObjects[newPlayer.Id] = &newPlayer
	// 	clients[client].GameObjects[newPlayer.Id] = &newPlayer
	// 	playerObjects[newPlayer.Id] = &newPlayer
	// 	physicsComponents[newPlayer.Id] = newPlayer.PhysicsComp

	// 	packet := BuildObjectPackage("createUnit", &newPlayer)
	// 	clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	// 	fmt.Println("B")
	// } else if message.Event == "move" {
	// 	// The client requested moving
	// 	log.Println("moving with this packet")
	// 	log.Println(string(message.Data))

	// 	updateData := ReadMMessage(message.Data)
	// 	fmt.Printf("%+v Read this message\n", updateData)
	// 	// gameObjects[updateData.Id].

	// 	// if mover, ok := gameObjects[updateData.Id].(Mover); ok {

	// 	// gameObjects[updateData.Id].(updateData.XVel, updateData.YVel)
	// 	// }
	// }
	// else if message.Event == "update" {
	// 	updateData := ReadCreateMessage(message.Data)
	// 	log.Println("In update loop")
	// 	// gameObjects[updateData.Id].Update()

	// 	packet := BuildObjectPackage("update", *gameObjects[updateData.Id])
	// 	clientBackend.BroadCastPackets(packet, ExcludeClient(client))
	// }
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
