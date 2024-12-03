package database

import (
	"api/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func New() {
	DB = initDB()
	CreateMockUsers()
	CreateMockConversation()
	CreateMockGames()
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

	database.AutoMigrate(&models.User{}, &models.TwoFactorAuth{}, &models.FriendShip{}, &models.Message{}, &models.GameHistory{})

	return database
}

func CreateMockUsers() {
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
	AddFriendsToUsers()
}

func AddFriendsToUsers() {
	friendships := map[string][]string{
		"hichame": {"maxime", "omar"},
		"yanis":   {"omar"},
	}

	for nickname, friends := range friendships {
		var user models.User

		if err := DB.First(&user, "nickname = ?", nickname).Error; err != nil {
			log.Printf("Failed to get user %s to add friends: %v", nickname, err)
			continue
		}

		for _, friendNickname := range friends {
			var friend models.User
			if err := DB.First(&friend, "nickname = ?", friendNickname).Error; err != nil {
				log.Printf("Failed to get user %s to add in %s friends list: %v", friendNickname, user.Nickname, err.Error())
				continue
			}

			friendsRelation := models.FriendShip{
				UserID:        user.ID,
				FriendID:      friend.ID,
				MutualFriends: true,
			}

			if err := DB.Create(&friendsRelation).Error; err != nil {
				log.Printf("Failed to add friend %s to user %s: %v", friendNickname, user.Nickname, err.Error())
			} else {
				log.Printf("Added friend %s to user %s", friendNickname, nickname)
			}
		}
	}

}


func GetMockedMessages() []models.Message {
	return []models.Message{{
		SenderID:   1,
		ReceiverID: 4,
		Content:    "Salut ID 4 !",
		CreatedAt:  time.Now().Add(time.Hour * (-4)),
	}, {
		SenderID:   4,
		ReceiverID: 1,
		Content:    "Hey salut ID 1",
		CreatedAt:  time.Now().Add(time.Hour * (-3)),
	}, {
		SenderID:   4,
		ReceiverID: 1,
		Content:    "Comment vas tu ?",
		CreatedAt:  time.Now().Add(time.Hour*(-3) + time.Minute*2),
	}, {
		SenderID:   1,
		ReceiverID: 4,
		Content:    "Bien et toi ?",
		CreatedAt:  time.Now().Add(time.Hour * (-2)),
	}, {
		SenderID:   1,
		ReceiverID: 4,
		Content:    "Super !",
		CreatedAt:  time.Now().Add(time.Hour * (-1)),
	}}
}

func CreateMockConversation() {
	conversation := GetMockedMessages()
	for _, message := range conversation {
		if err := DB.Create(&message).Error; err != nil {
			log.Printf("Failed to add message %s from user id %d to user id %d: %v", message.Content, message.SenderID, message.ReceiverID, err.Error())
		}
	}
}

func CreateMockGames() {
	games := []models.GameHistory{
		{
			Player1ID: 1,
			Player2ID: 2,
			WinnerID:  1,
			Score1:    5,
			Score2:    3,
		},
		{
			Player1ID: 1,
			Player2ID: 4,
			WinnerID:  4,
			Score1:    2,
			Score2:    5,
		},
	}

	for _, game := range games {
		if err := DB.Create(&game).Error; err != nil {
			log.Printf("Failed to create game history: %v", err)
		} else {
			log.Printf("Created mock game between players %d and %d", game.Player1ID, game.Player2ID)
		}
	}
}
