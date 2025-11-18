package main

import (
	"fmt"
	"log"
	"net"
)

const inputFilePath = "messages.txt"
const network = "tcp"
const port = "localhost:42069"

func main() {
	//start tcp sever listening on port
	fmt.Printf("Starting to listen on port %s\n", port)
	listener, err := net.Listen(network, port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", port)
	}
	fmt.Println("Listening for TCP traffic on", port)

	for {
		//wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error connecting: %s\n", err.Error())
			continue
		}

		fmt.Println("Connection has been accepted from", conn.RemoteAddr())
		//read lines from channel from connection and print them
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "has been closed...")
	}
}
