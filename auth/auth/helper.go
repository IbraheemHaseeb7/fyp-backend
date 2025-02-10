package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserToken struct {
	Name	 			string
	Email    			string
	RegistrationNumber 	string
}

func NewUserToken(name, email, registrationNumber string) *UserToken {
	return &UserToken{
		Name: name,
		Email:    email,
		RegistrationNumber: registrationNumber,
	}
}

func CreateToken(user UserToken, minutes time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name": user.Name,
			"email":    user.Email,
			"registrationNumber":    user.RegistrationNumber,
			"exp":      time.Now().Add(time.Minute * minutes).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RefreshToken(tokenString string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("Claims could not be extracted")
	}

	name := claims["name"].(string)
	email := claims["email"].(string)
	registrationNumber := claims["registrationNumber"].(string)

	newToken, err := CreateToken(*NewUserToken(name, email, registrationNumber), 60)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateToken(*NewUserToken(name, email, registrationNumber), 1440)
	if err != nil {
		return "", "", err
	}

	return newToken, refreshToken, nil
}

func VerifyToken(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid Token")
	}

	return nil
}

func GetClaimsFromToken(tokenString string) (map[string]any, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return make(map[string]any), err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return make(map[string]any), fmt.Errorf("Could not extract claims from JWT token")
	}

	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
