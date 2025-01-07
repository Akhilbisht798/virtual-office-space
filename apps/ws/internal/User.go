package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type UserConn struct {
	conn   *websocket.Conn
	Id     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Sprite string  `json:"sprite"`
}

type GetUserSpriteRequest struct {
	UserID string `json:"userID"`
}

type GetSpriteResponse struct {
	Sprite string `json:"sprite"`
}

func GetUser(userId string) (string, error) {
	serverUrl := os.Getenv("SERVER")
	if serverUrl == "" {
		return "", errors.New("server url not found in the env")
	}
	url := serverUrl + "/api/v1/getSprite"
	payload := GetUserSpriteRequest{
		UserID: userId,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	defer resp.Body.Close()

	var data GetSpriteResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	fmt.Println("Sprite is: ", data.Sprite)
	return data.Sprite, nil
}
