package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const network = "tcp"
const port = ":42069"

func main() {
	//start tcp sever listening on port
	fmt.Printf("Starting to listen on port localhost%s\n", port)
	listener, err := net.Listen(network, port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		//wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error connecting: %s\n", err.Error())
			continue
		}
		fmt.Println("Connection has been accepted from", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		printRequestLine(r)

		fmt.Println("Connection to ", conn.RemoteAddr(), "has been closed...")
	}
}

func printRequestLine(r *request.Request) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	//channel for sending lines
// 	linesChan := make(chan string)
// 	//go routine for reading lines from file and sending lines to channel
// 	go func() {
// 		//defer closing file and closing channel
// 		defer f.Close()
// 		defer close(linesChan)
// 		//keep track of current bytes read from file
// 		currentLine := ""
// 		//read loop
// 		for {
// 			//create 8 byte buffer to read file contents into
// 			buffer := make([]byte, 8)
// 			n, err := f.Read(buffer)
// 			if err != nil {
// 				//verify its end of file and exit gracefully
// 				//print leftover bytes
// 				if errors.Is(err, io.EOF) {
// 					if currentLine != "" {
// 						linesChan <- currentLine
// 					}
// 					break
// 				}
// 				//if other error return
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
// 			//string of read bytes
// 			str := string(buffer[:n])
// 			//split by newline
// 			parts := strings.Split(str, "\n")
// 			//loop over parts to print each part
// 			//if only one part loop is skipped
// 			for i := 0; i < len(parts)-1; i++ {
// 				//send aggregated string plus first split part of current str
// 				currentLine = currentLine + parts[i]
// 				linesChan <- currentLine
// 				//reset currentLine after sending its contents
// 				currentLine = ""
// 			}
// 			//add str to currentLine, if only one part adds it
// 			//if more than one part it adds the leftover after splitting
// 			currentLine += parts[len(parts)-1]
// 		}
// 	}()
// 	return linesChan
// }
