package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const network = "udp"
const serverAddress = "localhost:42070"

func main() {
	//create udp address
	udpAddr, err := net.ResolveUDPAddr(network, serverAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to serverAddress %s for UDP traffic: %s\n", serverAddress, err)
		os.Exit(1)
	}

	fmt.Println("Ready to send UDP traffic on", serverAddress)

	//start connection to remote address, local address is nil
	conn, err := net.DialUDP(network, nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting connection: %s", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddress)

	//create reader
	input := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		//read stream of bytes from stdin
		message, err := input.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v", err)
		}

		//send stream of bytes through connection
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message through connection: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
