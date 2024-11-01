package models

import "time"

type FriendShip struct {
	UserID        uint `gorm:"primaryKey"`
	FriendID      uint `gorm:"primaryKey"`
	MutualFriends bool `gorm:"not null"`
}

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SenderID   uint      `json:"senderId" gorm:"not null"`
	ReceiverID uint      `json:"receiverId" gorm:"not null"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
}
