package internal

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

var Rooms *RoomManager

func join(conn *websocket.Conn, payload map[string]interface{}) {
	user := &User{
		conn: conn,
		Id:   payload["userId"].(string),
		X:    payload["x"].(float64),
		Y:    payload["y"].(float64),
	}
	room := payload["roomId"].(string)

	if Rooms == nil {
		Rooms = NewRoomManager()
	}
	Rooms.AddUserToRoom(room, user)

	// send user list of the user in rooms.
	message := Message{
		Type: "space-joined",
		Payload: map[string]interface{}{
			"spawn": map[string]interface{}{
				"x": user.X,
				"y": user.Y,
			},
			"users": Rooms.GetUsersInRoom(room, user.Id),
		},
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, jsonMessage)

	// broadcast to all the user this user have joined.
	message = Message{
		Type: "user-join",
		Payload: map[string]interface{}{
			"userId": user.Id,
			"x":      user.X,
			"y":      user.Y,
		},
	}
	jsonMessage, err = json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	Rooms.BroadcastToRoom(room, user.Id, jsonMessage)
}

func move(conn *websocket.Conn, payload map[string]interface{}) {
	room := payload["roomId"].(string)
	userId := payload["userId"].(string)
	x := payload["x"].(float64)
	y := payload["y"].(float64)

	var user *User

	if room, exists := Rooms.rooms[room]; exists {
		if u, exists := room[userId]; exists {
			user = u
		}
	}

	// var xDisplacement = math.Abs(float64(user.X) - float64(x))
	// var yDisplacement = math.Abs(float64(user.Y) - float64(y))
	// log.Println("Working till here.")
	// log.Printf("%v - %v", float64(user.X), float64(x))
	// log.Printf("%v - %v", float64(user.Y), float64(y))

	// Movement Rejected
	// if xDisplacement > 2 || yDisplacement > 2 {
	// 	return
	// }
	// log.Println("Working till here 2.")

	user.X = float64(x)
	user.Y = float64(y)

	message := Message{
		Type: "movement",
		Payload: map[string]interface{}{
			"userId": userId,
			"x":      x,
			"y":      y,
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	Rooms.BroadcastToRoom(room, userId, jsonMessage)
}
