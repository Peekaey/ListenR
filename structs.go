package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/gorilla/websocket"
)

type connectionStatus struct {
	connected bool
	err       error
}

type SignalRResponse struct {
	data string
}

type SignalRError struct {
	err error
}

type BubbleteaModel struct {
	conn      *websocket.Conn
	connected bool
	messages  []string
	status    string
	err       error
	quitting  bool
	url       string

	// for the prompt screen
	promptMode bool
	ti         textinput.Model
}
