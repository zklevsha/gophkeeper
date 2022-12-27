package gserver

import (
	"context"
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/db"
)

// testServer represents gRPC server used in tests
type testServer struct {
	auth authServer
	pdata pdataServer
}

const dsnDefault = "postgres://gophkeeper:gophkeeper@localhost:5532/gophkeeper_test?sslmode=disable"

func setUp() testServer {
	ctx := context.Background()
	key := "secret"
	// added so github action will be able to connect to test database
	var dsn = os.Getenv("GK_DB_TEST_DSN")
	if (dsn == "" ){
		dsn = dsnDefault
	}
	// connecting to DB
	c := db.Connector{Ctx: ctx, DSN: dsn}
	err := c.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}

	// running migrations
	err = db.RunMigrations(dsn, "up")
	if err != nil {
		log.Fatalf("cannot run up migrations: %s", err.Error())
	}

	return testServer{
		auth: authServer{db: c, key: key},
		pdata: pdataServer{db: c, key: key},
	}
}

func tearDown(s testServer) {
	err := db.RunMigrations(s.auth.db.DSN, "down")
	if err != nil {
		log.Fatalf("cannot run down migrations: %s", err.Error())
	}
	s.auth.db.Close()
}
