// File that holds all the server to server handler functions
package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"

	// "github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

// TODO: Learn better go data structures
var ServerConnections map[*websocket.Conn]ServerConnection

type ServerConnection struct {
	// an object that represents a conection with another game server
	// has where the server is and the connection to the server
	Connection *websocket.Conn
}

func HandleServerEvent(event []byte, client *websocket.Conn) {
	var message Message
	json.Unmarshal(event, &message)

	if message.Event == "coordinationEvent" {
		Trace.Println("coordination event")
		Trace.Println(message.Data)
	}
}

// Initializes data that this server will store about other connected server
func initializeServerData(conn *websocket.Conn) {
	server := ServerConnection{Connection: conn}
	ServerConnections[conn] = server
}

func initializeServerVars() {
	ServerConnections = make(map[*websocket.Conn]ServerConnection)
}

func newClientConnection(connectionUrl string) (conn *websocket.Conn) {
	u, err := url.Parse(connectionUrl)
	if err != nil {
		Error.Println(err)
	}

	Trace.Printf("url for new client connection is %+v\n", u)

	Info.Println("Attempting to dial " + u.Host)

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		Error.Printf("Failed to dial %s. Error: %+v\n", err)
		return nil
	}

	wsHeaders := http.Header{
		"Origin": {u.Host},
		// your milage may differ
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)

	if err != nil {
		Error.Printf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
	}

	Info.Printf("Successfully connected to %+v", wsConn.RemoteAddr())

	return wsConn
}
