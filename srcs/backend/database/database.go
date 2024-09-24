package database

import (
	"api/models"
	"log"
	"os"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func New() {
    DB = initDB()
    createMockUsers()
}

func initDB() *gorm.DB  {
     dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRES_DB"),
        "5432",
    )    

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalln(err)
    }
    database.AutoMigrate(&models.User{})
    
    return database
} 

func createMockUsers() {
    users := []models.User{
        {Nickname: "Hichame", Email: "hichame@42LH.fr", Password: "hichame42LH"},
        {Nickname: "Maxime", Email: "maxime@42LH.fr", Password: "maxime42LH"},
        {Nickname: "Yanis", Email: "yanis@42LH.fr", Password: "yanis42LH"},
        {Nickname: "Omar", Email: "omar@42LH.fr", Password: "omar42LH"},
    }

    for i := range users {
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(users[i].Password), 10)
        users[i].Password = string(hashedPassword)
        if err := DB.Create(&users[i]).Error; err != nil {
            log.Printf("Failed to create user %s: %v", users[i].Nickname, err)
        } else {
            log.Printf("Created mock user: %s", users[i].Nickname)
        }
    }
}
