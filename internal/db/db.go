package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/errs"
)

// Connector encapsulates DB communication logic
type Connector struct {
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
func (c *Connector) Init(ctx context.Context) error {

	p, err := pgxpool.Connect(ctx, c.DSN)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	c.Pool = p
	c.initialized = true
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
func (c *Connector) Register(ctx context.Context, user client.User) (int64, error) {
	err := c.checkInit()
	if err != nil {
		return -1, err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	// Check if user don`t exists
	var counter int
	sql := `SELECT COUNT(id) FROM users WHERE email=$1;`
	err = conn.QueryRow(ctx, sql, user.Email).Scan(&counter)
	if err != nil {
		return -1, fmt.Errorf("failed to query users table: %s", err.Error())
	}
	if counter != 0 {
		return -1, errs.ErrUserAlreadyExists
	}

	// adding new user
	var id int64
	sql = `INSERT INTO users (email, password)
		   VALUES($1, $2)
		   RETURNING id;`
	err = conn.QueryRow(ctx, sql, user.Email, user.Password).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("failed to create user id DB: %s", err.Error())
	}
	return id, nil
}

// GetUser searches for user by email
func (c *Connector) GetUser(ctx context.Context, email string) (client.User, error) {
	err := c.checkInit()
	if err != nil {
		return client.User{}, err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return client.User{}, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	var user client.User
	usersSQL := `SELECT id, email, password FROM users WHERE email=$1`
	row := conn.QueryRow(ctx, usersSQL, email)

	switch err := row.Scan(&user.ID, &user.Email, &user.Password); err {
	case pgx.ErrNoRows:
		return client.User{}, errs.ErrUserAuth
	case nil:
		return user, nil
	default:
		e := fmt.Errorf("unknown error while accesing database: %s", err.Error())
		return client.User{}, e
	}
}

// PrivateAdd adds private data in database for specific userID
func (c *Connector) PrivateAdd(ctx context.Context, userID int64, pdata Pdata) (int64, error) {
	err := c.checkInit()
	if err != nil {
		return -1, err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	// get type id
	var typeID int64
	sql := `SELECT id
		    FROM private_types
			WHERE name=$1;`
	row := conn.QueryRow(ctx, sql, pdata.Type)
	err = row.Scan(&typeID)
	if err != nil {
		return -1, fmt.Errorf("cant get private_type id %s", err.Error())
	}

	// Check if pdata don`t exists
	var counter int
	sql = `SELECT COUNT(id)
		   FROM private_data
		   WHERE user_id=$1 AND name=$2;`
	err = conn.QueryRow(ctx, sql, userID, pdata.Name).Scan(&counter)
	if err != nil {
		return -1, fmt.Errorf("failed to query users table: %s", err.Error())
	}
	if counter != 0 {
		return -1, errs.ErrPdataAlreatyEsists
	}

	var id int64
	sql = `INSERT INTO private_data (name, user_id, type_id, khash_base64, data_base64)
		   VALUES($1, $2, $3, $4, $5)
		   RETURNING id;`
	err = conn.QueryRow(ctx, sql,
		pdata.Name, userID, typeID, pdata.KeyHash, pdata.PrivateData).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("failed to insert data to db: %s", err.Error())
	}

	return id, nil
}

// PrivateGet retrives private data from database
func (c *Connector) PrivateGet(ctx context.Context, userID int64, pdataID int64) (Pdata, error) {

	err := c.checkInit()
	if err != nil {
		return Pdata{}, err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return Pdata{}, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	sql := `SELECT a.id, a.name, b.name, a.khash_base64, a.data_base64
			FROM private_data AS a
			INNER JOIN private_types AS b
			ON a.type_id=b.id
			WHERE a.id=$1 AND a.user_id=$2;`

	row := conn.QueryRow(ctx, sql, pdataID, userID)
	var pdata = Pdata{}

	switch err = row.Scan(&pdata.ID, &pdata.Name, &pdata.Type, &pdata.KeyHash, &pdata.PrivateData); err {
	case pgx.ErrNoRows:
		return Pdata{}, errs.ErrPdataNotFound
	case nil:
		return pdata, nil
	default:
		e := fmt.Errorf("unknown error while accesing database: %s", err.Error())
		return Pdata{}, e
	}
}

// PrivateUpdate updates private data in database
func (c *Connector) PrivateUpdate(ctx context.Context, userID int64, pdata Pdata) error {
	err := c.checkInit()
	if err != nil {
		return err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	// inserting data
	sql := `UPDATE private_data
			SET name = $1, khash_base64 = $2, data_base64 = $3
			WHERE id = $4 AND user_id = $5`
	res, err := conn.Exec(ctx, sql, pdata.Name, pdata.KeyHash,
		pdata.PrivateData, pdata.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to insert data to db: %s", err.Error())
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrPdataNotFound
	}
	if rowsAffected > 1 {
		return fmt.Errorf("rows affected: %d", rowsAffected)
	}
	return nil
}

// PrivateList returns list of user`s private entries
func (c *Connector) PrivateList(ctx context.Context, userID int64, ptype string) ([]PdataEntry, error) {
	err := c.checkInit()
	if err != nil {
		return nil, err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	sql := `SELECT a.id, a.name
			FROM private_data as a
			INNER JOIN private_types as b
			ON a.type_id=b.id
			WHERE a.user_id=$1 AND b.name=$2
			ORDER BY a.name`
	rows, err := conn.Query(ctx, sql, userID, ptype)
	if err != nil {
		return nil, fmt.Errorf("db query error: %s", err.Error())
	}

	var pdataList []PdataEntry
	for rows.Next() {
		var entry PdataEntry
		err := rows.Scan(&entry.ID, &entry.Name)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %s", err.Error())
		}
		pdataList = append(pdataList, entry)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("row scan error: %s", rows.Err())
	}

	return pdataList, nil
}

// PrivateDelete deletes private data from database
func (c *Connector) PrivateDelete(ctx context.Context, userID int64, pdataID int64) error {

	err := c.checkInit()
	if err != nil {
		return err
	}

	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %s", err.Error())
	}
	defer conn.Release()

	sql := `DELETE FROM private_data
			WHERE id=$1 AND user_id=$2`

	res, err := conn.Exec(ctx, sql, pdataID, userID)
	if err != nil {
		return fmt.Errorf("query error: %s", err.Error())
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrPdataNotFound
	}
	if rowsAffected > 1 {
		return fmt.Errorf("oy vey, i`ve deleted %d rows", rowsAffected)
	}
	return nil
}
