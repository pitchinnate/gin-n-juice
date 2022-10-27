package models

import (
	"fmt"
	"gin-n-juice/config"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestUserModel(t *testing.T) {
	hashedPassword := ""
	t.Run("Test HashPassword", func(t *testing.T) {
		password := "testing"
		result, err := HashPassword(password)
		if err != nil {
			t.Error("Error hashing a password", err.Error())
		}
		if len(result) < 60 {
			t.Error("Password hash seems to be to small")
		}
		hashedPassword = result
	})
	t.Run("Test HashPassword", func(t *testing.T) {
		result := CheckPasswordHash("testing", hashedPassword)
		if !result {
			t.Error("Check has failed")
		}
	})
	t.Run("Test GenerateJwt", func(t *testing.T) {
		user := User{
			Email:         "test@test.com",
			Password:      hashedPassword,
			Admin:         false,
			EmailVerified: nil,
		}
		user.ID = 1

		tokenString := GenerateJwt(user)

		if len(tokenString) < 10 {
			t.Error("Token seems to be to small")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			hmacSampleSecret := []byte(config.ENCRYPT_KEY)
			return hmacSampleSecret, nil
		})

		if err != nil {
			t.Error("Error parsing token", err.Error())
		} else {
			_, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				t.Error("Could not get claims from token")
			}
		}
	})
}
