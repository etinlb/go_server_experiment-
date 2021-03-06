package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
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

// Connection structures
var connections map[*websocket.Conn]int // Maps the connection object to the client id
var clientIdMap map[int]*ClientData

// map that keeps track of what data came from what client

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
	Info.Println("Cleaning up connection from %s", conn.RemoteAddr())
	Info.Println(len(gameObjects))
	clientId := connections[conn]
	clientData := clientIdMap[clientId]
	for id, _ := range clientData.GameObjects {
		Trace.Println("deleting from gameObjects map, id: %s", id)
		// TODO: Need to delete all references to this object...Idk how to best
		// do that.
		delete(gameObjects, id)
		delete(playerObjects, id)
		delete(physicsComponents, id)
	}

	delete(clientIdMap, clientData.ClientId)
	delete(connections, conn)
	Info.Println(len(gameObjects))

	conn.Close()
	printGameObjectMap()
}

func initializeGameData() {
	// keyed by id
	gameObjects = make(map[string]GameObject)
	playerObjects = make(map[string]*Player)
	physicsComponents = make(map[string]*PhysicsComponent)
}

func initializeLogger() {
	// TODO: Read a config file
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	// InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

// TODO: SHould this be in server vars?
func initializeConnectionData() {
	// TODO: Access if we need the clients variable
	// clients = make(map[*websocket.Conn]*ClientData)
	connections = make(map[*websocket.Conn]int)
	clientIdMap = make(map[int]*ClientData)
}

func main() {
	initializeLogger()

	port := flag.Int("port", 8080, "port to serve on")
	// TODO: have this address passed from the other server
	dir := flag.String("directory", "../web/", "directory of web files")
	interactive := flag.Bool("i", false, "Run with an interactive shell")
	flag.Parse()

	// =========Game Initializations============================================
	initializeGameData()

	// =========Connection Initializations======================================
	initializeConnectionData()

	//=========Backend Initializations==========================================
	initializeServerVars()

	clientBackend = backend.NewBackendController(HandleClientEvent,
		cleanUpSocket,
		initializeClientData)

	// =========Connection handlers=============================================
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

	Info.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	moveChannel, addChannel, broadcastAddChannel := StartGameLoop()

	// Add channels to the channel coordinator
	channelCoordinator = ComunicationChannels{
		moveChannel:         moveChannel,
		addChannel:          addChannel,
		broadcastAddChannel: broadcastAddChannel}

	if *interactive {
		go runHttpServer(addr)
		runInteractiveMode(channelCoordinator)
	} else {
		// function runs until an exit is called
		runHttpServer(addr)
	}
}

func runHttpServer(addr string) {
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	Warning.Println(err.Error())
}

// register with the master server and get a list of neighbors to start connections with
func jackIn(port int) NeighborServerList {
	jsonStr := "{\"port\":" + strconv.Itoa(port) + "}"
	Trace.Println("Jacking in with " + jsonStr)
	var jsonByte = []byte(jsonStr)

	req, err := http.NewRequest("POST", masterUrl, bytes.NewBuffer(jsonByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Error.Printf("%v\n", err)
		panic(err)
	}
	defer resp.Body.Close()

	var neighbors NeighborServerList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&neighbors)

	if err != nil {
		Error.Printf("%v\n", err)
		panic(err)
	}

	Trace.Printf("Neighbor servers are: %+v\n", neighbors)

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
		Info.Println(obj)
	}
}
