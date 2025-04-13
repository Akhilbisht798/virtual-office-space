package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akhilbisht798/ws/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	http.HandleFunc("/ws", internal.WebSocketHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server is running on port: ", port)
	http.ListenAndServe(":"+port, nil)
}
