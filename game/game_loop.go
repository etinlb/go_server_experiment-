package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

func GameLoop() {
	tickChan := time.NewTicker(time.Millisecond * 16) // about 16 milliseconds for 60 fps a second

	go func() {
		for range tickChan.C {
			for _, gameObj := range gameObjects {
				// TODO: figure out the implication of multi threading this yeah, this is
				fmt.Printf("%+v In game Loop\n", gameObj)
				gameObj.Update()
				// TODO: This is dumb
				packet := BuildObjectPackage("update", gameObj)
				clientBackend.BroadCastPackets(packet, make(map[*websocket.Conn]bool))
			}
		}
	}()

	fmt.Println("Started Game Loop Go Routine at")
}
