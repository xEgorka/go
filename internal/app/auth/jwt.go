package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims contains claims, inluding own one: userID.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// Token settings.
const (
	tokenExp  = 30 * time.Minute // expiration period
	secretKey = "supersecretkey" // secret key
)

// BuildJWTString creates a token and returns it as a string.
func BuildJWTString(usr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: usr, // own claim
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserID extracts UserID from token string.
func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v",
					t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil || !token.Valid {
		return claims.UserID, err
	}
	return claims.UserID, nil
}
