package main

import (
	// "flag"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

const masterUrl = "http://localhost:4000/jackIn"

// various object maps to keep track of different types of objects
var gameObjects map[string]GameObject
var playerObjects map[string]*Player
var physicsComponents map[string]*PhysicsComponent

// Communication coordinator
var channelCoordinator ComunicationChannels

var connections map[*websocket.Conn]bool

// map that keeps track of what data came from what client
var clients map[*websocket.Conn]ClientData

var clientBackend backend.BackendController
var serverBackend backend.BackendController

var neighbors NeighborServerList

type NeighborServer struct {
	Port int `json:"port"`
	Ip   string
}

type NeighborServerList struct {
	Servers []NeighborServer `json:"servers"`
}

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

func initializeGameData() {
	// keyed by id
	gameObjects = make(map[string]GameObject)
	playerObjects = make(map[string]*Player)
	physicsComponents = make(map[string]*PhysicsComponent)
}

// TODO: SHould this be in server vars?
func initializeConnectionData() {
	clients = make(map[*websocket.Conn]ClientData)
	connections = make(map[*websocket.Conn]bool)
}

func main() {
	port := flag.Int("port", 8080, "port to serve on")
	// TODO: have this address passed from the other server
	dir := flag.String("directory", "../web/", "directory of web files")
	flag.Parse()

	// =========Game Initializations============
	initializeGameData()

	// =========Connection Initializations============
	initializeConnectionData()

	//=========Backend Initializations============
	initializeServerVars()

	clientBackend = backend.NewBackendController(HandleClientEvent,
		cleanUpSocket,
		initializeClientData)

	// =========Connection handlers===================
	serverBackend = backend.NewBackendController(HandleServerEvent,
		cleanUpSocket,
		initializeServerData) // why doesn't go format align the arguments!

	neighbors = jackIn(*port)
	setUpNeighborConnections(neighbors)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", clientBackend.WsHandler)
	// the socket to read incoming connections from the master server
	http.HandleFunc("/masterSocket", serverBackend.WsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	moveChannel, addChannel := GameLoop()

	// Add channels to the channel coordinator
	channelCoordinator = ComunicationChannels{moveChannel: moveChannel, addChannel: addChannel}
	// fmt.Println(channelCoordinator)

	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

// register with the master server and get a list of neighbors to start connections with
func jackIn(port int) NeighborServerList {
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

	var neighbors NeighborServerList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&neighbors)

	if err != nil {
		fmt.Print("%v\n", err)
		panic(err)
	}

	fmt.Printf("%+v\n", neighbors)

	return neighbors
}

func setUpNeighborConnections(neighbors NeighborServerList) {
	for _, neighbor := range neighbors.Servers {
		url := "http://" + neighbor.Ip + ":" + strconv.Itoa(neighbor.Port) + "/masterSocket"
		serverBackend.NewWebsocket(url)
	}
}

func printGameObjectMap() {
	for _, obj := range gameObjects {
		log.Println(obj)
	}
}
