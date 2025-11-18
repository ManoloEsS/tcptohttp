package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	//channel for sending lines
	linesChan := make(chan string)
	//go routine for reading lines from file and sending lines to channel
	go func() {
		//defer closing file and closing channel
		defer f.Close()
		defer close(linesChan)
		//keep track of current bytes read from file
		currentLine := ""
		//read loop
		for {
			//create 8 byte buffer to read file contents into
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				//verify its end of file and exit gracefully
				//print leftover bytes
				if errors.Is(err, io.EOF) {
					if currentLine != "" {
						linesChan <- currentLine
					}
					break
				}
				//if other error return
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			//string of read bytes
			str := string(buffer[:n])
			//split by newline
			parts := strings.Split(str, "\n")
			//loop over parts to print each part
			//if only one part loop is skipped
			for i := 0; i < len(parts)-1; i++ {
				//send aggregated string plus first split part of current str
				currentLine = currentLine + parts[i]
				linesChan <- currentLine
				//reset currentLine after sending its contents
				currentLine = ""
			}
			//add str to currentLine, if only one part adds it
			//if more than one part it adds the leftover after splitting
			currentLine += parts[len(parts)-1]
		}
	}()
	return linesChan
}
