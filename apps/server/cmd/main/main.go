package main

import (
	"log"
	"net/http"
	"os"

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

	//r.HandleFunc("/api/v1/getSprite", middleware.ApplyMiddleware(apis.GetSprite, middleware.EnableCORS))
	r.HandleFunc("/api/v1/uploadSprite", middleware.ApplyMiddleware(apis.UploadSprite, middleware.EnableCORS))

	r.HandleFunc("/api/v1/getMaps", middleware.ApplyMiddleware(apis.GetMaps, middleware.EnableCORS))
	r.HandleFunc("/api/v1/getSpaces", middleware.ApplyMiddleware(apis.GetSpaces, middleware.EnableCORS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server Listening at Port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
