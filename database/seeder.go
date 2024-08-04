package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	*gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique;not null"`
}

func Seed() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	DB = db

}
