package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	HashedPassword string `json:"password"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(register RegisterUser) (User, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(register.Password), 10)
	
	if error != nil {
		return User{}, error
	}

	return User{
		Username: register.Username,
		HashedPassword: string(hashedPassword),
	}, nil 
}

func ValidatePassword(hashedPassword, plaintextPassword string) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))

	return error == nil
}

func CreateToken (user User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user": user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := "secret-key"

	tokenString, error := token.SignedString([]byte(secret))

	if error != nil {
		return ""
	}

	return tokenString
}