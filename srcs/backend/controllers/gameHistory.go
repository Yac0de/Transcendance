package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"api/models"
	"api/database"
    "api/prometheus"
    "gorm.io/gorm"
	"fmt"
)


type GameHistoryInput struct {
    Player1ID uint64 `json:"player1_id"`
    Player2ID uint64 `json:"player2_id"`
    WinnerID  uint64 `json:"winner_id"`
    Score1    int    `json:"Score1"`
    Score2    int    `json:"Score2"`
}

func SaveGameHistory(c *gin.Context) {

    var input GameHistoryInput
    if err := c.BindJSON(&input); err != nil {
        fmt.Printf("Binding error: %v\n", err)
        fmt.Printf("Partial input: %+v\n", input)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input data",
            "details": err.Error(),
        })
        return
    }

    fmt.Printf("Successfully bound input: %+v\n", input)

    gameHistory := models.GameHistory{
        Player1ID: input.Player1ID,
        Player2ID: input.Player2ID,
        WinnerID:  input.WinnerID,
        Score1:    input.Score1,
        Score2:    input.Score2,
    }

    if err := database.DB.Create(&gameHistory).Error; err != nil {
        fmt.Printf("DB error: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save game history",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Game history saved successfully",
        "data": gameHistory,
    })
    prometheus.IncrementPlayedGames()
}


func GetUserGameHistory(c *gin.Context) {
    nickname := c.Param("nickname")
    
    var user models.User
    if err := database.DB.Where("nickname = ?", nickname).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "User not found",
                "details": fmt.Sprintf("No user found with nickname: %s", nickname),
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch user",
            "details": err.Error(),
        })
        return
    }

    var gameHistories []models.GameHistory
    if err := database.DB.
        Where("player1_id = ? OR player2_id = ?", user.ID, user.ID).
        Preload("Player1", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "display_name", "nickname", "avatar")
        }).
        Preload("Player2", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "display_name", "nickname", "avatar")
        }).
        Preload("Winner", func(db *gorm.DB) *gorm.DB {
            return db.Select("id", "display_name", "nickname", "avatar")
        }).
        Order("created_at desc").
        Find(&gameHistories).Error; err != nil {
            fmt.Printf("Database error: %v\n", err) 
            c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch game history",
            "details": err.Error(),
        })
        return
    }

    type GameHistoryResponse struct {
        models.GameHistory
        IsWinner bool `json:"is_winner"`
    }

    var response []GameHistoryResponse
    for _, game := range gameHistories {
        response = append(response, GameHistoryResponse{
            GameHistory: game,
            IsWinner:    game.WinnerID == uint64(user.ID),
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "data": response,
    })
}