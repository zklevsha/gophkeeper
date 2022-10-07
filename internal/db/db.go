package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// Connector encapsulates DB communication logic
type Connector struct {
	Ctx         context.Context
	Pool        *pgxpool.Pool
	DSN         string
	initialized bool
}

func (d *Connector) checkInit() error {
	if !d.initialized {
		err := fmt.Errorf("Connector is not initiliazed (run Connector.Init() to initilize)")
		return err
	}
	return nil
}

// Init connects to DB and initilizes connection pool
func (d *Connector) Init() error {
	p, err := pgxpool.Connect(d.Ctx, d.DSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	d.Pool = p
	d.initialized = true
	return nil
}

// Close closes connection to DB
func (d *Connector) Close() {
	if d.initialized {
		d.Pool.Close()
		d.initialized = false
	}
}

// Register adds user to database and returns new user`s id
func (d *Connector) Register(user structs.User) (int, error) {
	err := d.checkInit()
	if err != nil {
		return -1, err
	}

	conn, err := d.Pool.Acquire(d.Ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	// Check if user don`t exists
	var counter int
	sql := `SELECT COUNT(id) FROM users WHERE email=$1;`
	err = conn.QueryRow(d.Ctx, sql, user.Email).Scan(&counter)
	if err != nil {
		return -1, fmt.Errorf("failed to query users table: %s", err.Error())
	}
	if counter != 0 {
		return -1, structs.ErrUserAlreadyExists
	}

	// adding new user
	var id int
	sql = `INSERT INTO users (email, password)
		   VALUES($1, $2)
		   RETURNING id;`
	err = conn.QueryRow(d.Ctx, sql, user.Email, user.Password).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("failed to create user id DB: %s", err.Error())
	}
	return id, nil
}

// CreateTables creates all project tables.
// This function mainly used in tests
func (d *Connector) CreateTables() error {
	conn, err := d.Pool.Acquire(d.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	usersSQL := `CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		email VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL);;`

	_, err = conn.Exec(d.Ctx, usersSQL)
	if err != nil {
		return fmt.Errorf("cant create users table: %s", err.Error())
	}

	return nil
}

// DropTables drops all project`s tables.
// This function mainly used in tests.
func (d *Connector) DropTables() error {
	conn, err := d.Pool.Acquire(d.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	usersSQL := "DROP TABLE IF EXISTS users"
	_, err = conn.Exec(d.Ctx, usersSQL)
	if err != nil {
		return fmt.Errorf("cant drop counters table: %s", err.Error())
	}

	return nil
}

// GetUser returns user info from database
func (d *Connector) GetUser(id int) (structs.User, error) {
	conn, err := d.Pool.Acquire(d.Ctx)
	defer conn.Release()
	if err != nil {
		return structs.User{}, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	var user structs.User
	usersSQL := `SELECT id, email, password FROM users WHERE id=$1`
	row := conn.QueryRow(d.Ctx, usersSQL, id)

	switch err := row.Scan(&user.Id, &user.Email, &user.Password); err {
	case pgx.ErrNoRows:
		return structs.User{}, structs.ErrUserAuth
	case nil:
		return user, nil
	default:
		e := fmt.Errorf("unknown error while accesing database: %s", err.Error())
		return structs.User{}, e
	}

	// row := conn.QueryRow(d.Ctx, sql, creds.Login)

	// switch err := row.Scan(&id, &password); err {
	// case pgx.ErrNoRows:
	// 	return -1, structs.ErrUserAuth
	// case nil:
	// 	return id, nil
	// default:
	// 	e := fmt.Errorf("unknown error while authenticating user: %s", err.Error())
	// 	return -1, e
	// }

}
