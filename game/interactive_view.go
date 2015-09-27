package main

import (
	"bufio"
	"fmt"
	"os"
)

func runInteractiveMode(channelCoordinator ComunicationChannels) {
	fmt.Printf("Running in interative mode!\n ")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		input, _ := reader.ReadString('\n')
		switch input {
		case "help":
			helpCommand()
		default:
			helpCommand()
		}
		fmt.Printf("Read %s", input)
	}
}

func helpCommand() {
	fmt.Printf("Commands: There are different categories of commands\n")
	fmt.Printf("    show <var>   # var can be any of the following: \n")
	fmt.Printf("                                                    gameObjects \n")
	fmt.Printf("                                                    playerObjects \n")
	fmt.Printf("                                                    physicsComponents \n")
	fmt.Printf("    send <event> # Send an event through the channel coordinator \n")
	fmt.Printf("                 # command is interactive letting you fill out the parameters of the event  \n")
}
