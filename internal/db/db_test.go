package db

import (
	"context"
	"log"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

const dsn = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper_test"

var ctx = context.Background()

// func tableExists(tname string, d Connector) (bool, error) {
// 	conn, err := d.Pool.Acquire(d.Ctx)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to acquire connection: %s", err.Error())
// 	}
// 	defer conn.Release()

// 	var count int
// 	sql := `SELECT COUNT(table_name) FROM information_schema.tables WHERE table_name=$1;`
// 	row := conn.QueryRow(d.Ctx, sql, tname)
// 	err = row.Scan(&count)
// 	if err != nil {
// 		return false, fmt.Errorf("sql %s have returned an error: %s", sql, err.Error())
// 	}

// 	switch count {
// 	case 0:
// 		return false, nil
// 	case 1:
// 		return true, nil
// 	default:
// 		return false, fmt.Errorf("sql %s have returned more than 1 COUNT: %d", sql, count)
// 	}
// }

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
