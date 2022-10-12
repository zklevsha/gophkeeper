package structs

// MemStorage stores runtime state of client
type MemStorage struct {
	Token     string
	MasterKey string
}

func (m *MemStorage) SetToken(token string) {
	m.Token = token
}
