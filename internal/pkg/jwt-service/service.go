package jwtService

import (
	"authorization_service/internal/models/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey string = "<KEY>"
var accessTTL time.Duration = time.Hour
var refreshTTL time.Duration = time.Hour * 24

type JwtPair struct {
	Access  string
	Refresh string
}

func ConfigureService(secret string, accessTtl time.Duration, refreshTtl time.Duration) {
	secretKey = secret
	accessTTL = accessTtl
	refreshTTL = refreshTtl
}

func CreateTokenPair(user entity.User) (*JwtPair, error) {
	key := []byte(secretKey)

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString(key)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		ID:        user.Id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &JwtPair{Access: accessToken, Refresh: refreshToken}, nil
}

func VerifyToken(token string) (*jwt.RegisteredClaims, bool) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, false
	}
	return &claims, true
}
