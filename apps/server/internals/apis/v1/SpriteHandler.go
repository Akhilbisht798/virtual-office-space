package apis

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Akhilbisht798/office/server/internals/cloud"
	"github.com/Akhilbisht798/office/server/internals/db"
	"gorm.io/gorm"
)

type UploadSpriteRequest struct {
	Name string `json:"name"`
}

// What happen if i create a entry but client will not upload.
// for now handle that in client side so it won't happen.
func UploadSprite(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}

	var data UploadSpriteRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding body: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sprite db.Avatar
	res := db.Database.Where("name = ?", data.Name).First(&sprite)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Println("Error assessing the database: ", res.Error)
		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	if sprite.Name != "" {
		log.Println("Error sprite already exists")
		ReturnError(w, "Sprite with this name already exits.", http.StatusBadRequest)
		return
	}

	sprite.Name = data.Name
	res = db.Database.Create(&sprite)
	if res.Error != nil {
		log.Println("Error assessing the database: ", res.Error)
		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	bucket := os.Getenv("BUCKET")
	key := "avatar/" + sprite.Name
	url, err := cloud.PutPreSignedUrl(bucket, key, 10)
	if err != nil {
		log.Println("Error decoding body: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": url.URL,
	})
}

// type GetSpriteRequest struct {
// 	UserID string `json:"userID"`
// }

// func GetSprite(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		ReturnError(w, "use post method", http.StatusBadRequest)
// 		return
// 	}
// 	var data GetSpriteRequest
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		log.Println("Error decoding body: ", err)
// 		ReturnError(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var user db.User

// 	res := db.Database.Where("id = ?", data.UserID).First(&user)
// 	if res.Error != nil {
// 		log.Println("Error find item: ", err)
// 		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var avatar db.Avatar
// 	res = db.Database.Where("id = ?", user.AvatarID).First(&avatar)
// 	if res.Error != nil {
// 		log.Println("Error find item: ", err)
// 		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	bucket := os.Getenv("BUCKET")
// 	url, err := cloud.GetPreSignedUrl(bucket, avatar.Name, 10)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"sprite": url.URL,
// 	})
// }
