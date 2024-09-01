package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func rootHandler(conn net.Conn) error {
	response := HTTPResponse{
		Status:  "200 OK",
		Headers: map[string]string{},
		Body:    "",
	}
	return writeResponse(conn, response)
}

func notFoundHandler(conn net.Conn) error {
	response := HTTPResponse{
		Status:  "404 Not Found",
		Headers: map[string]string{},
		Body:    "",
	}
	return writeResponse(conn, response)
}

func echoHandler(conn net.Conn, path string, encoding string) error {
	echoStr := strings.TrimPrefix(path, "/echo/")
	response := HTTPResponse{
		Status: "200 OK",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: echoStr,
	}

	if encoding == "gzip" {
		response.Headers["Content-Encoding"] = "gzip"

		compressedData, compressedLength, err := compressGzip(response.Body)
		if err != nil {
			return err
		}

		response.Body = compressedData
		response.Headers["Content-Length"] = fmt.Sprintf("%d", compressedLength)
	} else {
		response.Headers["Content-Length"] = fmt.Sprintf("%d", len(echoStr))
	}

	return writeResponse(conn, response)
}

func userAgentHandler(conn net.Conn, userAgent string) error {
	response := HTTPResponse{
		Status: "200 OK",
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(userAgent)),
		},
		Body: userAgent,
	}
	return writeResponse(conn, response)
}

func getFilesHandler(conn net.Conn, path string) error {
	fileName := strings.TrimPrefix(path, "/files/")
	dir := os.Args[2]

	content, err := os.ReadFile(dir + fileName)
	if err != nil {
		return notFoundHandler(conn)
	} else {
		response := HTTPResponse{
			Status: "200 OK",
			Headers: map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": fmt.Sprintf("%d", len(content)),
			},
			Body: string(content),
		}
		return writeResponse(conn, response)
	}
}

func postFilesHandler(conn net.Conn, path string, content string) error {
	fileName := strings.TrimPrefix(path, "/files/")
	dir := os.Args[2]

	if err := os.WriteFile(dir+fileName, []byte(content), 0644); err != nil {
		return notFoundHandler(conn)
	} else {
		response := HTTPResponse{
			Status:  "201 Created",
			Headers: map[string]string{},
			Body:    "",
		}
		return writeResponse(conn, response)
	}
}
