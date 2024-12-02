package controllers

import (
	"api/database"
	"api/models"
	"api/prometheus"
	"api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Auth(ctx *gin.RouterGroup) {
	ctx.POST("/signin", SignIn)
	ctx.POST("/signup", SignUp)
	ctx.POST("/signout", SignOut)
	ctx.GET("/generate2FA", Generate2FAcode)
	ctx.GET("/2FA-status", GetUser2FAStatus)
	ctx.POST("/verify2FA", Verify2FAcode)
}

func GetUser2FAStatus(ctx *gin.Context) {
	id, exists := ctx.get("userid")

	userid, ok := id.(uint)
	if !exists || !ok {
		ctx.json(http.statusunauthorized, gin.h{"error": "unauthorized: you must be logged in to access this resource."})
		return
	}

	var twoFactor models.TwoFactorAuth
	if err := database.DB.Where("user_id = ?", userId).First(&twoFactor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.json(http.StatusOK, gin.h{"status": false})
			return
		}
	}

	if twoFactor.IsActive == true {
		ctx.json(http.StatusOK, gin.h{"status": true})
		return
	}

	ctx.json(http.StatusOK, gin.h{"status": false})
	return
}

func Verify2FAcode(ctx *gin.Context) {
	id, exists := ctx.get("userid")

	userid, ok := id.(uint)
	if !exists || !ok {
		ctx.json(http.statusunauthorized, gin.h{"error": "unauthorized: you must be logged in to access this resource."})
		return
	}

	var data map[string]interface{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	code, exists := data["code"].(string)
	if !exists || len(code) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "You must provide the 2FA code"})
		return
	}

	var twoFactor models.TwoFactorAuth
	if err := database.DB.Where("user_id = ?", userId).First(&twoFactor).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "2FA code was not generated for this user"})
		return
	}

	valid := totp.Validate(code, twoFactor.Secret)
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "2FA code verified successfully"})
}

func Generate2FAcode(ctx *gin.Context) {
	id, exists := ctx.Get("UserId")

	userId, ok := id.(uint)
	if !exists || !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: You must be logged in to access this resource."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Transcendance",
		AccountName: user.Nickname,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	var twoFactor models.TwoFactorAuth
	if err := database.DB.Where("user_id = ?", userId).First(&twoFactor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			twoFactor = models.TwoFactorAuth{
				UserID:   userId,
				IsActive: false,
				Secret:   key.Secret(),
			}
			if err := database.DB.Create(&twoFactor).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save TwoFactorAuth"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query TwoFactorAuth"})
			return
		}
	} else {
		twoFactor.Secret = key.Secret()
		if err := database.DB.Save(&twoFactor).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update TwoFactorAuth"})
			return
		}
	}

	ctx.Data(http.StatusCreated, "image/png", png)
}

func SignOut(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", 0, "/", "", false, true)
	ctx.JSON(http.StatusCreated, gin.H{"succes": "Logout successfully"})
	prometheus.DecrementActiveUsers()
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
		prometheus.RecordLoginAttempt(false)
		return
	}

	var user models.User
	result := database.DB.Where("nickname = ?", strings.ToLower(input.Nickname)).First(&user)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist. Please check your nickname."})
			prometheus.RecordLoginAttempt(false)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while checking the user. Please try again later."})
		prometheus.RecordLoginAttempt(false)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password. Please try again."})
		prometheus.RecordLoginAttempt(false)
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while generating the token. Please try again."})
		prometheus.RecordLoginAttempt(false)
		return
	}

	ctx.SetCookie("access_token", token, 0, "/", "", false, true)
	ctx.JSON(http.StatusAccepted, gin.H{"success": "User connected"})
	prometheus.RecordLoginAttempt(true)
	prometheus.IncrementActiveUsers()
}
