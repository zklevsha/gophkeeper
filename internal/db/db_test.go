package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

const dsn = "postgres://gophkeeper:gophkeeper@localhost:5532/gophkeeper_test?sslmode=disable"
const migrationsFolder = "file://../../db/migrations"

var ctx = context.Background()

func runMigrations(direction string) error {
	migrate, err := migrate.New(migrationsFolder, dsn)
	if err != nil {
		return fmt.Errorf("cannot init migrate object: %s", err.Error())
	}
	defer migrate.Close()
	switch direction {
	case "up":
		return migrate.Up()
	case "down":
		return migrate.Down()
	default:
		return fmt.Errorf("bad direction parameter: %s (only up/down are supported)", direction)
	}
}

func setUp() Connector {
	c := Connector{Ctx: ctx, DSN: dsn}
	err := c.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}

	// running migrations
	err = runMigrations("up")
	if err != nil {
		log.Fatalf("cannot run up migrations: %s", err.Error())
	}
	return c
}

func tearDown(c Connector) {
	err := runMigrations("down")
	if err != nil {
		log.Fatalf("cannot run down migrations: %s", err.Error())
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
	userID, err := c.Register(structs.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}
	masterKey := structs.MasterKey{Key: helpers.GetRandomSrt(32)}
	masterKey.SetHash()

	// setup for Upass
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
		pdata structs.Pdata
	}{
		{name: "Test UPass",
			pdata: structs.Pdata{
				Name:        "test_upass",
				Type:        "upass",
				KeyHash:     base64.StdEncoding.EncodeToString(masterKey.KeyHash[:]),
				PrivateData: base64.StdEncoding.EncodeToString(upassEnc)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := c.PrivateAdd(userID, tc.pdata)
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
	userID, err := c.Register(structs.User{Email: "vasya@test.ru",
		Password: "password"})
	if err != nil {
		t.Fatalf("cant register a test user: %s", err.Error())
	}
	masterKey := structs.MasterKey{Key: helpers.GetRandomSrt(32)}
	masterKey.SetHash()

	// setup for Upass
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
	pdata := structs.Pdata{
		Name:        "test_upass",
		Type:        "upass",
		KeyHash:     base64.StdEncoding.EncodeToString(masterKey.KeyHash[:]),
		PrivateData: base64.StdEncoding.EncodeToString(upassEnc)}
	err = c.PrivateAdd(userID, pdata)
	if err != nil {
		t.Fatalf("cant add pdata to database: %s", err.Error())
	}

	// TestCases
	tt := []struct {
		name string
		want structs.Pdata
	}{
		{
			name: "Test UPass",
			want: pdata,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have, err := c.PrivateGet(userID, tc.want.Type)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if have != tc.want {
				t.Errorf("pdata mismatch: want %v, have %v", tc.want, have)
			}
		})
	}
}
