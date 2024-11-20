package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity. In production, limit to trusted origins.
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Run the command and stream its output
	cmd := exec.Command("ping", "google.com", "-t") // Example: Ping command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error getting stdout pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("Error starting command:", err)
		return
	}

	// Read command output line by line and send it to the WebSocket
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			message := string(buf[:n])
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Error writing to WebSocket:", err)
				break
			}
		}
		if err != nil {
			log.Println("Command finished:", err)
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Error waiting for command to finish:", err)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
