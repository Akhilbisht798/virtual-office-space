package internal

var Rooms *RoomManager
var CallManager *Calls

func init() {
	Rooms = NewRoomManager()
	CallManager = NewCallManager()
}
