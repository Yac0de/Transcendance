package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	ID uint
	// ADD CREATEDAT FIELD ?
}

func CreateToken(id uint) (string, error) {
	userToken := UserToken{
		ID: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": userToken.ID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNATURE")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecryptToken(AccessToken string) (*UserToken, error) {
	key := []byte(os.Getenv("JWT_SIGNATURE"))
	token, err := jwt.Parse(AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signature method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["ID"].(float64)
		if !ok {
			return nil, fmt.Errorf("ID manquant ou de type incorrect")
		}

		return &UserToken{
			ID: uint(id),
		}, nil
	}
	return nil, fmt.Errorf("Invalid token")
}
