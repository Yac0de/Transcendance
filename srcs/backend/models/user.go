package models

// import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primary_key;autoIncrement"`
	Nickname string `json:"nickname" gorm:"unique;not null" binding:"required,min=3"`
	Email    string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password string `json:"password" gorm:"not null" binding:"required,min=6"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type CreateUserDto struct {
	Nickname string `json:"nickname" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserDto struct {
	Nickname string `json:"nickname,omitempty" binding:"omitempty,min=3"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type SignInDto struct {
	Nickname string `json:"nickname" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
