package main

import (
	// "flag"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

const masterUrl = "http://localhost:4000/jackIn"

var gameObjects map[string]*GameObject

var connections map[*websocket.Conn]bool

// map that keeps track of what data came from what client
var clients map[*websocket.Conn]ClientData

var clientBackend backend.BackendController

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
	dir := flag.String("directory", "../web/", "directory of web files")
	flag.Parse()

	// XXX
	jackIn(*port)

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

	clientBackend = backend.NewBackendController(HandleClientEvent,
		cleanUpSocket,
		initializeClientData)

	// =========Connection handlers===================
	serverBackend = backend.NewBackendController(HandleServerEvent,
		cleanUpSocket,
		initializeServerData) // why doesn't go format align the arguments!

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", clientBackend.WsHandler)
	// the socket to read incoming connections from the master server
	http.HandleFunc("/masterSocket", serverBackend.WsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

// register with the master server and get a list of neighbors to start connections with
func jackIn(port int) {
	jsonStr := "{\"port\":" + strconv.Itoa(port) + "}"
	log.Println(jsonStr)
	var jsonByte = []byte(jsonStr)
	req, err := http.NewRequest("POST", masterUrl, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Print("%v", resp.Body)
	fmt.Print("%v", err)
}

func setUpNeighborConnections() {

}

func printGameObjectMap() {
	for _, obj := range gameObjects {
		log.Println(*obj)
	}
}
