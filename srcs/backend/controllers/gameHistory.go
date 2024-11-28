package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"api/models"
	"strconv"
	"api/database"
    "io"
    "bytes"
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
    // Log le JSON brut reçu
    rawData, _ := io.ReadAll(c.Request.Body)
    fmt.Printf("Raw request body: %s\n", string(rawData))
    c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

    var input GameHistoryInput
    if err := c.ShouldBindJSON(&input); err != nil {
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
}


func GetUserGameHistory(c *gin.Context) {
    userID := c.Param("userId")
    fmt.Printf("Fetching games for user ID: %s\n", userID)
    
    // Vérifier si id est valide
    id, err := strconv.ParseUint(userID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }

    var gameHistories []models.GameHistory
    
    // Utiliser directement database.DB
    if err := database.DB.Where("player1_id = ? OR player2_id = ?", id, id).
    Order("created_at desc").
    Find(&gameHistories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch game history",
            "details": err.Error(),
        })
        return
    }

    // Préparer la réponse avec is_winner
    type GameHistoryResponse struct {
        models.GameHistory
        IsWinner bool `json:"is_winner"`
    }

    var response []GameHistoryResponse
    for _, game := range gameHistories {
        response = append(response, GameHistoryResponse{
            GameHistory: game,
            IsWinner:    game.WinnerID == id,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "data": response,
    })
}