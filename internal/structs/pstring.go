package structs

// PricateString represents user`s private string
type Pstring struct {
	Name   string            `json:"name"`
	String string            `json:"string"`
	Tags   map[string]string `json:"tags,omitempty"`
}
