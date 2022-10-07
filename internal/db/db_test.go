package db

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

const dsn = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper_test"

var ctx = context.Background()

func setUp() Connector {
	d := Connector{Ctx: ctx, DSN: dsn}
	err := d.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}
	err = d.CreateTables()
	if err != nil {
		log.Fatalf("Cant create tables: %s", err.Error())
	}
	return d
}

func tearDown(d Connector) {
	err := d.DropTables()
	if err != nil {
		log.Fatalf("DropTables have returned ad error: %s", err.Error())
	}
}

func TestRegister(t *testing.T) {
	// setup
	d := setUp()
	defer d.Close()
	defer tearDown(d)

	// run tests
	t.Run("Register", func(t *testing.T) {
		user := structs.User{Email: "vasa@test.ru", Password: "test"}
		id, err := d.Register(user)
		if err != nil {
			t.Fatalf("Register() have returned an error: %s", err.Error())
		}

		conn, err := d.Pool.Acquire(d.Ctx)
		if err != nil {
			t.Fatalf("failed to acquire connection: %s", err.Error())
		}
		defer conn.Release()

		var email string
		sql := `SELECT email FROM users WHERE id=$1`
		err = conn.QueryRow(d.Ctx, sql, id).Scan(&email)
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
	d := setUp()
	defer d.Close()
	defer tearDown(d)

	// Get user that exists
	want := structs.User{Email: "vasya@test.ru", Password: "secret"}
	id, err := d.Register(want)
	if err != nil {
		t.Fatalf("cant register new user: %s", err.Error())
	}
	want.Id = id
	have, err := d.GetUser(id)
	if err != nil {
		t.Fatalf("cant get user: %s", err.Error())
	}
	if want != have {
		t.Errorf("user mismatch: have: %v want: %v", have, want)
	}

	// Get user that does`t exists
	_, err = d.GetUser(1004932)
	if !errors.Is(err, structs.ErrUserAuth) {
		t.Errorf("err != structs.ErrUserAuth: %v", err)
	}
}
