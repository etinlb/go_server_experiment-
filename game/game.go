package main

import (
	// "flag"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/etinlb/go_game/backend"
	"github.com/gorilla/websocket"
)

var gameObjects map[string]*GameObject

var connections map[*websocket.Conn]bool

// map that keeps track of what data came from what client
var clients map[*websocket.Conn]ClientData

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

func initializeConectionVaribles(conn *websocket.Conn) {
	// initialize the connection
	connections[conn] = true
	clients[conn] = ClientData{Client: conn, GameObjects: make(map[string]*GameObject)}
	SyncClient(conn)
}

func main() {
	port := flag.Int("port", 8080, "port to serve on")
	// TODO: have this address passed from the other server
	neighborPort := flag.Int("neighbor", 8081, "port to connect to neighbor on")
	dir := flag.String("directory", "../web/", "directory of web files")

	flag.Parse()

	// TODO: WTF IS THIS SHIT
	portAsString := strconv.Itoa(*neighborPort)

	// =========Game Initializations============
	// keyed by id
	gameObjects = make(map[string]*GameObject)
	clients = make(map[*websocket.Conn]ClientData)

	// =========Connection Initializations============
	connections = make(map[*websocket.Conn]bool)

	// =========Connection handlers===================
	backend.InitHandlerFunctions(HandleEvent, cleanUpSocket, initializeConectionVaribles)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", backend.WsHandler)
	// the socket to read incoming connections from the master server
	http.HandleFunc("/masterSocket", backend.WsHandler)
	log.Println(portAsString)
	connectionUrl := "http://localhost:" + portAsString + "/masterSocket"
	x := newClientConnection(connectionUrl)
	log.Println(x)
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
