package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/errs"
)


const dsnDefault = "postgres://gophkeeper:gophkeeper@localhost:5532/gophkeeper_test?sslmode=disable"

var ctx = context.Background()



func setUp() Connector {
	// added so github action will be able to connect to test database
	var dsn = os.Getenv("GK_DB_TEST_DSN")
	if (dsn == "" ){
		dsn = dsnDefault
	}
	c := Connector{Ctx: ctx, DSN: dsn}
	err := c.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}

	// running migrations
	err = RunMigrations(dsn, "up")
	if err != nil {
		log.Fatalf("cannot run up migrations: %s", err.Error())
	}
	return c
}

func tearDown(c Connector) {
	err := RunMigrations(c.DSN, "down")
	if err != nil {
		log.Fatalf("cannot run down migrations: %s", err.Error())
	}
	c.Close()
}

func pdataConvert(ptype string, pname string, input interface{}) (Pdata, error) {
	masterKey := client.MasterKey{Key: getRandomSrt(32)}
	masterKey.SetHash()

	pdataBytes, err := json.Marshal(input)
	if err != nil {
		return Pdata{}, fmt.Errorf("cant marshall upass: %s", err.Error())
	}
	upassEnc, err := enc.EncryptAES(pdataBytes, []byte(masterKey.Key))
	if err != nil {
		return Pdata{}, fmt.Errorf("cant encrypt upass: %s", err.Error())
	}
	pdata := Pdata{
		Name:        pname,
		Type:        ptype,
		KeyHash:     base64.StdEncoding.EncodeToString(masterKey.KeyHash[:]),
		PrivateData: base64.StdEncoding.EncodeToString(upassEnc)}
	return pdata, nil
}

func TestRegister(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// run tests
	t.Run("Register", func(t *testing.T) {
		user := client.User{Email: "vasa@test.ru", Password: "test"}
		id, err := c.Register(user)
		if err != nil {
			t.Fatalf("Register() have returned an error: %s", err.Error())
		}

		conn, err := c.Pool.Acquire(c.Ctx)
		if err != nil {
			t.Fatalf("failed to acquire connection: %s", err.Error())
		}
		defer conn.Release()

		var email string
		sql := `SELECT email FROM users WHERE id=$1`
		err = conn.QueryRow(c.Ctx, sql, id).Scan(&email)
		if err != nil {
			t.Fatalf("failed to query users table: %s", err.Error())
		}
		if email != user.Email {
			t.Errorf("user data mismatch: have: %v, want: %v", email, user.Email)
		}
	})

}

func TestGetUser(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// Get user that exists
	want := client.User{Email: "vasya@test.ru", Password: "secret"}
	id, err := c.Register(want)
	if err != nil {
		t.Fatalf("cant register new user: %s", err.Error())
	}
	want.ID = id
	have, err := c.GetUser(want.Email)
	if err != nil {
		t.Fatalf("cant get user: %s", err.Error())
	}
	if want != have {
		t.Errorf("user mismatch: have: %v want: %v", have, want)
	}

	// Get user that does`t exists
	_, err = c.GetUser("john")
	if !errors.Is(err, errs.ErrUserAuth) {
		t.Errorf("err != structs.ErrUserAuth: %v", err)
	}
}

func TestAddPrivate(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)
	userID, err := c.Register(client.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}
	masterKey := client.MasterKey{Key: getRandomSrt(32)}
	masterKey.SetHash()

	// setup for Upass
	upass := client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err := pdataConvert("upass", "upass_test", upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}

	// TestCases
	tt := []struct {
		name  string
		pdata Pdata
	}{
		{name: "Test UPass",
			pdata: pdata,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := c.PrivateAdd(userID, tc.pdata)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}

func TestGetPrivate(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)
	userID, err := c.Register(client.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}
	masterKey := client.MasterKey{Key: getRandomSrt(32)}
	masterKey.SetHash()

	// setup for Upass
	upass := client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err := pdataConvert("upass", "upass_test", upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	pdataID, err := c.PrivateAdd(userID, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}
	pdata.ID = pdataID

	// TestCases
	tt := []struct {
		name string
		want Pdata
	}{
		{
			name: "Test UPass",
			want: pdata,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have, err := c.PrivateGet(userID, tc.want.ID)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if have != tc.want {
				t.Errorf("pdata mismatch: want %v, have %v", tc.want, have)
			}
		})
	}
}

func TestUpdatePrivate(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// adding test user
	userID, err := c.Register(client.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}

	// setup for Upass
	//before
	upassBefore := client.UPass{Name: "before", Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdataBefore, err := pdataConvert("upass", "upass_test", upassBefore)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	privateID, err := c.PrivateAdd(userID, pdataBefore)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}
	pdataBefore.ID = privateID
	// after
	upassAfter := client.UPass{Name: "after", Username: "userNew",
		Password: "passwordNew", Tags: map[string]string{"test": "testNew"}}
	pdataAfter, err := pdataConvert("upass", "upass_test", upassAfter)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	pdataAfter.ID = privateID

	tt := []struct {
		name string
		want Pdata
	}{
		{
			name: "Update UPass",
			want: pdataAfter,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := c.PrivateUpdate(userID, tc.want)
			if err != nil {
				t.Fatalf(err.Error())
			}
			have, err := c.PrivateGet(userID, tc.want.ID)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if have != tc.want {
				t.Errorf("pdata mismatch: want %v, have %v", tc.want, have)
			}
		})
	}
}

func TestPrivateList(t *testing.T) {

	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// adding test user
	userID, err := c.Register(client.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}

	// adding first Pdata
	pnameFirst := "upass1"
	upass := client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err := pdataConvert("upass", pnameFirst, upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	pdataFirstID, err := c.PrivateAdd(userID, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}

	// adding second pdata
	pnameSecond := "upass2"
	upass = client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err = pdataConvert("upass", pnameSecond, upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	pdataSecondID, err := c.PrivateAdd(userID, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}

	// TestCases
	tt := []struct {
		name  string
		ptype string
		want  []PdataEntry
	}{
		{
			name:  "Test Upass",
			ptype: "upass",
			want: []PdataEntry{
				{Name: pnameFirst, ID: pdataFirstID},
				{Name: pnameSecond, ID: pdataSecondID},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have, err := c.PrivateList(userID, tc.ptype)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if len(have) != len(tc.want) {
				t.Fatalf("result array lengh is wrong have: %d, want: %d",
					len(have), len(tc.want))
			}
			for idx, entryHave := range have {
				entryWant := tc.want[idx]
				if entryWant != entryHave {
					t.Errorf("entry mismatch have: %v, want: %v", entryHave, entryWant)
				}
			}
		})
	}

}

func TestPdataDelete(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// adding test user
	userID, err := c.Register(client.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}

	// adding another user
	userIDSecond, err := c.Register(client.User{Email: "petya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}

	// adding test Pdata
	pnameFirst := "upass_test"
	upass := client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err := pdataConvert("upass", pnameFirst, upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	firstID, err := c.PrivateAdd(userID, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}

	// adding Pdata of second user
	pnameSecond := "upass_second"
	upass = client.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	pdata, err = pdataConvert("upass", pnameSecond, upass)
	if err != nil {
		t.Fatalf("cant convert upass to pdata: %s", err.Error())
	}
	secondID, err := c.PrivateAdd(userIDSecond, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}

	// Case1: deleting pdata
	err = c.PrivateDelete(userID, firstID)
	if err != nil {
		t.Errorf("PrivateDelete returned an error: %s", err.Error())
	}

	// Case2 : getting notexistent pdata
	err = c.PrivateDelete(userID, 99)
	if !errors.Is(err, errs.ErrPdataNotFound) {
		t.Errorf("err != structs.ErrUserAuth: %v", err)
	}

	// Case3: make sure that user one not see user two Private data
	err = c.PrivateDelete(userID, secondID)
	if err == nil {
		t.Error("User one can delete pdata of second user")
	}
	if !errors.Is(err, errs.ErrPdataNotFound) {
		t.Errorf("err != structs.ErrUserAuth: %v", err)
	}
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
