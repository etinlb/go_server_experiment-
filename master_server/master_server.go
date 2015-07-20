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
	"strings"

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

	// clientData.Ip = req.RemoteAddr
	clientData.Ip = strings.Split(req.RemoteAddr, ":")[0]

	serverList := readServerList(serverFile)
	filteredList := filterOutClientServer(clientData, serverList)
	// remove the currently connecting server in case it was
	// TODO: DO way more logic to get it's neighbors and such

	// fmt.Printf("%v+\n", filteredList)
	// fmt.Printf("%v", req.RemoteAddr)
	// fmt.Printf("%v", req.Host)

	response, err := json.Marshal(filteredList)

	log.Println(string(response))

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	serverList = addServerIfNotDuplicate(clientData, serverList)
	writeServerList(serverList)

}

// TODO: integrate this with the addserer if not duplicate, they are doing the same thing
// return a slice without the passed in client
func filterOutClientServer(serverToFilter ClientServer, serverList ClientServerList) ClientServerList {
	var newList ClientServerList
	newList.Servers = serverList.Servers[:0]
	// fmt.Printf("%v+\n", newList)
	// fmt.Printf("%v+\n", serverToFilter)

	for _, server := range serverList.Servers {

		fmt.Printf("%+v\n", server)
		fmt.Printf("%+v\n", serverToFilter)
		if server.Port != serverToFilter.Port || server.Ip != serverToFilter.Ip {
			log.Println("here")
			newList.Servers = append(newList.Servers, server)
		}
	}
	return newList
}

func addServerIfNotDuplicate(serverToAdd ClientServer, serverList ClientServerList) ClientServerList {
	for _, server := range serverList.Servers {
		if server.Port == serverToAdd.Port && server.Ip == serverToAdd.Ip {
			return serverList
		}
	}

	var newServerList ClientServerList
	newServerList.Servers = append(serverList.Servers, serverToAdd)

	return newServerList
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
