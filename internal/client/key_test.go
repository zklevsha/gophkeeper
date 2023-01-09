package client

import (
	"os"
	"testing"
)

func TestKey(t *testing.T) {

	// GenerateKey
	kpath, err := kgenerate("/tmp")
	if err != nil {
		t.Fatalf("GenerateKey have returned an error: %s", err.Error())
	}

	// LoadKey
	key, err := kload(kpath)
	if err != nil {
		t.Fatalf("LoadKey have returned an errorL %s", err.Error())
	}
	have := len(key)
	if keyLength != have {
		t.Errorf("key length mismatch: have: %d, want: %d", have, keyLength)
	}

	// DeleteKey
	err = os.Remove(kpath)
	if err != nil {
		t.Fatalf("failed to remove %s: %s", kpath, err.Error())
	}
}
