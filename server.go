// Server

package main

import (
	// Package for buffered I/O
	"encoding/gob" // Package for encoding and decoding data
	"fmt"          // Package for formatted I/O
	"net"          // Package for network operations
	"os/exec"      // Package for executing commands
	"strings"      // Package for string operations
)

// Message represents the structure of the message received from the client
type Message struct {
	Content string // Command field holds the command to be executed
}

// Response represents the structure of the response message to be sent to the client
type Response struct {
	Output       string // Output field holds the command execution output
	ErrorMessage string // ErrorMessage field holds the error message, if any
}

// handleConnection is a function that handles the communication with a client connection
func handleConnection(conn net.Conn) {
	// Create a decoder to decode the binary data received from the client connection
	decoder := gob.NewDecoder(conn)

	// Create an encoder to encode the response message and send it back to the client
	encoder := gob.NewEncoder(conn)

	for {
		// Create a new empty message to hold the decoded message content
		var message Message

		// Decode the message received from the client into the message variable
		err := decoder.Decode(&message)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			break
		}

		fmt.Println("Received command:", message.Content) // Print the received command

		// Create a response message with empty output and error message fields
		response := Response{}

		// Execute the command and capture the output and error streams
		cmd := exec.Command("bash", "-c", message.Content)
		output, err := cmd.Output()
		if err != nil {
			response.ErrorMessage = err.Error()
		}
		response.Output = string(output)

		// Send the response message back to the client
		err = encoder.Encode(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
			break
		}

		// Check if the client wants to exit
		if strings.ToLower(message.Content) == "exit" {
			break
		}
	}

	conn.Close() // Close the connection with the client
}

func main() {
	// Listen for incoming connections on TCP port "localhost:1234"
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close() // Close the listener before exiting the main function

	fmt.Println("Server started. Waiting for connections...")

	// Accept and handle client connections in a loop
	for {
		// Accept a new client connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		// Handle the client connection concurrently in a separate goroutine
		go handleConnection(conn)
	}
}
