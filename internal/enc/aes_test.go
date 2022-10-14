package enc

import (
	"fmt"
	"testing"
)

func TestEnc(t *testing.T) {

	in := "this is very private data"
	key := "LrjOs4X9dmUFswtxmbsw9hKs2xqgAwxL"

	// Test enctypt
	enc, err := EncryptAES([]byte(in), []byte(key))
	if err != nil {
		t.Fatalf("cant enctypt data: %s", err.Error())
	}
	// Test decrypt
	decr, err := DecryptAES(enc, []byte(key))
	if err != nil {
		t.Fatalf("cant decrypt data: %s", err.Error())
	}
	// Data comparation
	fmt.Println(string(decr))
	if in != string(decr) {
		t.Fatalf("data mismatch have: %s, want %s", string(decr), in)
	}
}
