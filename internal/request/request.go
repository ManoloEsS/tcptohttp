package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	crlf       = "\r\n"
	bufferSize = 8
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	//read bytes from reader
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	newRequest := Request{State: 1}

	for newRequest.State != 0 {

		if len(buff) <= readToIndex {
			newBuff := make([]byte, len(buff), cap(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		read, err := reader.Read(buff[readToIndex:])
		if err == io.EOF {
			newRequest.State = 0
			break
		}
		readToIndex += read

		parsed, err := newRequest.parse(buff)
		if err != nil {
			return nil, err
		}

		newBuff := make([]byte, len(buff), cap(buff))
		copy(newBuff, buff)
		buff = newBuff

		readToIndex -= parsed

	}

	return &newRequest, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, nil
}

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

func (r *Request) parse(data []byte) (int, error) {
	if r.State == 1 {
		parseRequestLine(data)
	}
}
