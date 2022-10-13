package structs

// MemStorage stores runtime state of client
type MemStorage struct {
	Token     string
	MasterKey MasterKey
}

func (m *MemStorage) SetToken(token string) {
	m.Token = token
}

func (m *MemStorage) SetMasterKey(key string, keyPath string) {
	m.MasterKey.Key = key
	m.MasterKey.KeyPath = keyPath
	m.MasterKey.SetHash()
}
