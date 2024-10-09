package controllers

import (
	"api/database"
	"api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FriendShip(ctx *gin.RouterGroup) {
	ctx.GET("/list", GetFriendList)
	ctx.GET("/requests", GetFriendRequests)
	ctx.POST("/add/:friendId", AddFriend)
	ctx.POST("/accept/:friendId", AcceptFriend)
}

func AcceptFriend(ctx *gin.Context) {
	friend := ctx.Param("friendId")
	friendId, err := strconv.ParseUint(friend, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid format of friend id"})
		return
	}

	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var friendship models.FriendShip
	if err := database.DB.First(&friendship, "user_id = ? AND friend_id = ?", friendId, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Friendship does not exists"})
		return
	}

	if friendship.MutualFriends == true {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Friendship already accpeted"})
		return
	}

	friendship.MutualFriends = true

	if err := database.DB.Save(&friendship).Error; err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Succes": "Friend accepted"})
}

func AddFriend(ctx *gin.Context) {
	friend := ctx.Param("friendId")
	friendId, err := strconv.ParseUint(friend, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid format of friend id"})
		return
	}

	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var friendUser models.User
	if err := database.DB.First(&friendUser, friendId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	var existingFriend models.FriendShip
	if err := database.DB.First(&existingFriend, "(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", id, friendId, friendId, id).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Friendship already exists"})
		return
	}

	friends := models.FriendShip{
		UserID:        id,
		FriendID:      friendUser.ID,
		MutualFriends: false,
	}

	if err := database.DB.Create(&friends).Error; err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Succes": "Request sent"})
}

func GetFriendList(ctx *gin.Context) {
	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var friendsRequested []models.UserResponse

	// Search friends user by friend.user_id field
	err := database.DB.Raw(`
		SELECT u.id, u.display_name, u.nickname, u.avatar 
		FROM friend_ships f JOIN "users" u ON f.friend_id = u.id
		WHERE f.user_id = ? AND f.mutual_friends = true
		`, id).Scan(&friendsRequested).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Search friends user by friend.friend_id field
	var friendsAccepted []models.UserResponse
	err = database.DB.Raw(`
		SELECT u.id, u.display_name, u.nickname, u.avatar 
		FROM friend_ships f JOIN "users" u ON f.user_id = u.id
		WHERE f.friend_id = ? AND f.mutual_friends = true
		`, id).Scan(&friendsAccepted).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	friends := append(friendsRequested, friendsAccepted...)
	ctx.JSON(http.StatusOK, friends)
}

func GetFriendRequests(ctx *gin.Context) {
	userId, exists := ctx.Get("UserId")
	id, ok := userId.(uint)
	if exists == false || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var requests []models.UserResponse
	err := database.DB.Raw(
		`
		SELECT u.id, u.display_name, u.nickname, u.avatar
		FROM friend_ships f JOIN "users" u ON f.user_id = u.id
		WHERE f.friend_id = ? AND f.mutual_friends = false
		`, id).Scan(&requests).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}
