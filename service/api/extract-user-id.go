package api

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func extractUserId(tokenCode string) (string, error) {
	//Decode the token
	token, err := jwt.Parse(tokenCode, func(token *jwt.Token) (interface{}, error) {
		if _, flag := token.Method.(*jwt.SigningMethodHMAC); !flag {
			return nil, errors.New("Signing method error")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("Invalid token")
	}

	//Extract claims from the token
	claims, flag := token.Claims.(jwt.MapClaims) //map[string]interface{}
	if !flag {
		return "", errors.New("Token extraction error")
	}

	userID, flag := claims["user_id"].(string)
	if !flag {
		return "", errors.New("Cast claims error")
	}

	return userID, nil
}
