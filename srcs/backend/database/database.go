package database

import (
	"api/models"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func New() {
	DB = initDB()
	createMockUsers()
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		"5432",
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Printf("ERR DOCKER")
		log.Fatalln(err)
	}

	database.AutoMigrate(&models.User{})

	return database
}

func createMockUsers() {
	users := []models.User{
		{Nickname: "Hichame", DisplayName: "hichame", Password: "hichame42LH"},
		{Nickname: "Maxime", DisplayName: "maxime", Password: "maxime42LH"},
		{Nickname: "Yanis", DisplayName: "yanis", Password: "yanis42LH"},
		{Nickname: "Omar", DisplayName: "omar", Password: "omar42LH"},
	}

	for i := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(users[i].Password), 10)
		users[i].Nickname = strings.ToLower(users[i].Nickname)
		users[i].Password = string(hashedPassword)
		if err := DB.Create(&users[i]).Error; err != nil {
			log.Printf("Failed to create user %s: %v", users[i].Nickname, err)
		} else {
			log.Printf("Created mock user: %s", users[i].Nickname)
		}
	}
}
