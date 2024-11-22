package main

import (
	"log"
	"net/http"

	"github.com/Akhilbisht798/office/server/internals/apis"
	"github.com/Akhilbisht798/office/server/internals/db"
	"github.com/gorilla/mux"
)

func main() {
	err := db.ConnectDB()
	if err != nil {
		log.Println("server closing.")
		return
	}
	r := mux.NewRouter()
	r.HandleFunc("/", apis.RootHandler)
	log.Println("Server Listening at Port: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
