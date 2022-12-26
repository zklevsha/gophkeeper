package client

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

const keyLength = 32

func kgenerate(kdir string) (string, error) {
	randomStr := getRandomSrt(keyLength)
	kname := fmt.Sprintf("gkeeper-%d", time.Now().UnixNano())
	kpath := path.Join(kdir, kname)
	err := os.WriteFile(kpath, []byte(randomStr), 0600)
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

func keyGenerate(mstorage *MemStorage) error {
	keyPath, err := kgenerate(mstorage.MasterKeyDir)
	if err != nil {
		return fmt.Errorf("cant create key file: %s", err.Error())
	}
	log.Printf("key saved at %s", keyPath)
	if getYN("Do you want lo load key?") {
		err = keyLoad(keyPath, mstorage)
		return err
	}
	return nil
}

func keyLoad(kpath string, mstorage *MemStorage) error {
	if kpath == "" {
		keys, err := listDir(mstorage.MasterKeyDir)
		if err != nil {
			return fmt.Errorf("cannot read keychain directory(%s): %s",
				mstorage.MasterKeyDir, err.Error())

		}
		if len(keys) == 0 {
			log.Printf("you dont have any keychain directory(%s)", mstorage.MasterKeyDir)
			if getYN("Do you want to generate one?") {
				return keyGenerate(mstorage)
			}
			return nil
		}
		kpath = inputSelect("Select key to load", keys)
	}
	key, err := kload(kpath)
	if err != nil {
		return fmt.Errorf("cant load key: %s", err.Error())
	}
	mstorage.SetMasterKey(key, kpath)
	log.Printf("key %s was loaded", kpath)
	return nil
}

func keyPrint(mstorage *MemStorage) {
	log.Println(mstorage.MasterKey.Str())
}
