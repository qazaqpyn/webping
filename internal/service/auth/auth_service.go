package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secret = "sadfasdfkjoi1j23124"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username != "admin" || password != "admin" {
		return "", errors.New("invalid username or password")
	}
	token, err := s.GenerateTokens(ctx, username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateTokens(ctx context.Context, username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Subject:   username,
	})
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("token are not *TokenClaim types")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid object")
	}

	return subject, nil
}
