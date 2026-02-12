package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type JWTMaker struct {
	secretKey   string
	redisClient *redis.Client
}

func NewJWTMaker(secretKey string, redisClient *redis.Client) *JWTMaker {
	return &JWTMaker{secretKey: secretKey, redisClient: redisClient}
}

func (maker *JWTMaker) CreateToken(userID string, email string, role string, duration time.Duration, isRefreshToken bool) (string, *models.UserClaims, error) {
	claims, err := models.NewUserClaims(userID, email, role, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}

	if isRefreshToken {
		err = maker.redisClient.Set(context.Background(), tokenStr, userID, duration).Err()
	}

	return tokenStr, claims, err
}

func (maker *JWTMaker) VerifyToken(token string) (*models.UserClaims, error) {
	payload, err := jwt.ParseWithClaims(token, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := payload.Claims.(*models.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
