package main

import (
	"log"
	"net/http"

	"github.com/akhilbisht798/ws/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	http.HandleFunc("/ws", internal.WebSocketHandler)
	log.Println("Server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
