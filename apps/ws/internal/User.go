package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type UserConn struct {
	conn     *websocket.Conn
	Id       string  `json:"id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Sprite   string  `json:"sprite"`
	Username string  `json:"username"`
}

type GetUserRequest struct {
	UserID string `json:"userID"`
}

type GetUserResponse struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
}

func GetUser(userId string) (string, error) {
	serverUrl := os.Getenv("SERVER")
	if serverUrl == "" {
		return "", errors.New("server url not found in the env")
	}
	url := serverUrl + "/api/v1/getUser"
	payload := GetUserRequest{
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

	var data GetUserResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	//fmt.Println("Sprite is: ", data.Sprite)
	return data.Username, nil
}
