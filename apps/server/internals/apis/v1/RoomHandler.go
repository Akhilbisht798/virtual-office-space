package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/Akhilbisht798/office/server/internals/cloud"
	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name   string `json:"name" validate:"required"`
	MapId  string `json:"mapId" validate:"required"`
	Public bool   `json:"public"`
	Jwt    string `json:"jwt" validate:"required"`
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
	log.Println("Type of Public: ", reflect.TypeOf(data.Public))

	token, err := jwt.ParseWithClaims(data.Jwt, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Println("Error with cookie: ", err.Error())
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	id := claims.Issuer

	var mapUsed db.Map
	res := db.Database.First(&mapUsed, "id = ?", data.MapId)
	if res.Error != nil {
		log.Println("Error finding map: ", err)
		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	mapName := strings.Split(mapUsed.Name, ".")[0]
	thumbnail := "maps/thumbnail/" + mapName + ".png"

	room := db.Space{
		ID:        uuid.New().String(),
		Name:      data.Name,
		MapID:     data.MapId,
		Public:    data.Public,
		Thumbnail: &thumbnail,
		UserID:    id,
	}
	log.Printf("Room to be saved: %+v\n", room)

	result := db.Database.Create(&room)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Println("roomID", room.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"roomId": room.ID,
	})
}

type DeleteRoomRequest struct {
	RoomId string `json:"roomID" validate:"required"`
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}

	var data DeleteRoomRequest
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

	result := db.Database.Delete(&db.Space{}, "id = ?", data.RoomId)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "space deleted successfully",
	})
}

func GetAllRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}

	var spaces []db.Space
	result := db.Database.Where("public = ?", true).Find(&spaces)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("fetched spaces: %#v", spaces)
	response := map[string]interface{}{
		"spaces": spaces,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println("Error: ", err.Error())
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}
}

type JoinRoomRequest struct {
	RoomID string `json:"roomID" validate:"required"`
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}

	var data JoinRoomRequest
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

	log.Println("roomID", data.RoomID)
	space := &db.Space{}
	//result := db.Database.First(space, data.RoomID)
	result := db.Database.Where("id = ?", data.RoomID).First(space)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	maps := &db.Map{}
	result = db.Database.Where("id = ?", space.MapID).First(maps)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}
	// get all presigned url for the name
	bucket := os.Getenv("BUCKET")

	mapKey := "maps/" + maps.Name
	url, err := cloud.GetPreSignedUrl(bucket, mapKey, 60)
	if err != nil {
		log.Println("Error getting presigned url: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url":     url.URL,
		"spaceID": space.ID,
	})
}
