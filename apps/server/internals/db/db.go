package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID       string  `gorm:"primaryKey;unique;default:uuid_generate_v4()"` // Use uuid_generate_v4() or equivalent for cuid
	Username string  `gorm:"unique;not null"`
	Password string  `gorm:"unique;not null"`
	AvatarID *string `gorm:"default:null"` // Nullable field
	Role     Role    `gorm:"type:role"`    // Enum reference
}

// Space model
type Space struct {
	ID        string  `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	Name      string  `gorm:"not null"`
	Width     int     `gorm:"not null"`
	Height    *int    `gorm:"default:null"` // Nullable field
	Thumbnail *string `gorm:"default:null"` // Nullable field
}

// SpaceElement model
type SpaceElement struct {
	ID        string `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	ElementID string `gorm:"not null"`
	SpaceID   string `gorm:"not null"`
	X         int    `gorm:"not null"`
	Y         int    `gorm:"not null"`
}

// Element model
type Element struct {
	ID       string `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	Width    int    `gorm:"not null"`
	Height   int    `gorm:"not null"`
	ImageURL string `gorm:"not null"`
}

// Map model
type Map struct {
	ID     string `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	Width  int    `gorm:"not null"`
	Height int    `gorm:"not null"`
	Name   string `gorm:"not null"`
}

// MapElement model
type MapElement struct {
	ID        string  `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	MapID     string  `gorm:"not null"`
	ElementID *string `gorm:"default:null"` // Nullable field
	X         *int    `gorm:"default:null"` // Nullable field
	Y         *int    `gorm:"default:null"` // Nullable field
}

// Avatar model
type Avatar struct {
	ID       string  `gorm:"primaryKey;unique;default:uuid_generate_v4()"`
	ImageURL *string `gorm:"default:null"` // Nullable field
	Name     *string `gorm:"default:null"` // Nullable field
}

// Role enum type
type Role string

const (
	Admin  Role = "Admin"
	Client Role = "User"
)

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

func ConnectDB() error {
	dbUrl := os.Getenv("DB_URL")
	db, err := gorm.Open(
		postgres.Open(dbUrl),
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
	return nil
}
