package internal

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/websocket"
)

var Rooms *RoomManager
var CallManager *Calls

// TODO: make handle payload properly.
// maybe later.
func join(conn *websocket.Conn, payload map[string]interface{}) {
	log.Println("join method is called.")
	cookie, ok := payload["jwt"].(string)
	if !ok {
		log.Println("cookie payload not available")
		return
	}
	spirite, ok := payload["spirite"].(string)
	if !ok {
		log.Println("spirite not available")
		return
	}
	log.Println("spirites: ", spirite)

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Println("Error with cookie: ", err.Error())
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	id := claims.Issuer

	username, err := GetUser(id)
	if err != nil {
		log.Println("Error Getting User: ", err.Error())
		return
	}

	user := &UserConn{
		conn:     conn,
		Id:       id,
		X:        payload["x"].(float64),
		Y:        payload["y"].(float64),
		Sprite:   spirite,
		Username: username,
	}

	//log.Printf("user to be added: %+v\n", user)
	room := payload["roomId"].(string)
	log.Println("room: ", room)

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
			"userId":  user.Id,
			"spirite": user.Sprite,
			"users":   Rooms.GetUsersInRoom(room, user.Id),
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
			"userId":   user.Id,
			"x":        user.X,
			"y":        user.Y,
			"spirite":  user.Sprite,
			"username": user.Username,
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

	var user *UserConn

	if room, exists := Rooms.rooms[room]; exists {
		if u, exists := room.users[userId]; exists {
			user = u
		}
	}

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

func makeCall(conn *websocket.Conn, payload map[string]interface{}) {
	roomId := payload["roomId"].(string)
	userId := payload["userId"].(string)
	remoteUserId := payload["remoteUserId"].(string)
	callId := payload["callId"].(string)

	log.Println("Making a call with CallID", callId)

	if CallManager == nil {
		CallManager = NewCallManager()
	}

	CallManager.SendChannelID(conn, userId, remoteUserId, roomId, callId)
	log.Printf("user %s wants to make a video call to user %s", userId, remoteUserId)
}

func callAccepted(conn *websocket.Conn, payload map[string]interface{}) {
	channelID := payload["channelId"].(string)
	userId := payload["userId"].(string)

	if CallManager == nil {
		CallManager = NewCallManager()
	}

	CallManager.CallAccepted(channelID, userId)
}

func leaveCall(conn *websocket.Conn, payload map[string]interface{}) {
	channelId := payload["channelId"].(string)
	userId := payload["userId"].(string)

	if CallManager == nil {
		return
	}

	CallManager.LeaveCall(channelId, userId)
}
