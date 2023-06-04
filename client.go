// Client

package main

import (
	"bufio"        // Package for buffered I/O
	"encoding/gob" // Package for encoding and decoding data
	"fmt"          // Package for formatted I/O
	"net"          // Package for network operations
	"os"           // Package for operating system functionality
)

type Message struct {
	Content string // Command field holds the command to be sent to the server
}

type Response struct {
	Output string // Output field holds the command execution output received from the server
}

func main() {
	for {
		// Prompt the user for an input message
		fmt.Print("Enter a command (or type 'exit' to quit): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()

		// Check if the user wants to exit
		if command == "exit" {
			break
		}

		// Establish a connection to the server at "localhost:1234" using TCP protocol
		conn, err := net.Dial("tcp", "localhost:1234")
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		defer conn.Close() // Close the connection before exiting the main function

		// Create an encoder to encode Go values into a binary format to be sent over the network connection
		encoder := gob.NewEncoder(conn)

		// Create a message with the user's input command
		message := Message{
			Content: command,
		}

		// Encode the message and send it to the server through the connection
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding message:", err)
			return
		}

		// Create a decoder to decode the response message received from the server
		decoder := gob.NewDecoder(conn)

		// Create a new empty response to hold the decoded response message content
		var response Response

		// Decode the response message received from the server into the response variable
		err = decoder.Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return
		}

		fmt.Println("Server response:", response.Output) // Print the response received from the server
	}
}