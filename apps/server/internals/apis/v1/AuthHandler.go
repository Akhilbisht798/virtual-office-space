package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

func ReturnError(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": err,
	})
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}
	var data SignUpRequest
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

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating hash: ", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := db.User{
		Username: data.Username,
		Password: string(hasedPassword),
	}

	result := db.Database.Create(&user)
	if result.Error != nil {
		log.Println("Error: ", result.Error.Error())
		ReturnError(w, result.Error.Error(), http.StatusBadRequest)
		return
	}
	db.Database.First(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"userId": user.ID,
	})
}

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ReturnError(w, "use post method", http.StatusBadRequest)
		return
	}
	var data SignInRequest
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
	var user db.User
	db.Database.Where("username = ?", data.Username).First(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		log.Println("Error wrong password", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: user.ID,
		ExpiresAt: &jwt.Time{
			Time: time.Now().Add(time.Hour * 200),
		},
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Error in getting a token", err)
		ReturnError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
