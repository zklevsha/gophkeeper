package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Generate generates new JWT token
func Generate(userid int, key string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userid,
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * 5).Unix(),
	})
	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}

// GetUserID exstracts user id from JWT token
func GetUserID(tokenString string, key string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["id"].(float64)), nil
	} else {
		return -1, err
	}
}
