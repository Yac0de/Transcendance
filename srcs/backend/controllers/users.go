package controllers

import (
	"api/database"
	"api/models"
	// "log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Users(ctx *gin.RouterGroup) {
    ctx.GET("", GetAllUsers)
    ctx.POST("", CreateUser)

    ctx.GET("/:userId", GetUser)
    ctx.PUT("/:userId", UpdateUser)
    ctx.DELETE("/:userId", DeleteUser)
}

func GetAllUsers(ctx *gin.Context) {
    var users []models.User
    if err := database.DB.Find(&users).Error; err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, users)
}

func CreateUser(ctx *gin.Context) {
    var input models.CreateUserDto

    err := ctx.ShouldBindJSON(&input);
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data" })
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newUser := models.User{
        Nickname: input.Nickname,
        Email: input.Email,
        Password: string(hashedPassword),
    }

    if err := database.DB.Create(&newUser).Error; err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx.Status(http.StatusAccepted)
}

func GetUser(ctx *gin.Context) {
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

    ctx.JSON(http.StatusOK, target)
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

    err = ctx.ShouldBindJSON(&input);
    if err != nil || reflect.DeepEqual(input, (models.UpdateUserDto{})) {
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data" })
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
        ctx.JSON(500, gin.H{"error": err.Error()})
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

