package internal

import (
	"log"
	"os"

	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserConn struct {
	conn   *websocket.Conn
	Id     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Sprite string  `json:"sprite"`
}

type User struct {
	ID       string  `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"` // Use (uuid()) or equivalent for cuid
	Username string  `gorm:"unique;not null"`
	Password string  `gorm:"unique;not null"`
	AvatarID *string `gorm:"default:null"` // Nullable field
}

func GetUser(userId string) (User, error) {
	var user User
	db, err := ConnectDB()
	if err != nil {
		log.Println(err)
		return user, err
	}
	res := db.Where("id = ?", userId).First(&user)
	if res.Error != nil {
		log.Println(res.Error)
		return user, res.Error
	}
	return user, nil
}

func ConnectDB() (*gorm.DB, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(
		mysql.Open(dbUrl),
		&gorm.Config{},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

// User model
//type User struct {
//	ID       string  `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"` // Use (uuid()) or equivalent for cuid
//	Username string  `gorm:"unique;not null"`
//	Password string  `gorm:"unique;not null"`
//	AvatarID *string `gorm:"default:null"` // Nullable field
//	Spaces   []Space `gorm:"foreignKey:UserID"`
//}
