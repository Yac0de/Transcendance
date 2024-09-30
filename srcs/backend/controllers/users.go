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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Users(ctx *gin.RouterGroup) {
	ctx.GET("", GetUser)
	ctx.GET("/all", GetAllUsers)
	ctx.GET("/:userId", GetUser)

	ctx.PUT("/update-profile", UpdateProfile)
	ctx.POST("/upload-avatar", UploadAvatar)
	ctx.PUT("/:userId", UpdateUser)
	ctx.DELETE("/:userId", DeleteUser)
}

func GetAllUsers(ctx *gin.Context) {
	var users []models.UserResponse

	if err := database.DB.Raw("SELECT id, nickname, email FROM users").Scan(&users).Error; err != nil {
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
	result := database.DB.Raw("SELECT id, nickname, email, avatar FROM users WHERE id = ?", id).Scan(&user)
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

	var user models.User
	result := database.DB.First(&user, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for key := range form.Value {
		// Unsecure (any data validation)
		switch strings.ToLower(key) {
		case "nickname":
			user.Nickname = form.Value[key][0]
		case "email":
			user.Email = form.Value[key][0]
			// Handle password case
		}
	}

	if form.File["avatar"] != nil {
		filename, err := SaveAvatar(ctx, form.File["avatar"][0])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(user.Avatar) > 0 {
			err = os.Remove(filepath.Join("./avatars", user.Avatar))
			fmt.Println("Can not remove file: ", user.Avatar, err)
		}
		user.Avatar = filename
	}

	if err := database.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "user profile updated successfully",
	})
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

func UploadAvatar(ctx *gin.Context) {
	id, exists := ctx.Get("UserId")
	userId, ok := id.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	timestamp := time.Now().Unix()
	newFilename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	file.Filename = newFilename
	fmt.Println(userId)

	err = ctx.SaveUploadedFile(file, filepath.Join("./avatars", file.Filename))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	result := database.DB.Exec("UPDATE users SET avatar = ? WHERE ID = ?", file.Filename, userId)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Can't update table"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Files uploaded successfully"})
}

func UpdateUser(ctx *gin.Context) {
	id, err := GetUserIdToUINT(ctx.Params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var target models.User
	if err := database.DB.First(&target, id).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input models.UpdateUserDto

	err = ctx.ShouldBindJSON(&input)
	if err != nil || reflect.DeepEqual(input, (models.UpdateUserDto{})) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		}
		return
	}

	if input.Email != "" && input.Email != target.Email {
		target.Email = input.Email
	}

	if input.Nickname != "" && input.Nickname != target.Nickname {
		target.Nickname = input.Nickname
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		target.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&target).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusAccepted)
}

func DeleteUser(ctx *gin.Context) {
	id, err := GetUserIdToUINT(ctx.Params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Delete(&models.User{}, id)

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func GetUserIdToUINT(params gin.Params) (uint, error) {
	idStr := params.ByName("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
