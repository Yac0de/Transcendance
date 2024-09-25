package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Auth(ctx *gin.RouterGroup) {
	ctx.POST("/signin", SignIn)
	ctx.POST("/signup", SignUp)
}

func SignUp(ctx *gin.Context) {
	var input models.CreateUserDto

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.CreateToken(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("access_token", token, 0, "/", "", false, true)

	ctx.JSON(http.StatusCreated, gin.H{"succes": "User created"})
}

func SignIn(ctx *gin.Context) {
	var input models.SignInDto

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	var user models.User
	result := database.DB.Raw("SELECT * FROM users WHERE nickname ILIKE ? LIMIT 1", input.Nickname).Scan(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("access_token", token, 0, "/", "", false, true)

	// ctx.Status(http.StatusNoContent)
	ctx.JSON(http.StatusAccepted, gin.H{"succes": "User connected"})
}
