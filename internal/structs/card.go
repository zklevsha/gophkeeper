package structs

// Card represents users Credit card entry
type Card struct {
	Name   string            `json:"name"`
	Number string            `json:"number"`
	Holder string            `json:"holder"`
	Expire string            `json:"expire"`
	CVC    string            `json:"cvv"`
	Tags   map[string]string `json:"tags,omitempty"`
}
