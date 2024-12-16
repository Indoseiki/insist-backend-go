package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var ACCESS_KEY = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
var REFRESH_KEY = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

func GenerateAccessToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(ACCESS_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(1 * 24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(REFRESH_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyAccessToken(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return ACCESS_KEY, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return 0, errors.New("token expired")
		}
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}

func VerifyRefreshToken(refreshToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return REFRESH_KEY, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return 0, errors.New("token expired")
		}
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}
