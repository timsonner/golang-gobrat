// Client

package main

import (
	"bufio"        // Package for buffered I/O
	"encoding/gob" // Package for encoding and decoding data
	"fmt"          // Package for formatted I/O
	"net"          // Package for network operations
	"os"           // Package for OS functions
	"strings"      // Package for string operations
)

// Message represents the structure of the message to be sent to the server
type Message struct {
	Content string // Command field holds the command to be executed
}

// Response represents the structure of the response received from the server
type Response struct {
	Output       string // Output field holds the command execution output
	ErrorMessage string // ErrorMessage field holds the error message, if any
}

func main() {
	// Connect to the server at TCP address "localhost:1234"
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close() // Close the connection before exiting the main function

	// Create an encoder to encode the message and send it to the server
	encoder := gob.NewEncoder(conn)

	// Create a decoder to decode the response received from the server
	decoder := gob.NewDecoder(conn)

	// Create a scanner to read user input from the standard input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter a command (or type 'exit' to quit): ")
		scanner.Scan()
		command := scanner.Text()

		// Create a message with the user input command
		message := Message{
			Content: command,
		}

		// Encode and send the message to the server
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding message:", err)
			break
		}

		// Check if the user wants to exit
		if strings.ToLower(command) == "exit" {
			break
		}

		// Create an empty response to hold the decoded response from the server
		var response Response

		// Decode the response received from the server into the response variable
		err = decoder.Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			break
		}

		// Print the server response
		if response.ErrorMessage != "" {
			fmt.Println("Server response:", response.ErrorMessage)
		} else {
			fmt.Println("Server response:\n\n", response.Output)
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading user input:", scanner.Err())
	}
}
