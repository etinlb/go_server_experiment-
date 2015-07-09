// File that holds all the server to server handler functions
package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	// TODO: Fix this so it's not just the general interface object
	var message Message
	json.Unmarshal(event, &message)

	if message.Event == "coordinationEvent" {

		log.Println("coordination event")
		log.Println(message.Data)
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

// log.Println(portAsString)
// connectionUrl := "http://localhost:" + portAsString + "/masterSocket"
// x := newClientConnection(connectionUrl)
// log.Println(x)
func newClientConnection(connectionUrl string) (conn *websocket.Conn) {
	u, err := url.Parse(connectionUrl)
	if err != nil {
		log.Println(err)
	}
	log.Println(u)

	log.Println(u.Host)
	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Println(err)
		return nil
	}

	wsHeaders := http.Header{
		"Origin": {u.Host},
		// your milage may differ
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)

	if err != nil {
		fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
	}
	return wsConn
}
