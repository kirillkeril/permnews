package jwt_service

import (
	"authorization_service/internal/app/models/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtService struct {
	secretKey  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type JwtPair struct {
	Access  string
	Refresh string
}

func (s JwtService) CreateTokenPair(user entity.User) (*JwtPair, error) {
	key := []byte(s.secretKey)

	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
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

func (s JwtService) VerifyToken(token string) (*jwt.RegisteredClaims, bool) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, false
	}
	return &claims, true
}
