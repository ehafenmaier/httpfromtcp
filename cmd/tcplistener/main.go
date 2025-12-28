package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"io"
	"net"
	"strings"
)

func main() {
	port := "42069"
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error opening TCP listener on port %s: %v\n", port, err)
		return
	}
	defer l.Close()
	fmt.Printf("Listening on TCP port %s\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			conn.Close()
			continue
		}

		rl := r.RequestLine
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", rl.Method, rl.RequestTarget, rl.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range r.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Printf("%s\n", r.Body)
	}
}

// Reads lines from a file and returns them as a channel of strings.
func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		currentLine := ""

		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err == io.EOF {
				if currentLine != "" {
					lines <- currentLine
				}
				return
			}

			currentLine += string(b[:n])

			if strings.Contains(currentLine, "\n") {
				parts := strings.Split(currentLine, "\n")
				lines <- parts[0]
				currentLine = "" + parts[1]
			}
		}
	}()

	return lines
}
