package controllers

import (
	"api/database"
	"api/models"
	"api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Auth(ctx *gin.RouterGroup) {
	ctx.POST("/signin", SignIn)
	ctx.POST("/signup", SignUp)
	ctx.POST("/signout", SignOut)
}

func SignOut(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", 0, "/", "", false, true)
	ctx.JSON(http.StatusCreated, gin.H{"succes": "Logout successfully"})
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

	var existingUser models.User
	if err := database.DB.Where("nickname = ?", strings.ToLower(input.Nickname)).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Nickname already taken"})
		return
	}

	newUser := models.User{
		Nickname:    strings.ToLower(input.Nickname),
		DisplayName: input.Nickname,
		Password:    string(hashedPassword),
		Avatar:      "",
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
			"error": "Invalid input data. Please check the fields and try again.",
		})
		return
	}

	var user models.User
	result := database.DB.Where("nickname = ?", strings.ToLower(input.Nickname)).First(&user)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist. Please check your nickname."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while checking the user. Please try again later."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password. Please try again."})
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while generating the token. Please try again."})
		return
	}

	ctx.SetCookie("access_token", token, 0, "/", "", false, true)
	ctx.JSON(http.StatusAccepted, gin.H{"success": "User connected"})
}

