package structs

type UPass struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Tags     map[string]string `json:"tags,omitempty"`
}
