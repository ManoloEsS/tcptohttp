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
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	//read bytes from reader
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	//parse RequestLine fields from read bytes into string
	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}

	//create new Request struct from parsed fields
	newRequest := Request{
		RequestLine: *requestLine,
	}

	return &newRequest, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
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
