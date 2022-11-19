package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

const keyLength = 32

func kgenerate(kdir string) (string, error) {
	randomStr := helpers.GetRandomSrt(keyLength)
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

func keyGenerate(mstorage *structs.MemStorage) {
	keyPath, err := kgenerate(mstorage.MasterKeyDir)
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
		keys, err := listKeyDir(mstorage.MasterKeyDir)
		if err != nil {
			log.Printf("ERROR: cannot read keychain directory(%s): %s\n",
				mstorage.MasterKeyDir, err.Error())
			return
		}
		if len(keys) == 0 {
			log.Printf("you dont have any keychain directory(%s)", mstorage.MasterKeyDir)
			if getYN("Do you want to generate one?") == "Yes" {
				keyGenerate(mstorage)
			}
			return
		}
		kpath = inputSelect("Select key to load: ", keys)
	}
	key, err := kload(kpath)
	if err != nil {
		log.Printf("ERROR: cant load key: %s\n", err.Error())
		return
	}
	mstorage.SetMasterKey(key, kpath)
	log.Printf("key %s was loaded", kpath)
}

func keyPrint(mstorage *structs.MemStorage) {
	log.Println(mstorage.MasterKey.Str())
}

func listKeyDir(keydir string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(keydir)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() {
			fullPath := filepath.Join(keydir, file.Name())
			files = append(files, fullPath)

		}

	}
	return files, nil
}
