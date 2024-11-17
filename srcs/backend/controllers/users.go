package controllers

import (
	"api/database"
	"api/models"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func Users(ctx *gin.RouterGroup) {
	ctx.GET("/all", GetAllUsers)
	ctx.GET("/", func(ctx *gin.Context) {
		if ctx.Query("nickname") != "" {
			GetUserByNickname(ctx)
		} else {
			GetUserByID(ctx)
		}
	})
	ctx.GET("", GetUser)

	ctx.PUT("/update-profile", UpdateProfile)
	ctx.PUT("/change-password", ChangePassword)

	FriendShip(ctx.Group("/friendships")) // /users/friends/...
	ctx.DELETE("/delete-account", DeleteAccount)
}

func GetAllUsers(ctx *gin.Context) {
	var users []models.UserResponse

	if err := database.DB.Raw("SELECT id, nickname, display_name FROM users").Scan(&users).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var user models.UserResponse
	result := database.DB.Raw("SELECT id, display_name, nickname, avatar FROM users WHERE id = ?", id).Scan(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func GetUserByID(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.UserResponse
	result := database.DB.Raw("SELECT id, display_name, nickname, avatar FROM users WHERE id = ?", id).Scan(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func GetUserByNickname(ctx *gin.Context) {
	nickname := ctx.Query("nickname")
	if nickname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.UserResponse
	result := database.DB.Raw("SELECT id, display_name, nickname, avatar FROM users WHERE nickname = ?", nickname).Scan(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func UpdateProfile(ctx *gin.Context) {
	id, exists := ctx.Get("UserId")

	userId, ok := id.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	code, err := UpdateUser(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nickname already exists, please choose another one."})
			return
		}
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "user profile updated successfully",
	})
}

func UpdateUser(ctx *gin.Context, id uint) (int, error) {
	validate := validator.New()
	var user models.User

	result := database.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return http.StatusNotFound, result.Error
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return http.StatusNotFound, err
	}

	for key := range form.Value {
		switch strings.ToLower(key) {
		case "nickname":
			user.Nickname = strings.ToLower(form.Value[key][0])
		case "displayname":
			user.DisplayName = form.Value[key][0]
			// Handle password case
		}
	}

	err = validate.Struct(user)
	if err != nil {
		var e string
		for _, err := range err.(validator.ValidationErrors) {
			e += fmt.Sprintf("Error validation for '%s': %s\n", err.Field(), err.Tag())
		}
		return http.StatusBadRequest, fmt.Errorf(e)
	}

	if form.File["avatar"] != nil {
		filename, err := SaveAvatar(ctx, form.File["avatar"][0])
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if len(user.Avatar) > 0 {
			err = os.Remove(filepath.Join("./avatars", user.Avatar))
		}
		user.Avatar = filename
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}

func SaveAvatar(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("Avatar not found")
	}
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	file.Filename = filename

	err := ctx.SaveUploadedFile(file, filepath.Join("./avatars", file.Filename))
	if err != nil {
		return "", err
	}

	return filename, nil
}

func DeleteAccount(ctx *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, exists := ctx.Get("UserId")
	userId, ok := id.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func ChangePassword(ctx *gin.Context) {
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id, exists := ctx.Get("UserId")
	userId, ok := id.(uint)
	if !exists || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	if len(req.CurrentPassword) < 6 || len(req.NewPassword) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password length (6 character min)"})
		return

	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password doesn't match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash the new password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Password updated successfully"})
}
