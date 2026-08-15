package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "localhost:7777"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to Raven DB server at %s: %v\n", serverAddr, err)
		os.Exit(1)
	}
	defer conn.Close()

	connReader := bufio.NewReader(conn)
	// Read server welcome message
	greeting, err := connReader.ReadString('\n')
	if err == nil {
		fmt.Print(greeting)
	}

	stdinScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("raven> ")

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())
		if input == "" {
			fmt.Print("raven> ")
			continue
		}
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Bye!")
			break
		}

		// Send command to server
		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Printf("Error writing to server: %v\n", err)
			break
		}

		// Read response from server
		response, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection closed by server: %v\n", err)
			break
		}

		fmt.Print(response)
		fmt.Print("raven> ")
	}
}
