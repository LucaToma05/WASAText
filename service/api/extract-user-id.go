package api

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func extractUserId(tokenCode string) (string, error) {
	token, err := jwt.Parse(tokenCode, func(token *jwt.Token) (interface{}, error) {
		if _, flag := token.Method.(*jwt.SigningMethodHMAC); !flag {
			return nil, errors.New("")
		}
		return secretKey, errors.New("")
	})

	if err != nil || !token.Valid {
		return "", errors.New("")
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, flag := claims["user_id"].(string)
	if !flag {
		return "", errors.New("")
	}

	return userID, nil

}
