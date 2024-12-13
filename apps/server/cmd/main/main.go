package main

import (
	"log"
	"net/http"

	"github.com/Akhilbisht798/office/server/internals/apis/middleware"
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

	r.HandleFunc("/api/v1/signup", middleware.ApplyMiddleware(apis.SignUp, middleware.EnableCORS))
	r.HandleFunc("/api/v1/signin", middleware.ApplyMiddleware(apis.SignIn, middleware.EnableCORS))

	r.HandleFunc("/api/v1/createroom", middleware.ApplyMiddleware(apis.CreateRoom, middleware.EnableCORS))
	r.HandleFunc("/api/v1/deleteroom", middleware.ApplyMiddleware(apis.DeleteRoom, middleware.EnableCORS))
	r.HandleFunc("/api/v1/getrooms", middleware.ApplyMiddleware(apis.GetAllRoom, middleware.EnableCORS))
	r.HandleFunc("/api/v1/joinroom", middleware.ApplyMiddleware(apis.JoinRoom, middleware.EnableCORS))

	log.Println("Server Listening at Port: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
