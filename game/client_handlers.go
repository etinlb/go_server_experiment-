// Holds all the client handler functiosn
package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"math/rand"
)

// Client events is data sent from the client to the server
func HandleClientEvent(event []byte, client *websocket.Conn) {
	var message Message
	json.Unmarshal(event, &message)
	sourceId := connections[client]
	channelCoordinator.ProcessEvents(message.Event, message.Data, sourceId)
}

func initializeClientData(conn *websocket.Conn) {
	// initialize the connection
	// connections[conn] = true
	clientData := ClientData{Client: conn, GameObjects: make(map[string]GameObject)}
	// clients[conn] =
	clientId := AddClientDataToMap(clientIdMap, &clientData)
	clientData.ClientId = clientId
	connections[conn] = clientId

	SyncClient(conn)
}

func ExcludeClient(client *websocket.Conn) map[*websocket.Conn]bool {
	// makes a map with only this one client to pass to sendPackets
	connections := make(map[*websocket.Conn]bool)
	connections[client] = true

	return connections
}

func AddClientDataToMap(mapToAdd map[int]*ClientData, clientToAdd *ClientData) int {
	x := rand.Int()
	for {
		if _, ok := mapToAdd[x]; !ok {
			mapToAdd[x] = clientToAdd
			return x
		}
		x = rand.Int()
	}
}
