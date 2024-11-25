package main

import (
	"log"
	"net/http"

	"github.com/Akhilbisht798/office/server/internals/apis/v1"
	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	err := db.ConnectDB()
	if err != nil {
		log.Println("server closing.")
		return
	}
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/signup", apis.SignUp)
	r.HandleFunc("/api/v1/signin", apis.SignIn)

	r.HandleFunc("/api/v1/createroom", apis.CreateRoom)
	r.HandleFunc("/api/v1/deleteroom", apis.DeleteRoom)
	r.HandleFunc("/api/v1/getrooms", apis.GetAllRoom)

	log.Println("Server Listening at Port: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
