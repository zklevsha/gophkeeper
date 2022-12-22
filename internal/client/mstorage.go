package client

// MemStorage stores runtime state of client
type MemStorage struct {
	Token        string
	MasterKey    MasterKey
	MasterKeyDir string
	PfilesDir    string
}

// SetToken sets/updates token
func (m *MemStorage) SetToken(token string) {
	m.Token = token
}

// SetMasterKey sets/updates MasterKey
func (m *MemStorage) SetMasterKey(key string, keyPath string) {
	m.MasterKey.Key = key
	m.MasterKey.KeyPath = keyPath
	m.MasterKey.SetHash()
}
