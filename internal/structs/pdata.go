package structs

// `SELECT a.name, b.name, a.khash_base64, a.data_base64
// 			FROM private_data AS a
// 			INNER JOIN private_types AS b
// 			ON a.type_id=b.id;`

// Pdata represendts raw and encrypted user`s private data
// received from database
type Pdata struct {
	Name        string
	Type        string
	KeyHash     string
	PrivateData string
}
