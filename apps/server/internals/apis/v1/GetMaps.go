package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator"
)

func GetMaps(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ReturnError(w, "use get method", http.StatusBadRequest)
		return
	}

	var maps []db.Map
	res := db.Database.Find(&maps)
	if res.Error != nil {
		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("fetched spaces: %#v", maps)
	response := map[string]interface{}{
		"maps": maps,
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

type GetSpaceRequest struct {
	Jwt string `json:"jwt" validate:"required"`
}

func GetSpaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use get method", http.StatusBadRequest)
		return
	}
	var data GetSpaceRequest
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

	token, err := jwt.ParseWithClaims(data.Jwt, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Println("Error with cookie: ", err.Error())
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)
	id := claims.Issuer

	var spaces []db.Space

	res := db.Database.Where("user_id = ?", id).Find(&spaces)
	if res.Error != nil {
		ReturnError(w, res.Error.Error(), http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"spaces": spaces,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println("Error: ", err.Error())
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}
}
