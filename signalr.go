package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

func connectToSignalR(url string) tea.Msg {

	var webSocketUrl = url
	reqHeader := http.Header{}
	// Check if any groups/query parameters have been passed through
	var urlSplit = strings.Split(url, " ")
	if len(urlSplit) > 1 {
		// Get the base Url, then we split and sort the rest into headers based off a predetermined structure
		webSocketUrl = urlSplit[0]

		// Split into header groups delimited by commas
		var headerGroups = strings.Split(strings.Join(urlSplit[1:], " "), ",")

		for _, headerGroups := range headerGroups {
			// Split into key and value based off the colon delimiter
			var headerKV = strings.SplitN(headerGroups, ":", 2)
			if len(headerKV) == 2 {
				reqHeader.Set(strings.TrimSpace(headerKV[0]), strings.TrimSpace(headerKV[1]))
			}
		}

	}

	c, _, err := websocket.DefaultDialer.Dial(webSocketUrl, reqHeader)
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
