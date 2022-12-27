package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/errs"
	"github.com/zklevsha/gophkeeper/internal/pb"
)

// reqCheck checks that user is logged in and has master key loaded
func reqCheck(mstorage *MemStorage) error {
	if mstorage.MasterKey.Key == "" {
		return errs.ErrMasterKeyIsMissing

	}
	if mstorage.Token == "" {
		return errs.ErrLoginRequired
	}
	return nil
}

// listPnames getting list of available pnames
// list returned as <pname>: <pid> map
func listPnames(ctx context.Context, gclient *Gclient, ptype string) (map[string]int64, error) {
	listResponse, err := gclient.Pdata.ListPdata(ctx, &pb.ListPdataRequest{Ptype: ptype})
	if err != nil {
		return map[string]int64{}, err
	}

	entries := make(map[string]int64)
	for _, e := range listResponse.PdataEtnry {
		entries[e.Name] = e.ID
	}

	return entries, nil
}

// toPdata converts Upass/Card/Pfile/Pstring to Pdata
func toPdata(ptype string, input interface{}, key MasterKey) (*pb.Pdata, error) {
	var name string
	var encoded []byte
	var err error
	switch ptype {
	case "upass":
		up := input.(UPass)
		name = up.Name
		encoded, err = json.Marshal(up)
		if err != nil {
			return nil, fmt.Errorf("cant encode upass: %s", err.Error())
		}
	case "card":
		card := input.(Card)
		name = card.Name
		encoded, err = json.Marshal(card)
		if err != nil {
			return nil, fmt.Errorf("cant encode card: %s", err.Error())
		}
	case "pstring":
		pstring := input.(Pstring)
		name = pstring.Name
		encoded, err = json.Marshal(pstring)
		if err != nil {
			return nil, fmt.Errorf("cant encode pstring: %s", err.Error())
		}
	case "pfile":
		pfile := input.(Pfile)
		name = pfile.Name
		buf := bytes.Buffer{}
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(pfile)
		if err != nil {
			return nil, fmt.Errorf("cant encode pfile: %s", err.Error())
		}
		encoded = buf.Bytes()
	default:
		return nil, fmt.Errorf("%s is not supported", ptype)
	}
	encrypted, err := enc.EncryptAES(encoded, []byte(key.Key))
	if err != nil {
		return nil, fmt.Errorf("cant encrypt data: %s", err.Error())
	}
	pdata := pb.Pdata{Pname: name, Ptype: ptype, Pdata: encrypted, KeyHash: key.KeyHash[:]}
	return &pdata, nil
}

// fromPdata converts Pdata to Upass/Card/Pfile/Pstring
func fromPdata(pdata *pb.Pdata, key MasterKey) (interface{}, error) {
	// check if we using correct master key
	if string(pdata.KeyHash) != string(key.KeyHash) {
		e := fmt.Errorf(`key hash mismatch
					hash of the key that used to encrypt pdata: %v\
					masterkey hash: %v`, pdata.KeyHash, key.KeyHash)
		return nil, e
	}

	// decrypt
	decrypted, err := enc.DecryptAES(pdata.Pdata, []byte(key.Key))
	if err != nil {
		return nil, fmt.Errorf("error cant decrypt data: %s", err.Error())
	}

	// decode
	switch pdata.Ptype {
	case "upass":
		var up UPass
		err = json.Unmarshal(decrypted, &up)
		if err != nil {
			return nil, fmt.Errorf("error cant decode upass JSON to struct: %s", err.Error())
		}
		return up, nil
	case "card":
		var card Card
		err = json.Unmarshal(decrypted, &card)
		if err != nil {
			return nil, fmt.Errorf("error cant decode card JSON to struct: %s", err.Error())
		}
		return card, nil
	case "pstring":
		var pstring Pstring
		err = json.Unmarshal(decrypted, &pstring)
		if err != nil {
			return nil, fmt.Errorf("error cant decode pstring json to struct: %s", err.Error())
		}
		return pstring, nil
	case "pfile":
		var pfile Pfile
		dec := gob.NewDecoder(bytes.NewReader(decrypted))
		err := dec.Decode(&pfile)
		if err != nil {
			return nil, fmt.Errorf("error cant decode to struct: %s", err.Error())
		}
		return pfile, nil
	default:
		return nil, fmt.Errorf("ptype: %s is not supported", pdata.Ptype)
	}
}



// listDir reads directory (non recurcevly) and returns full paths of files
func listDir(dirPath string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(dirPath)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() {
			fullPath := filepath.Join(dirPath, file.Name())
			files = append(files, fullPath)

		}

	}
	return files, nil
}


func getRandomSrt(strLen int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, strLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

