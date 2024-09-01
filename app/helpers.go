package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func parseRequestLine(line string) (string, string, string, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return "", "", "", errors.New("invalid request line format")
	}
	return parts[0], parts[1], parts[2], nil
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			break
		}
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = headerParts[1]
		}
	}
	return headers
}

func parseBody(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\r\n")
}

func parseRequest(input string) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})

	lines := strings.Split(input, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	// request line
	method, path, protocol, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}
	parsed["method"] = method
	parsed["path"] = path
	parsed["protocol"] = protocol

	// headers
	headers := parseHeaders(lines[1:])
	parsed["headers"] = headers

	// body
	body := parseBody(lines[len(headers)+2:])
	parsed["body"] = body

	return parsed, nil

}

func compressGzip(data string) (string, int, error) {
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)

	_, err := gzipWriter.Write([]byte(data))
	if err != nil {
		return "", 0, fmt.Errorf("failed to write to gzip writer: %w", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return "", 0, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return compressedBody.String(), compressedBody.Len(), nil
}

func getSupportedEncoding(encodingHeader string) string {
	supportedEncodings := map[string]bool{"gzip": true}
	encodings := strings.Split(encodingHeader, ",")

	for _, e := range encodings {
		trimmedEncoding := strings.TrimSpace(e)
		if supportedEncodings[trimmedEncoding] {
			return trimmedEncoding
		}
	}

	return ""
}

func writeResponse(conn net.Conn, response HTTPResponse) error {
	defer conn.Close()
	var responseBuilder strings.Builder

	// status
	responseBuilder.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", response.Status))
	// headers and body
	if len(response.Headers) > 0 || response.Body != "" {
		for key, value := range response.Headers {
			responseBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
		responseBuilder.WriteString("\r\n") // end headers
		responseBuilder.WriteString(response.Body)
	} else {
		responseBuilder.WriteString("\r\n") // end response
	}

	_, err := conn.Write([]byte(responseBuilder.String()))
	return err
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg+":", err.Error())
		os.Exit(1)
	}
}
