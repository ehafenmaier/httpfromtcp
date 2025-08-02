package main

import (
	"fmt"
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

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection closed by client")
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
