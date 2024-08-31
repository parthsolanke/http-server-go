package main

import (
	"fmt"
	"net"
	"os"
)

var END_HEADER = "\r\n "

func main() {
	fmt.Println("My Logs: ")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handleError(err, "Failed to bind to port 4221")

	conn, err := l.Accept()
	handleError(err, "Error accepting connection")

	// read request
	buffer := make([]byte, 4096)
	_, err = conn.Read(buffer)
	handleError(err, "Error reading request")

	// parse request
	req, err := parseRequest(string(buffer))
	handleError(err, "Error parsing request")

	// route request
	switch req["method"] {
	case "GET":
		switch req["path"] {
		case "/":
			err := rootHandler(conn)
			handleError(err, "Error writing response for rootHandler()")
		default:
			res := "HTTP/1.1 404 Not Found\r\n" + END_HEADER
			_, err := conn.Write([]byte(res))
			handleError(err, "Error writing response 404")
		}
	default:
		fmt.Println(fmt.Sprintf("%v", req["method"]) + "Method not allowed")
		os.Exit(1)
	}
}
