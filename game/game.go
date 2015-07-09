package main

import (
	// "flag"
	"flag"
	"fmt"
	"log"
	// "strconv"
	"net/http"

	"github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

var gameObjects map[string]*GameObject

var connections map[*websocket.Conn]bool

// map that keeps track of what data came from what client
var clients map[*websocket.Conn]ClientData

var clientBackend backend.BackendController
var serverBackend backend.BackendController

func cleanUpSocket(conn *websocket.Conn) {
	for id, _ := range clients[conn].GameObjects {
		log.Println("deleting from gameObjects map")
		delete(gameObjects, id)
	}

	delete(clients, conn)
	delete(connections, conn)

	conn.Close()
	printGameObjectMap()
}

func main() {
	port := flag.Int("port", 8080, "port to serve on")
	// TODO: have this address passed from the other server
	// neighborPort := flag.Int("neighbor", 8081, "port to connect to neighbor on")
	dir := flag.String("directory", "../web/", "directory of web files")

	flag.Parse()

	// TODO: WTF IS THIS SHIT
	// portAsString := strconv.Itoa(*neighborPort)

	// =========Game Initializations============
	// keyed by id
	gameObjects = make(map[string]*GameObject)
	clients = make(map[*websocket.Conn]ClientData)

	// =========Connection Initializations============
	connections = make(map[*websocket.Conn]bool)
	initializeServerVars()

	//=========Backend Initializations============

	clientBackend = backend.NewBackendController(HandleClientEvent, cleanUpSocket, initializeClientData)
	// =========Connection handlers===================
	serverBackend = backend.NewBackendController(HandleServerEvent, cleanUpSocket, initializeServerData)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", clientBackend.WsHandler)
	// the socket to read incoming connections from the master server
	http.HandleFunc("/masterSocket", serverBackend.WsHandler)

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
