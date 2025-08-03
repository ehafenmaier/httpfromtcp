package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Resolve the UDP address for localhost on port 42069
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error resolving UDP address: %v\n", err)
	}

	// Dial the UDP address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Error dialing UDP address: %v\n", err)
		return
	}
	defer conn.Close()

	// Create a standard io Reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	for {
		// Read a line from standard input
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from stdin: %v\n", err)
			return
		}

		// Send the line to the UDP server
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error sending data to UDP server: %v\n", err)
			return
		}
	}
}
