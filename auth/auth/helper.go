package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserToken struct {
	ID                 float64
	Name               string
	Email              string
	RegistrationNumber string
}

func NewUserToken(id float64, name, email, registrationNumber string) *UserToken {
	return &UserToken{
		ID:                 id,
		Name:               name,
		Email:              email,
		RegistrationNumber: registrationNumber,
	}
}

func CreateToken(user UserToken, minutes time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":                 user.ID,
			"name":               user.Name,
			"email":              user.Email,
			"registrationNumber": user.RegistrationNumber,
			"exp":                time.Now().Add(time.Minute * minutes).Unix(),
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

	name, ok := claims["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("Could not assert type")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("Could not assert type")
	}
	registrationNumber, ok := claims["registrationNumber"].(string)
	if !ok {
		return "", "", fmt.Errorf("Could not assert type")
	}
	id, ok := claims["id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("Could not assert type")
	}

	newToken, err := CreateToken(*NewUserToken(id, name, email, registrationNumber), 60)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateToken(*NewUserToken(id, name, email, registrationNumber), 1440)
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
