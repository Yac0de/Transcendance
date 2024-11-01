package models

type User struct {
	ID          uint    `json:"id" gorm:"primary_key;autoIncrement"`
	DisplayName string  `json:"displayname" gorm:"not null" binding:"required,min=3" validate:"required,min=3,max=30"`
	Nickname    string  `json:"nickname" gorm:"unique;not null" binding:"required,min=3" validate:"required,min=3,max=20"`
	Password    string  `json:"password" gorm:"not null" binding:"required,min=6" validate:"required,min=6,max=50"`
	Avatar      string  `json:"avatar"`
	Friends     []*User `gorm:"many2many:friendShip;"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"displayname"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
}

type CreateUserDto struct {
	Nickname string `json:"nickname" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type UpdateUserDto struct {
	DisplayName string `json:"displayname,omitempty" binding:"omitempty,min=3,max=30"`
	Nickname    string `json:"nickname,omitempty" binding:"omitempty,min=3,max=20"`
	Password    string `json:"password,omitempty" binding:"omitempty,min=6,max=50"`
	Avatar      string `json:"avatar" binding:"omitempty"`
}

type SignInDto struct {
	Nickname string `json:"nickname" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}
