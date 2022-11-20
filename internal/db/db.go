package db

import (
	"context"
	"fmt"
	"os"

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

func (c *Connector) checkInit() error {
	if !c.initialized {
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
func (c *Connector) Close() {
	if c.initialized {
		c.Pool.Close()
		c.initialized = false
	}
}

// Register adds user to database and returns new user`s id
func (c *Connector) Register(user structs.User) (int64, error) {
	err := c.checkInit()
	if err != nil {
		return -1, err
	}

	conn, err := c.Pool.Acquire(c.Ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	// Check if user don`t exists
	var counter int
	sql := `SELECT COUNT(id) FROM users WHERE email=$1;`
	err = conn.QueryRow(c.Ctx, sql, user.Email).Scan(&counter)
	if err != nil {
		return -1, fmt.Errorf("failed to query users table: %s", err.Error())
	}
	if counter != 0 {
		return -1, structs.ErrUserAlreadyExists
	}

	// adding new user
	var id int64
	sql = `INSERT INTO users (email, password)
		   VALUES($1, $2)
		   RETURNING id;`
	err = conn.QueryRow(c.Ctx, sql, user.Email, user.Password).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("failed to create user id DB: %s", err.Error())
	}
	return id, nil
}

// CreateTables creates all project tables and populates it`s with data`
// This function mainly used in tests
func (c *Connector) CreateTables() error {
	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	// Recreating schema
	schemaPath := "../../sql/schema.sql"
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("cant read schema.sql: %s", err.Error())
	}
	_, err = conn.Exec(c.Ctx, string(schemaSQL))
	if err != nil {
		return fmt.Errorf("cant create db schema: %s", err.Error())
	}

	// Inserting data
	dataPath := "../../sql/data.sql"
	dataSQL, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("cant read schema.sql: %s", err.Error())
	}
	_, err = conn.Exec(c.Ctx, string(dataSQL))
	if err != nil {
		return fmt.Errorf("cant create db schema: %s", err.Error())
	}

	return nil
}

// DropTables drops all project`s tables.
// This function mainly used in tests.
func (c *Connector) DropTables() error {
	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	dropPath := "../../sql/drop.sql"
	dropSQL, err := os.ReadFile(dropPath)
	if err != nil {
		return fmt.Errorf("cant read schema.sql: %s", err.Error())
	}
	_, err = conn.Exec(c.Ctx, string(dropSQL))
	if err != nil {
		return fmt.Errorf("cant create db schema: %s", err.Error())
	}

	return nil
}

// GetUser search for user by email
func (c *Connector) GetUser(email string) (structs.User, error) {
	err := c.checkInit()
	if err != nil {
		return structs.User{}, err
	}

	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return structs.User{}, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	var user structs.User
	usersSQL := `SELECT id, email, password FROM users WHERE email=$1`
	row := conn.QueryRow(c.Ctx, usersSQL, email)

	switch err := row.Scan(&user.ID, &user.Email, &user.Password); err {
	case pgx.ErrNoRows:
		return structs.User{}, structs.ErrUserAuth
	case nil:
		return user, nil
	default:
		e := fmt.Errorf("unknown error while accesing database: %s", err.Error())
		return structs.User{}, e
	}
}

// PrivateAdd adds private data in database for specific userID
func (c *Connector) PrivateAdd(userID int64, pdata structs.Pdata) error {
	err := c.checkInit()
	if err != nil {
		return err
	}

	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	// get type id
	var typeID int64
	sql := `SELECT id
		    FROM private_types
			WHERE name=$1;`
	row := conn.QueryRow(c.Ctx, sql, pdata.Type)
	err = row.Scan(&typeID)
	if err != nil {
		return fmt.Errorf("cant get private_type id %s", err.Error())
	}

	// Check if pdata don`t exists
	var counter int
	sql = `SELECT COUNT(id)
		   FROM private_data
		   WHERE user_id=$1 AND name=$2;`
	err = conn.QueryRow(c.Ctx, sql, userID, pdata.Name).Scan(&counter)
	if err != nil {
		return fmt.Errorf("failed to query users table: %s", err.Error())
	}
	if counter != 0 {
		return structs.ErrPdataAlreatyEsists
	}

	// inserting data
	sql = `INSERT INTO private_data (name, user_id, type_id, khash_base64, data_base64)
			VALUES($1, $2, $3, $4, $5);`
	_, err = conn.Exec(c.Ctx, sql,
		pdata.Name, userID, typeID, pdata.KeyHash, pdata.PrivateData)
	if err != nil {
		return fmt.Errorf("failed to insert data to db: %s", err.Error())
	}

	return nil
}

// PrivateGet retrive private data from database
func (c *Connector) PrivateGet(userID int64, pname string) (structs.Pdata, error) {

	err := c.checkInit()
	if err != nil {
		return structs.Pdata{}, err
	}

	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return structs.Pdata{}, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	sql := `SELECT a.name, b.name, a.khash_base64, a.data_base64
			FROM private_data AS a
			INNER JOIN private_types AS b
			ON a.type_id=b.id
			WHERE a.user_id=$1 AND a.name=$2;`

	row := conn.QueryRow(c.Ctx, sql, userID, pname)
	var pdata = structs.Pdata{}

	switch err = row.Scan(&pdata.Name, &pdata.Type, &pdata.KeyHash, &pdata.PrivateData); err {
	case pgx.ErrNoRows:
		return structs.Pdata{}, structs.ErrPdataNotFound
	case nil:
		return pdata, nil
	default:
		e := fmt.Errorf("unknown error while accesing database: %s", err.Error())
		return structs.Pdata{}, e
	}
}

func (c *Connector) PrivateUpdate(userID int64, pdata structs.Pdata) error {
	err := c.checkInit()
	if err != nil {
		return err
	}

	conn, err := c.Pool.Acquire(c.Ctx)
	defer conn.Release()
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}

	// inserting data
	sql := `UPDATE private_data
			SET khash_base64 = $1, data_base64 = $2
			WHERE user_id = $3 and name = $4 `
	res, err := conn.Exec(c.Ctx, sql, pdata.KeyHash,
		pdata.PrivateData, userID, pdata.Name)
	if err != nil {
		return fmt.Errorf("failed to insert data to db: %s", err.Error())
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return structs.ErrPdataNotFound
	}
	if rowsAffected > 1 {
		return fmt.Errorf("rows affected: %d", rowsAffected)
	}
	return nil
}
