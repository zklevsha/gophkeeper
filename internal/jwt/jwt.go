package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// Generate generates new JWT token
func Generate(userid int64, key string) (structs.Jtoken, error) {
	now := time.Now()
	claims := structs.Claims{UserId: userid, Iat: now.Unix(),
		Exp: now.Add(time.Minute * 60).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  claims.UserId,
		"iat": claims.Iat,
		"exp": claims.Exp,
	})
	tokenString, err := token.SignedString([]byte(key))
	return structs.Jtoken{Claims: claims, Token: tokenString}, err
}

// Validate checks JWT token and converts to structs.Jtoken
func Validate(tokenString string, key string) (structs.Jtoken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims := structs.Claims{
			UserId: int64(claimsMap["id"].(float64)),
			Iat:    int64(claimsMap["iat"].(float64)),
			Exp:    int64(claimsMap["exp"].(float64)),
		}
		jtoken := structs.Jtoken{Token: tokenString, Claims: claims}
		return jtoken, nil
	} else {
		return structs.Jtoken{}, err
	}
}
