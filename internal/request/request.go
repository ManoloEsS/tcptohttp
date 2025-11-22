package request

import (
	"errors"
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	//read bytes from reader
	reqBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	//parse RequestLine fields from read bytes into string
	requestLine, err := parseRequestLine(string(reqBytes))
	if err != nil {
		return nil, err
	}

	//create new Request struct from parsed fields
	newRequest := Request{
		RequestLine: *requestLine,
	}

	return &newRequest, nil
}

func parseRequestLine(request string) (*RequestLine, error) {
	requestLines := strings.Split(request, "\r\n")
	if len(requestLines) < 1 {
		return nil, errors.New("request does not contain request line")
	}

	requestStartLineParts := strings.Split(requestLines[0], " ")
	if len(requestStartLineParts) != 3 {
		return nil, errors.New("invalid number of parts in request line")
	}

	method := requestStartLineParts[0]
	target := requestStartLineParts[1]
	rawHttpVersion := requestStartLineParts[2]

	for _, char := range method {
		if !unicode.IsUpper(char) {
			return nil, errors.New("invalid method")
		}
	}

	httpVersionParts := strings.Split(rawHttpVersion, "/")
	if httpVersionParts[1] != "1.1" {
		return nil, errors.New("invalid http version")
	}

	requestLine := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   httpVersionParts[1],
	}

	return &requestLine, nil
}
