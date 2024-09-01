package main

import (
	"fmt"
	"net"
	"strings"
)

type HTTPResponse struct {
	Status  string
	Headers map[string]string
	Body    string
}

func main() {
	fmt.Println("Server Logs: ")

	// listen on TCP port 4221
	// 0.0.0.0 -> all available network interfaces
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handleError(err, "Failed to bind to port 4221")

	for {
		// accept a new connection
		conn, err := l.Accept()
		handleError(err, "Error accepting connection")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// read request
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	handleError(err, "Error reading request")

	// parse request
	requestData := string(buffer[:n]) // trim unused buffer content
	req, err := parseRequest(requestData)
	path := req["path"].(string)
	handleError(err, "Error parsing request")

	// route request
	switch req["method"] {
	case "GET":
		switch {
		case path == "/":
			err := rootHandler(conn)
			handleError(err, "Error in rootHandler()")

		case strings.HasPrefix(path, "/echo/"):
			err := echoHandler(conn, path)
			handleError(err, "Error in echoHandler()")

		case path == "/user-agent":
			headers := req["headers"].(map[string]string)
			err := userAgentHandler(conn, headers["User-Agent"])
			handleError(err, "Error in userAgentHandler()")

		case strings.HasPrefix(path, "/files/"):
			err := getFilesHandler(conn, path)
			handleError(err, "Error in getFilesHandler()")

		default:
			err := notFoundHandler(conn)
			handleError(err, "Error in notFoundHandler()")
		}
	case "POST":
		switch {
		case strings.HasPrefix(path, "/files/"):
			body := req["body"].(string)
			err := postFilesHandler(conn, path, body)
			handleError(err, "Error in postFilesHandler()")

		default:
			err := notFoundHandler(conn)
			handleError(err, "Error in notFoundHandler()")
		}
	default:
		res := HTTPResponse{
			Status:  "405 Method not allowed",
			Headers: map[string]string{},
			Body:    "",
		}
		err := writeResponse(conn, res)
		handleError(err, "Error writing response 405")
	}
}
