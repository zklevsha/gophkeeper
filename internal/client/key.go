package client

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const keyLength = 32

func kgenerate(kdir string) (string, error) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	kname := fmt.Sprintf("gkeeper-%d", time.Now().UnixNano())
	kpath := path.Join(kdir, kname)
	err := os.WriteFile(kpath, b, 0600)
	if err != nil {
		return "", fmt.Errorf("cant create key file %s: %s", kname, err.Error())
	}
	return kpath, nil
}

func kload(kpath string) (string, error) {
	b, err := os.ReadFile(kpath)
	if err != nil {
		return "", fmt.Errorf("cant read %s: %s", kpath, err.Error())
	}
	return string(b), nil
}

func keyGenerate(mstorage *structs.MemStorage) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get HOME dir: %s", err.Error())
		return
	}
	keyPath, err := kgenerate(homeDir)
	if err != nil {
		log.Printf("cant create key file: %s", err.Error())
		return
	}
	log.Printf("key saved at %s", keyPath)
	if getYN("Do you want lo load key?") == "Yes" {
		keyLoad(keyPath, mstorage)
	}
}

func keyLoad(kpath string, mstorage *structs.MemStorage) {
	if kpath == "" {
		kpath = getInput("keyFile path: ", notEmpty, false)
	}
	key, err := kload(kpath)
	if err != nil {
		log.Printf("cant load key: %s", err.Error())
		return
	}
	mstorage.SetMasterKey(key, kpath)
	log.Printf("key %s was loaded", kpath)
}

func keyPrint(mstorage *structs.MemStorage) {
	log.Println(mstorage.MasterKey.Str())
}
