package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/docker/distribution/registry/auth"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"time"
)

type TokenService struct {
	signingKey []byte
}

func NewTokenService(configuration *config.Configuration) (*TokenService, error) {
	return &TokenService{signingKey: []byte(configuration.JwtSigningKey)}, nil
}

func (t *TokenService) CreateToken(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("user is null or empty")
	}
	claims := models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24 * 30 * 360)), // TODO: will change on OAuth later
			IssuedAt:  jwt.At(time.Now()),
		},
		Id:       user.Id,
		Username: user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.signingKey)
}

func (t *TokenService) ParseToken(accessToken string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, new(models.Claims), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, auth.ErrInvalidCredential
}
