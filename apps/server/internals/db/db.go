package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID       string  `gorm:"primaryKey;unique;default:(uuid())"` // Use (uuid()) or equivalent for cuid
	Username string  `gorm:"unique;not null"`
	Password string  `gorm:"unique;not null"`
	AvatarID *string `gorm:"default:null"` // Nullable field
}

// Space model
type Space struct {
	ID        string  `gorm:"primaryKey;unique;default:(uuid())"`
	Name      string  `gorm:"not null"`
	Width     int     `gorm:"not null"`
	Height    *int    `gorm:"default:null"` // Nullable field
	Thumbnail *string `gorm:"default:null"` // Nullable field
}

// SpaceElement model
type SpaceElement struct {
	ID        string `gorm:"primaryKey;unique;default:(uuid())"`
	ElementID string `gorm:"not null"`
	SpaceID   string `gorm:"not null"`
	X         int    `gorm:"not null"`
	Y         int    `gorm:"not null"`
}

// Element model
type Element struct {
	ID       string `gorm:"primaryKey;unique;default:(uuid())"`
	Width    int    `gorm:"not null"`
	Height   int    `gorm:"not null"`
	ImageURL string `gorm:"not null"`
}

// Map model
type Map struct {
	ID     string `gorm:"primaryKey;unique;default:(uuid())"`
	Width  int    `gorm:"not null"`
	Height int    `gorm:"not null"`
	Name   string `gorm:"not null"`
}

// MapElement model
type MapElement struct {
	ID        string  `gorm:"primaryKey;unique;default:(uuid())"`
	MapID     string  `gorm:"not null"`
	ElementID *string `gorm:"default:null"` // Nullable field
	X         *int    `gorm:"default:null"` // Nullable field
	Y         *int    `gorm:"default:null"` // Nullable field
}

// Avatar model
type Avatar struct {
	ID       string  `gorm:"primaryKey;unique;default:(uuid())"`
	ImageURL *string `gorm:"default:null"` // Nullable field
	Name     *string `gorm:"default:null"` // Nullable field
}

// GORM requires an `AutoMigrate` function to initialize models
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Space{},
		&SpaceElement{},
		&Element{},
		&Map{},
		&MapElement{},
		&Avatar{},
	)
}

var Database *gorm.DB

func ConnectDB() error {
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(
		mysql.Open(dbUrl),
		&gorm.Config{},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	err = Migrate(db)
	if err != nil {
		log.Println(err)
		return err
	}
	Database = db
	return nil
}
