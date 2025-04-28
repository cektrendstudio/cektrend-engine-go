package helper

import (
	"errors"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	"github.com/dgrijalva/jwt-go"
)

var (
	accessSecret  = []byte(utstring.Env("ACCESS_SECRET_KEY", "Sp3Sk1llT3stAccess"))
	refreshSecret = []byte(utstring.Env("REFRESH_SECRET_KEY", "Sp3Sk1llT3stRefresh"))
)

func GenerateAccessToken(userID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(accessSecret)
	return signedToken, expiresAt, err
}

func GenerateRefreshToken(userID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(14 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(refreshSecret)
	return signedToken, expiresAt, err
}

func VerifyToken(tokenString, typeToken string) (jwt.MapClaims, error) {
	errResponse := errors.New("Token-Invalid")
	var secret []byte

	switch typeToken {
	case "access":
		secret = accessSecret
	case "refresh":
		secret = refreshSecret
	default:
		return nil, errResponse
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errResponse
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errResponse
	}

	return claims, nil
}
