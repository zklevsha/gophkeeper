package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

const dsn = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper_test"

var ctx = context.Background()

func setUp() Connector {
	c := Connector{Ctx: ctx, DSN: dsn}
	err := c.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}
	err = c.CreateTables()
	if err != nil {
		log.Fatalf("Cant create tables: %s", err.Error())
	}
	return c
}

func tearDown(c Connector) {
	err := c.DropTables()
	if err != nil {
		log.Fatalf("DropTables have returned ad error: %s", err.Error())
	}
}

func TestRegister(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)

	// run tests
	t.Run("Register", func(t *testing.T) {
		user := structs.User{Email: "vasa@test.ru", Password: "test"}
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
	want := structs.User{Email: "vasya@test.ru", Password: "secret"}
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
	if !errors.Is(err, structs.ErrUserAuth) {
		t.Errorf("err != structs.ErrUserAuth: %v", err)
	}
}

func TestAddPrivate(t *testing.T) {
	// setup
	c := setUp()
	defer c.Close()
	defer tearDown(c)
	userId, err := c.Register(structs.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}

	// setup for Upass
	masterKey := structs.MasterKey{Key: helpers.GetRandomSrt(32)}
	masterKey.SetHash()
	upass := structs.UPass{Username: "user",
		Password: "password", Tags: map[string]string{"test": "test"}}
	upassBytes, err := json.Marshal(upass)
	if err != nil {
		t.Fatalf("cant marshall upass: %s", err.Error())
	}
	upassEnc, err := enc.EncryptAES(upassBytes, []byte(masterKey.Key))
	if err != nil {
		t.Fatalf("cant encrypt upass: %s", err.Error())
	}

	// TestCases
	tt := []struct {
		name  string
		ptype string
		pdata []byte
	}{
		{name: "Test UPass", ptype: "upass", pdata: upassEnc},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := c.PrivateAdd(tc.name, userId, tc.ptype, masterKey.KeyHash, tc.pdata)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
