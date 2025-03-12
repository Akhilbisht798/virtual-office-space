package internal

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Calls struct {
	rooms map[string][]string
}

func NewCallManager() *Calls {
	return &Calls{
		rooms: make(map[string][]string),
	}
}

func (c Calls) SendChannelID(conn *websocket.Conn, userId, remoteUserId, roomId, callId string) {
	isOnCall := false
	var onRoom string

	for roomId, users := range c.rooms {
		for _, user := range users {
			if user == remoteUserId {
				isOnCall = true
				onRoom = roomId
				break
			}
		}
		if isOnCall {
			break
		}
	}

	if !isOnCall {
		c.RequestCall(userId, remoteUserId, roomId, callId, Rooms)
		message := Message{
			Type: "call-req",
			Payload: map[string]interface{}{
				"channel": callId,
			},
		}
		jsonMessage, _ := json.Marshal(message)
		conn.WriteMessage(websocket.TextMessage, jsonMessage)
		return
	}
	message := Message{
		Type: "call-req",
		Payload: map[string]interface{}{
			"channel": onRoom,
		},
	}
	jsonMessage, _ := json.Marshal(message)
	conn.WriteMessage(websocket.TextMessage, jsonMessage)
}

func (c Calls) RequestCall(userId, remoteUserId, roomId, callId string, rm *RoomManager) {
	rmConn := rm.rooms[roomId].users[remoteUserId]
	message := Message{
		Type: "call-req",
		Payload: map[string]interface{}{
			"channel": callId,
		},
	}
	jsonMessage, _ := json.Marshal(message)
	rmConn.conn.WriteMessage(websocket.TextMessage, jsonMessage)
}

func (c Calls) CallAccepted(channelId, userId string) {
	if _, exits := c.rooms[channelId]; !exits {
		c.rooms[channelId] = []string{}
	}

	c.rooms[channelId] = append(c.rooms[channelId], userId)
}
