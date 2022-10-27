package models

import (
	"gin-n-juice/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string     `json:"email"`
	Password      string     `json:"-"`
	Admin         bool       `json:"admin"`
	EmailVerified *time.Time `json:"email_verified"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Password", password)
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwt(user User) string {
	mySigningKey := []byte(config.ENCRYPT_KEY)

	type MyCustomClaims struct {
		UserId uint   `json:"user_id"`
		Email  string `json:"email"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		user.ID,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + (3600 * 24 * 7),
			Issuer:    "gin-n-juice",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	return ss
}
