package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := ""

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("Unable to connect to provided url")
		log.Fatal("Dial error:", err)
	}

	// Means a valid hit
	defer conn.Close()
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		log.Println("Connected but error when setting read deadline")
		log.Fatal("SetReadDeadline error:", err)
	}

	// Send our handshake
	// Request that the server send it in JSON format
	handshakeMessage := `{"protocol":"json","version":1}` + "\x1e"

	// Send a text data message of our handshake
	err = conn.WriteMessage(websocket.TextMessage, []byte(handshakeMessage))
	if err != nil {
		log.Println("Connected but error when sending handshake message")
		log.Fatal("Write handshake error:", err)
	}

	// Read handshake response
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Handshake sent but error when reading response")
		log.Fatal("Read handshake response error:", err)
	}

	// Now we clear the deadline as want to listen to incoming broadcasts
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Println("Error when clearing read deadline")
		log.Fatal("Clear ReadDeadline error:", err)
	}

	log.Printf("Handshake response: %s", message)

	for {
		_, msg, error := conn.ReadMessage()
		if error != nil {
			log.Println("Error when reading broadcast message")
			log.Fatal("Read broadcast message error:", error)
		}
		log.Printf("Received broadcast message: %s", msg)
	}

}

// So I can run/debug while keeping unused variables
func UNUSED(x ...interface{}) {}
