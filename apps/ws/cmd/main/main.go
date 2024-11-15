package main

import (
	"log"
	"net/http"

	"github.com/akhilbisht798/ws/internal"
)

func main() {
	http.HandleFunc("/ws", internal.WebSocketHandler)
	log.Println("Server is running on port 5000")
	http.ListenAndServe(":5000", nil)
}
