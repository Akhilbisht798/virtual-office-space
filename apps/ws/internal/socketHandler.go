package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// TODO: only allow frontend connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer func() {
		userId := Rooms.RemoveUserFromRoom(conn)
		if userId != "" {
			CallManager.RemoveUserFromCall(userId)
		}
		conn.Close()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		if messageType == websocket.TextMessage {
			eventHandler(conn, message)
		}
	}
}

func eventHandler(conn *websocket.Conn, message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Println("Failed to unmarshal message:", err)
		return
	}

	switch msg.Type {
	case "join":
		join(conn, msg.Payload)
	case "move":
		move(conn, msg.Payload)
	case "make-call":
		makeCall(conn, msg.Payload)
	case "call-accept":
		callAccepted(conn, msg.Payload)
	case "leave-call":
		leaveCall(conn, msg.Payload)
	}
}
