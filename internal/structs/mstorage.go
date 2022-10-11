package structs

import "github.com/golang-jwt/jwt"

// MemStorage stores runtime state of client
type MemStorage struct {
	Token *jwt.Token
}
