// +build slowtest

package common

import (
	"context"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
)

func StartPostgreSQL() (pgURL string, cleaner func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("[StartPostgreSQL] dockertest.NewPool failed: %v", err)
	}

	resource, err := pool.Run(
		"postgres",
		"11",
		[]string{
			"POSTGRES_DB=insolar",
			"POSTGRES_PASSWORD=s3cr3t",
		},
	)
	if err != nil {
		log.Panicf("[StartPostgreSQL] pool.Run failed: %v", err)
	}

	// PostgreSQL needs some time to start.
	// Port forwarding always works, thus net.Dial can't be used here.
	connString := "postgres://postgres:s3cr3t@" + resource.GetHostPort("5432/tcp") + "/insolar?sslmode=disable"
	attempt := 0
	ok := false
	for attempt < 20 {
		attempt++
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			log.Infof("[StartPostgreSQL] pgx.Connect failed: %v, waiting... (attempt %d)", err, attempt)
			time.Sleep(1 * time.Second)
			continue
		}

		_ = conn.Close(context.Background())
		ok = true
		break
	}

	if !ok {
		_ = pool.Purge(resource)
		log.Panicf("[StartPostgreSQL] couldn't connect to PostgreSQL")
	}

	cleanerFunc := func() {
		// purge the container
		err := pool.Purge(resource)
		if err != nil {
			log.Panicf("[StartPostgreSQL] pool.Purge failed: %v", err)
		}
	}

	return connString, cleanerFunc
}
