// backend node for managing connections
package backend

import (
	// "flag"
	// "flag"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type cleanUpFunction func(*websocket.Conn)
type eventHandlerFunction func([]byte, *websocket.Conn)
type connectionHandlerFunction func(*websocket.Conn)

// a backend controller abstracts handling and managing websocket connections
type BackendController struct {
	EventHandler      eventHandlerFunction
	CleanUpHandler    cleanUpFunction
	ConnectionHandler connectionHandlerFunction

	connections map[*websocket.Conn]bool
}

func NewBackendController(event eventHandlerFunction, cleanUp cleanUpFunction,

	connections connectionHandlerFunction) BackendController {

	controller := BackendController{EventHandler: event, CleanUpHandler: cleanUp,
		ConnectionHandler: connections}

	controller.connections = make(map[*websocket.Conn]bool)

	return controller
}

func (b BackendController) BroadCastPackets(msg []byte, excludeList map[*websocket.Conn]bool) {

	for conn := range b.connections {
		if _, ok := excludeList[conn]; ok {
			continue
		}

		b.SendToClient(msg, conn)
	}
}

func (b BackendController) WsHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
	log.Println("getting a connection")
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(writer, "got a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	b.ConnectionHandler(conn)
	defer b.CleanUpHandler(conn)      // if this function ever exits, clean up the data
	defer delete(b.connections, conn) // if this function ever exits, clean up the data

	b.connections[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		b.EventHandler(msg, conn)
	}
}

func (b BackendController) SendToClient(msg []byte, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("Deleting")
		b.CleanUpHandler(conn)
	}
}

// adds a connection to the connection map.
func (b BackendController) AddNewConnection(conn *websocket.Conn) {
	b.connections[conn] = true
}

func (b BackendController) NewWebsocket(connectionUrl string) (*websocket.Conn, error) {
	u, err := url.Parse(connectionUrl)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Cannot parse connection url" + connectionUrl)
	}
	log.Println(u)

	log.Println(u.Host)
	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot dial " + u.Host)
	}

	wsHeaders := http.Header{
		"Origin": {u.Host},
		// your milage may differ
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)

	if err != nil {

		return nil, fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)

	}

	b.AddNewConnection(wsConn)

	return wsConn, nil
}
