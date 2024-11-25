package apis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/go-playground/validator"
)

type CreateRoomRequest struct {
	Name      string `json:"name" validate:"required"`
	MapId     string `json:"mapId" validate:"required"`
	Thumbnail string `json:"thumbnail`
	Public    bool   `json:"public"`
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}
	var data CreateRoomRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding body: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		log.Println("Error in the request body: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	room := db.Space{
		Name:      data.Name,
		MapID:     data.MapId,
		Public:    data.Public,
		Thumbnail: &data.Thumbnail,
	}

	result := db.Database.Create(&room)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"roomId": room.ID,
	})
}

//func DeleteRoom(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{
//		"userId": user.ID,
//	})
//}
//
//func JoinRoom(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{
//		"userId": user.ID,
//	})
//}
//
//func GetAllRoom(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(map[string]string{
//		"userId": user.ID,
//	})
//}
