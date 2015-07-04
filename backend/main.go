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

// map that keeps track of what data came from what client
var clients map[*websocket.Conn]ClientData

func broadCastPackets(msg []byte, excludeList map[*websocket.Conn]bool) {
	for conn := range connections {
		if _, ok := excludeList[conn]; ok {
			continue
		}

		sendToClient(msg, conn)
	}
}

func sendToClient(msg []byte, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("Deleting")
		cleanUpSocket(conn)
	}
}

func cleanUpSocket(conn *websocket.Conn) {
	log.Println("Cleaning up")
	fmt.Println(clients[conn])

	for id, _ := range clients[conn].GameObjects {
		log.Println("deleting from gameObjects map")
		delete(gameObjects, id)
	}

	delete(clients, conn)
	delete(connections, conn)

	conn.Close()
	printGameObjectMap()
}

func initializeConectionVaribles(conn *websocket.Conn) {
	// initialize the connection
	connections[conn] = true
	clients[conn] = ClientData{Client: conn, GameObjects: make(map[string]*GameObject)}
	SyncClient(conn)
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

	initializeConectionVaribles(conn)
	defer cleanUpSocket(conn) // if this function ever exits, clean up the data

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
	clients = make(map[*websocket.Conn]ClientData)

	// =========Connection Initializations============
	connections = make(map[*websocket.Conn]bool)

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

func printGameObjectMap() {
	for _, obj := range gameObjects {
		log.Println(*obj)
	}

}
