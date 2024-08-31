package main

import (
	"net"
)

func rootHandler(connection net.Conn) error {
	defer connection.Close()

	response := "HTTP/1.1 200 OK\r\n" + END_HEADER
	_, err := connection.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}
