package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

func connectToSignalR(url string) tea.Msg {

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return connectionStatus{connected: false, err: err}
	}

	conn = c // Store globally

	// Send handshake
	handshake := `{"protocol":"json","version":1}` + "\x1e"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(handshake)); err != nil {
		conn.Close()
		return connectionStatus{connected: false, err: err}
	}

	// Read handshake response
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, _, err = conn.ReadMessage()
	err = conn.SetReadDeadline(time.Time{})

	if err != nil {
		conn.Close()
		return connectionStatus{connected: false, err: err}
	}

	return connectionStatus{connected: true, err: nil}
}

func listenForMessages() tea.Msg {
	if conn == nil {
		return SignalRError{err: fmt.Errorf("no connection")}
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return SignalRError{err: err}
	}

	// Parse the message
	timestamp := time.Now().Format("15:04:05")
	return SignalRResponse{
		data: fmt.Sprintf("[%s] %s", timestamp, string(message)),
	}
}
