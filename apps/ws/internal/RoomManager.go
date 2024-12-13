package internal

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	rooms map[string]map[string]*UserConn // Map of roomID -> map of userID -> User
	mu    sync.RWMutex                    // Mutex for thread-safe operations
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]map[string]*UserConn),
	}
}

func (rm *RoomManager) PrintUsersInRoom(roomID string) {
	users, exists := rm.rooms[roomID]
	if !exists {
		log.Printf("Room with ID '%s' does not exist.\n", roomID)
		return
	}

	log.Printf("Users in room '%s':\n", roomID)
	log.Println("Number of users: ", len(users))
	for userID, _ := range users {
		log.Printf("- UserID: %s\n", userID)
	}
}

func (rm *RoomManager) AddUserToRoom(roomID string, user *UserConn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; !exists {
		rm.rooms[roomID] = make(map[string]*UserConn)
	}
	rm.rooms[roomID][user.Id] = user
	log.Printf("Added User %s to Room %s", user.Id, roomID)
}

type UserPosition struct {
	UserId string  `json:"userId"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Sprite string  `json:"sprite"`
}

// Return [{ userId, x, y }, { userId, x, y }]
func (rm *RoomManager) GetUsersInRoom(roomID string, currentUserId string) []UserPosition {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room := rm.rooms[roomID]

	users := []UserPosition{}

	for _, user := range room {
		if user.Id != currentUserId {
			users = append(users, UserPosition{
				UserId: user.Id,
				X:      float64(user.X),
				Y:      float64(user.Y),
				Sprite: user.Sprite,
			})
		}
	}
	return users
}

func (rm *RoomManager) RemoveUserFromRoom(conn *websocket.Conn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for room, users := range rm.rooms {
		for userId, user := range users {
			log.Println(user.Id)
			if user.conn == conn {
				log.Println("userId to be removed: ", user.Id)
				message := Message{
					Type: "user-left",
					Payload: map[string]interface{}{
						"userId": user.Id,
					},
				}
				delete(users, userId)

				if len(users) == 0 {
					delete(rm.rooms, room)
				}
				jsonMessage, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
					return
				}
				rm.BroadcastToRoom(room, user.Id, jsonMessage)
				return
			}
		}
	}
}

// No Error Handling here.
func (rm *RoomManager) BroadcastToRoom(roomID string, userID string, message []byte) {
	// rm.mu.RLock()
	// defer rm.mu.RUnlock()

	if room, exists := rm.rooms[roomID]; exists {
		for _, user := range room {
			if user.Id != userID {
				err := user.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					continue
				}
			}
		}
	}
}
