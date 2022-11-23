package structs

// Pdata represendts raw and encrypted user`s private data
// received from database
type Pdata struct {
	ID          int64
	Name        string
	Type        string
	KeyHash     string
	PrivateData string
}
