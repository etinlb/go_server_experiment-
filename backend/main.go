// backend node for managing connections
package backend

import (
	// "flag"
	// "flag"
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type cleanUpFunction func(*websocket.Conn)
type eventHandlerFunction func([]byte, *websocket.Conn)
type connectionHandlerFunction func(*websocket.Conn)

var cleanUpHandler cleanUpFunction
var eventHandler eventHandlerFunction
var conectionHandler connectionHandlerFunction

var connections map[*websocket.Conn]bool

func BroadCastPackets(msg []byte, connections map[*websocket.Conn]bool, excludeList map[*websocket.Conn]bool) {
	for conn := range connections {
		if _, ok := excludeList[conn]; ok {
			continue
		}

		SendToClient(msg, conn)
	}
}

func WsHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
	log.Println("getting a connection")

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(writer, "got a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	conectionHandler(conn)
	defer cleanUpHandler(conn)      // if this function ever exits, clean up the data
	defer delete(connections, conn) // if this function ever exits, clean up the data

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		eventHandler(msg, conn)
	}
}

func SendToClient(msg []byte, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("Deleting")
		cleanUpHandler(conn)
	}
}

func InitHandlerFunctions(event eventHandlerFunction, cleanUp cleanUpFunction, connections connectionHandlerFunction) {
	cleanUpHandler = cleanUp
	eventHandler = event
	conectionHandler = connections
}
