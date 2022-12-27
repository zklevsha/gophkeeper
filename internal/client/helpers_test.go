package client

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/errs"
	"github.com/zklevsha/gophkeeper/internal/pb"
)

func TestReqCheck(t *testing.T) {
	tt := []struct {
		name  string
		mstorage MemStorage
		errWant error
	}{
		{
			name: "Good mstorage",
			mstorage: MemStorage{Token: "ds';mvzx[k;mscvsdfsx", MasterKey: MasterKey{Key: "Test"}} ,
			errWant: nil,
		},

		{
			name: "No Master key",
			mstorage: MemStorage{Token: "sdsdsfsc"},
			errWant: errs.ErrMasterKeyIsMissing,
		},

		{
			name: "No Token",
			mstorage: MemStorage{MasterKey: MasterKey{Key: "Test"}},
			errWant: errs.ErrLoginRequired,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := reqCheck(&tc.mstorage)
			if !errors.Is(err, tc.errWant) {
				t.Errorf("error mismatch: have %s, want %s", err, tc.errWant)
			}
		})
	}
}

func TestToPdata(t *testing.T) {
	masterKey := MasterKey{Key: GetRandomSrt(32)}
	tt := []struct {
		name  string
		input interface{}
		inputType string
	}{
		{
			name: "Upass",
			inputType: "upass",
			input: UPass{
				Name: "Test upass",
				Username: "vasya", Password: "pupkin",
				Tags: map[string]string{"test":"test"},
			},
		},
		{
			name: "Card",
			inputType: "card",
			input: Card{
				Name: "Test Card",
				Number: "1111 3343 4344 0000",
				Holder: "Ivanov Ivan",
				Expire: "11/20",
				CVC: "123",
				Tags: map[string]string{"test":"test"},
			},
		},
		{
			name: "Pstring",
			inputType: "pstring",
			input: Pstring{
				Name: "Test Pstring",
				String: "Test pstring",
				Tags: map[string]string{"test":"test"},
			},
		},
		{
			name: "Pfile",
			inputType: "pfile",
			input: Pfile{
				Name: "Test pfile",
				Data: []byte("Test file"),
				Tags: map[string]string{"test":"test"},
			},
		},

	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ToPdata(tc.inputType, tc.input, masterKey)
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestFromPdata(t *testing.T) {
	masterKey := MasterKey{Key: GetRandomSrt(32)}

	// test Upass
	upassTest  := UPass{
		Name: "Test upass",
		Username: "vasya", Password: "pupkin",
		Tags: map[string]string{"test":"test"},
	}
	upassTestPdata, err := ToPdata("upass", upassTest, masterKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	upassJSON, err := json.Marshal(upassTest)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// test Card
	cardTest := Card{
		Name: "Test Card",
		Number: "1111 3343 4344 0000",
		Holder: "Ivanov Ivan",
		Expire: "11/20",
		CVC: "123",
		Tags: map[string]string{"test":"test"},
	}
	cardTestPdata, err := ToPdata("card", cardTest, masterKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	cardJSON, err := json.Marshal(cardTest)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// pstring
	pstringTest := Pstring{
		Name: "Test pstring",
		String: "secret string",
		Tags: map[string]string{"test":"test"},
	}
	pstringTestPdata, err := ToPdata("pstring", pstringTest, masterKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	pstringJSON, err := json.Marshal(pstringTest)
	if err != nil {
		t.Fatalf(err.Error())
	}


	tt := []struct {
		name  string
		input *pb.Pdata
		want string
	}{
		{
			name: "Upass",
			input: upassTestPdata,
			want: string(upassJSON),

		},
		{
			name: "Card",
			input: cardTestPdata,
			want: string(cardJSON),
		},
		{
			name: "Pstring",
			input: pstringTestPdata,
			want: string(pstringJSON),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have, err := FromPdata(tc.input, masterKey)
			if err != nil {
				t.Error(err.Error())
			}
			haveJSON, err := json.Marshal(have)
			if err != nil {
				t.Errorf("cant convert struct to JSON: %s", err.Error())
			}
			if string(haveJSON) != tc.want {
				t.Errorf("struct mismatch: have %s, want %s", string(haveJSON), tc.want)
			}
		})
	}
}

func TestGetRandomStr(t *testing.T) {
	wantSize := 32
	testStr := GetRandomSrt(wantSize)
	haveSize := len(testStr)
	if wantSize != haveSize {
		t.Errorf("string lenght mismatch: want: %d, have: %d",
			wantSize, haveSize)
	}
}

