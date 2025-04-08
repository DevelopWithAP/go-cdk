package types

import (
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