// Holds all the client handler functiosn
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

// Client events is data sent from the client to the server
func HandleClientEvent(event []byte, client *websocket.Conn) {
	var message Message
	json.Unmarshal(event, &message)

	fmt.Println("Handling client data")

	channelCoordinator.ProcessEvents(message.Event, message.Data)
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
