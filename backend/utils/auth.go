package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateToken(userID primitive.ObjectID, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (primitive.ObjectID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return primitive.NilObjectID, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid claims")
	}

	userIDText, ok := claims["userId"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("missing user id")
	}

	return primitive.ObjectIDFromHex(userIDText)
}
