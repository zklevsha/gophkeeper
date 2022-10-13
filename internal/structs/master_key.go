package structs

import (
	"crypto/sha256"
	"fmt"
)

// MasterKey represents master key
type MasterKey struct {
	Key     string
	KeyPath string
	KeyHash [32]byte
}

func (m *MasterKey) SetHash() {
	hash := sha256.Sum256([]byte(m.Key))
	m.KeyHash = hash
}

func (m *MasterKey) Str() string {
	return fmt.Sprintf("<MasterKey key:'%s', keyPath:'%s' keyHash: '%x'>",
		m.Key, m.KeyPath, m.KeyHash)
}
