package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// request state enum
type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

// Function that reads bytes from a io.Reader and parses then into a Request Line into a new Request
func RequestFromReader(reader io.Reader) (*Request, error) {
	//create buffer to store stream of bytes
	buff := make([]byte, bufferSize)
	//index of already read bytes limit
	readToIndex := 0

	//new request with requestStateInitialized state
	newRequest := Request{state: requestStateInitialized}

	//loop to read from reader while request is not 'requestStateDone'
	for newRequest.state != requestStateDone {
		//increase size of buffer if read to index is larger than len of buffer
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		//read from reader appending new bytes to the end of the buffer
		read, readErr := reader.Read(buff[readToIndex:])
		if readErr != nil {
			if !errors.Is(readErr, io.EOF) {
				return nil, readErr
			}
		}

		//move index of read bytes forward
		readToIndex += read

		//parse bytes into request until read to index
		//if succesful newRequest.state is updated to done
		parsed, err := newRequest.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		//copy unparsed and get rid of already processed bytes
		copy(buff, buff[parsed:])
		readToIndex -= parsed

		//if reached EOF, wether succesful or unsuccesful parsing we break from loop
		if readErr == io.EOF {
			break
		}
	}

	//check if successful parse
	if newRequest.state != requestStateDone {
		return nil, fmt.Errorf("incomplete request: reached EOF before parsing finished")
	}
	return &newRequest, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	//search for crlf in bytes
	idx := bytes.Index(data, []byte(crlf))

	//if not found return 0 bytes processed
	if idx == -1 {
		return nil, 0, nil
	}

	//if found we create request line from bytes
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	//succesful request line with amount of bytes processed is returned
	return requestLine, idx + 2, nil
}

// Function that creates a Request line from a string
func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	target := parts[1]
	rawHttpVersion := parts[2]

	for _, char := range method {
		if !unicode.IsUpper(char) {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	httpVersionParts := strings.Split(rawHttpVersion, "/")
	if len(httpVersionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	if httpVersionParts[0] != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP token: %s", httpVersionParts[0])
	}
	if httpVersionParts[1] != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersionParts[1])
	}

	requestLine := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   httpVersionParts[1],
	}

	return &requestLine, nil
}

// method that parses and returns a Request line from bytes
// if successful request state is updated
func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, parsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if parsed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateDone
		return parsed, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a requestStateDone state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
