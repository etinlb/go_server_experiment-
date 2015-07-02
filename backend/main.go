package main

import (
	// "flag"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var gameObjects map[string]*GameObject

var connections map[*websocket.Conn]bool
var broadcastPackets [][]byte // array of packets to send to all connected players

func sendPackets(msg []byte, excludeList map[*websocket.Conn]bool) {
	for conn := range connections {

		if _, ok := excludeList[conn]; ok {
			log.Println("don't send!")
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("send!")

			delete(connections, conn)
			break
		}
	}
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
	log.Println("getting a connection")

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(writer, "got a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	connections[conn] = true

	defer conn.Close() // if this function ever exits, close the connection
	// var messages =
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		HandleEvent(msg, conn)
	}
}

func main() {
	port := flag.Int("port", 8080, "port to serve on")
	dir := flag.String("directory", "../web/", "directory of web files")
	flag.Parse()

	// =========Game Initializations============
	// keyed by id
	gameObjects = make(map[string]*GameObject)

	// =========Connection Initializations============
	connections = make(map[*websocket.Conn]bool)
	// broadcastPackets = make([][]byte)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
