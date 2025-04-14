package internal

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	users map[string]*UserConn // Map of user -> User.
	mu    sync.RWMutex
}

type RoomManager struct {
	rooms map[string]*Room // Map of roomID -> map of userID -> User
	mu    sync.RWMutex     // Mutex for thread-safe operations
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) PrintUsersInRoom(roomID string) {
	room, exists := rm.rooms[roomID]
	if !exists {
		log.Printf("Room with ID '%s' does not exist.\n", roomID)
		return
	}

	log.Printf("Users in room '%s':\n", roomID)
	log.Println("Number of users: ", len(room.users))
	for userID, _ := range room.users {
		log.Printf("- UserID: %s\n", userID)
	}
}

func (rm *RoomManager) AddUserToRoom(roomID string, user *UserConn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; !exists {
		createdRoom := Room{
			users: make(map[string]*UserConn),
		}
		rm.rooms[roomID] = &createdRoom
	}
	rm.rooms[roomID].users[user.Id] = user
}

type UserPosition struct {
	UserId   string  `json:"userId"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Sprite   string  `json:"sprite"`
	Username string  `json:"username"`
}

// Return [{ userId, x, y }, { userId, x, y }]
func (rm *RoomManager) GetUsersInRoom(roomID string, currentUserId string) []UserPosition {
	rm.rooms[roomID].mu.Lock()
	defer rm.rooms[roomID].mu.Unlock()

	room := rm.rooms[roomID]

	users := []UserPosition{}

	for _, user := range room.users {
		if user.Id != currentUserId {
			users = append(users, UserPosition{
				UserId:   user.Id,
				X:        float64(user.X),
				Y:        float64(user.Y),
				Sprite:   user.Sprite,
				Username: user.Username,
			})
		}
	}
	return users
}

// Deleting a room need optimization.
func (rm *RoomManager) RemoveUserFromRoom(conn *websocket.Conn) string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for roomId, room := range rm.rooms {
		for userId, user := range room.users {
			if user.conn == conn {
				message := Message{
					Type: "user-left",
					Payload: map[string]interface{}{
						"userId": user.Id,
					},
				}
				delete(rm.rooms[roomId].users, userId)

				if len(rm.rooms[roomId].users) == 0 {
					delete(rm.rooms, roomId)
				}
				jsonMessage, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
					return ""
				}
				rm.BroadcastToRoom(roomId, user.Id, jsonMessage)
				return user.Id
			}
		}
	}
	return ""
}

// No Error Handling here.
func (rm *RoomManager) BroadcastToRoom(roomID string, userID string, message []byte) {
	// rm.mu.RLock()
	// defer rm.mu.RUnlock()

	if room, exists := rm.rooms[roomID]; exists {
		for _, user := range room.users {
			if user.Id != userID {
				err := user.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					continue
				}
			}
		}
	}
}
