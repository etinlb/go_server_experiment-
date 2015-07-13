package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	// "github.com/etinlb/go_game/backend"
	// "github.com/gorilla/websocket"
)

var serverFile string

type ClientServer struct {
	Port int `json:"port"`
	Ip   string
}

type clientPackage struct {
	// the struct that represents the json data received from the client
	Port int `json:"port"`
}

type ClientServerList struct {
	Servers []ClientServer `json:"servers"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var clientData ClientServer
	err := decoder.Decode(&clientData)

	if err != nil {
		panic(err)
	}

	log.Println("Client connecting from " + req.RemoteAddr + ": using port : " +
		strconv.Itoa(clientData.Port))

	clientData.Ip = req.RemoteAddr

	serverList := readServerList(serverFile)
	fmt.Printf("%v", serverList)

	response, err := json.Marshal(serverList)

	log.Println(string(response))

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	serverList.Servers = append(serverList.Servers, clientData)
	writeServerList(serverList)

}

// TODO: Write to a database rather than file
func readServerList(fileStr string) ClientServerList {
	file, err := ioutil.ReadFile(fileStr)

	if err != nil {
		log.Println(err)
	}

	var serverList ClientServerList
	json.Unmarshal(file, &serverList)

	return serverList
}

// TODO: Write to a database rather than file
func writeServerList(input ClientServerList) {

	j, jerr := json.MarshalIndent(input, "", "  ")
	if jerr != nil {
		fmt.Println("jerr:", jerr.Error())
	}

	err := ioutil.WriteFile(serverFile, j, 0644)
	log.Println(err)
}

func main() {
	//TODO: read from a database and not a json file
	list_file := flag.String("server-list", "server_list.json", "The file to "+
		"store the connected servers")
	flag.Parse()

	_, e := ioutil.ReadFile(*list_file)

	// make sure we can read the file
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	serverFile = *list_file

	http.HandleFunc("/jackIn", handler)

	addr := fmt.Sprintf("0.0.0.0:%d", 4000)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())

}
