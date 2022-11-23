package structs

type UPass struct {
	Name     string            `json:"name"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Tags     map[string]string `json:"tags,omitempty"`
}
