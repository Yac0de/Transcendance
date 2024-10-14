package models

type User struct {
	ID          uint    `json:"id" gorm:"primary_key;autoIncrement"`
	DisplayName string  `json:"displayname" gorm:"not null" binding:"required,min=3" validate:"required,min=3"`
	Nickname    string  `json:"nickname" gorm:"unique;not null" binding:"required,min=3" validate:"required,min=3"`
	Password    string  `json:"password" gorm:"not null" binding:"required,min=6" validate:"required,min=6"`
	Avatar      string  `json:"avatar"`
	Friends     []*User `gorm:"many2many:friendShip;"`
}

type FriendShip struct {
	UserID        uint `gorm:"primaryKey"`
	FriendID      uint `gorm:"primaryKey"`
	MutualFriends bool `gorm:"not null"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"displayname"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
}

type CreateUserDto struct {
	Nickname string `json:"nickname" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserDto struct {
	DisplayName string `json:"displayname,omitempty" binding:"omitempty,min=3"`
	Nickname    string `json:"nickname,omitempty" binding:"omitempty,min=3"`
	Password    string `json:"password,omitempty" binding:"omitempty,min=6"`
	Avatar      string `json:"avatar" binding:"omitempty"`
}

type SignInDto struct {
	Nickname string `json:"nickname" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
