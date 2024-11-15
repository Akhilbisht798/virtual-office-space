package internal

import "github.com/gorilla/websocket"

type User struct {
	conn *websocket.Conn
	Id   string  `json:"id"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}
