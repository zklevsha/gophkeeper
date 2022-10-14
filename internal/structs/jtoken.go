package structs

// Claims represent JWT token claims
type Claims struct {
	UserID int64
	Iat    int64
	Exp    int64
}

// Jtoken represents JWT token
type Jtoken struct {
	Token  string
	Claims Claims
}
