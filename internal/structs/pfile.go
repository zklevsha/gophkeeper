package structs

// Pfile represents user`s private file
type Pfile struct {
	Name string
	Data []byte
	Tags map[string]string
}
