package helpers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRandomSrt(strLen int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, strLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// ToPdata converts Upass/Cards... to Pdata
func ToPdata(ptype string, input interface{}, key structs.MasterKey) (*pb.Pdata, error) {
	var name string
	var encoded []byte
	var err error
	switch ptype {
	case "upass":
		up := input.(structs.UPass)
		name = up.Name
		encoded, err = json.Marshal(up)
		if err != nil {
			return nil, fmt.Errorf("cant encode upass: %s", err.Error())
		}
	case "card":
		card := input.(structs.Card)
		name = card.Name
		encoded, err = json.Marshal(card)
		if err != nil {
			return nil, fmt.Errorf("cant encode card: %s", err.Error())
		}
	case "pstring":
		pstring := input.(structs.Pstring)
		name = pstring.Name
		encoded, err = json.Marshal(pstring)
		if err != nil {
			return nil, fmt.Errorf("cant encode pstring: %s", err.Error())
		}
	case "pfile":
		pfile := input.(structs.Pfile)
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

func FromPdata(pdata *pb.Pdata, key structs.MasterKey) (interface{}, error) {
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
		var up structs.UPass
		err = json.Unmarshal(decrypted, &up)
		if err != nil {
			return nil, fmt.Errorf("error cant decode upass JSON to struct: %s", err.Error())
		}
		return up, nil
	case "card":
		var card structs.Card
		err = json.Unmarshal(decrypted, &card)
		if err != nil {
			return nil, fmt.Errorf("error cant decode card JSON to struct: %s", err.Error())
		}
		return card, nil
	case "pstring":
		var pstring structs.Pstring
		err = json.Unmarshal(decrypted, &pstring)
		if err != nil {
			return nil, fmt.Errorf("error cant decode pstring json to struct: %s", err.Error())
		}
		return pstring, nil
	case "pfile":
		var pfile structs.Pfile
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
