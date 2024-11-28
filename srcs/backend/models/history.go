package models

import "gorm.io/gorm"

type GameHistory struct {
    gorm.Model 
    Player1ID uint64   `json:"player1_id"`  
    Player2ID uint64   `json:"player2_id"`
    WinnerID  uint64   `json:"winner_id"`
    Score1    int      `json:"score1"`
    Score2    int      `json:"score2"`
}