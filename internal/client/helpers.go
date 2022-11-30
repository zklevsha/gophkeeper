package client

import (
	"errors"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

// reqCheck checks that user is logged in and has master key loaded
func reqCheck(mstorage *structs.MemStorage) error {
	if mstorage.MasterKey.Key == "" {
		return errors.New("master-key does not exists add it via key-generate/key-load commands")
	}
	if mstorage.Token == "" {
		return errors.New("login required (login)")
	}
	return nil
}
