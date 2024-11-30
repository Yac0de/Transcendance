package models

import "gorm.io/gorm"

type GameHistory struct {
    gorm.Model 
    Player1ID uint64   `json:"player1_id" gorm:"not null"`  
    Player2ID uint64   `json:"player2_id" gorm:"not null"`
    WinnerID  uint64   `json:"winner_id" gorm:"not null"`
    Score1    int      `json:"score1"`
    Score2    int      `json:"score2"`

    Player1 User `json:"player1" gorm:"foreignKey:Player1ID"`
    Player2 User `json:"player2" gorm:"foreignKey:Player2ID"`
    Winner User  `json:"winner" gorm:"foreignKey:WinnerID"`
}