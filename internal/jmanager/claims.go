package jmanager

// Claims represent JWT token claims
type Claims struct {
	UserID int64
	Iat    int64
	Exp    int64
}
