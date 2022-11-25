package structs

// PdataEntry represents single Pdata entry that returned by
// db.PdataList
type PdataEntry struct {
	ID   int64
	Name string
}
