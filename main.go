package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
	defer f.Close()

	ch := getLinesChannel(f)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
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
