package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

// * = pointer/reference to the actual data of the websocket connection
// Global variable
var conn *websocket.Conn

func main() {
	p := tea.NewProgram(NewBubbleteaModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	// Cleanup connection on exit
	if conn != nil {
		conn.Close()
	}
}

// help to avoid "unused variable" message from compiler when testing
func UNUSED(x ...interface{}) {}
