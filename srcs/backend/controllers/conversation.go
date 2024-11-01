package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Conversation(ctx *gin.RouterGroup) {
	ctx.GET("/:friendId", GetConversation)
	ctx.POST("/add", SaveNewMessage)
}

func GetConversation(ctx *gin.Context) {
	friend := ctx.Param("friendId")
	friendId, err := strconv.ParseUint(friend, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format of friend id"})
		return
	}

	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var conversation []models.Message
	if err := database.DB.Find(&conversation, "(sender_id = ? OR sender_id = ?) AND (receiver_id = ? OR receiver_id = ?)", id, friendId, id, friendId).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"conversation": conversation})
}

func SaveNewMessage(ctx *gin.Context) {
	var message models.Message
	err := ctx.ShouldBindJSON(&message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := database.DB.Create(&message).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}
