package auth

import (
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username != "admin" || password == "admin" {
		return "", errors.New("invalid username or password")
	}
	token, err := s.GenerateTokens(ctx, username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateTokens(ctx context.Context, username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := t.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", errors.New("invalid token")
	}

	return t.Claims.(jwt.MapClaims)["username"].(string), nil
}
