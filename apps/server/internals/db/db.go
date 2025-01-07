package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID       string  `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"` // Use (uuid()) or equivalent for cuid
	Username string  `gorm:"unique;not null"`
	Password string  `gorm:"unique;not null"`
	AvatarID *string `gorm:"default:null"` // Nullable field
	Spaces   []Space `gorm:"foreignKey:UserID"`
	Avatar   Avatar  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Maps from which spaces will be created
type Map struct {
	ID     string  `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"`
	Name   string  `gorm:"type:varchar(255);not null;unique"`
	Spaces []Space `gorm:"foreignKey:MapID"`
}

// Space or room created by a certain user to join to.
type Space struct {
	ID        string  `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"`
	Name      string  `gorm:"not null"`
	Thumbnail *string `gorm:"default:null"` // Nullable field
	MapID     string  `gorm:"type:varchar(255);not null"`
	UserID    string  `gorm:"type:varchar(255);not null"`
	Public    bool    `gorm:"not null"`
}

// Avatar
type Avatar struct {
	ID   string `gorm:"type:varchar(255);primaryKey;unique;default:(uuid())"`
	Name string `gorm:"unique;not null"` // Nullable field
}

// GORM requires an `AutoMigrate` function to initialize models
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Map{},
		&Space{},
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
