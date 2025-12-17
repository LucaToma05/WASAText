package api

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func extractUserId(tokenCode string) (int, error) {
	//Decode the token
	token, err := jwt.Parse(tokenCode, func(token *jwt.Token) (interface{}, error) {
		if _, flag := token.Method.(*jwt.SigningMethodHMAC); !flag {
			return nil, errors.New("signing method error")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	//Extract claims from the token
	claims, flag := token.Claims.(jwt.MapClaims) //map[string]interface{}
	if !flag {
		return 0, errors.New("token extraction error")
	}

	userID, flag := claims["user_id"].(int)
	if !flag {
		return 0, errors.New("cast claims error")
	}

	return userID, nil
}
