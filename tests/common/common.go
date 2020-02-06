//
// Copyright 2020 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// +build slowtest

package common

import (
	"context"
	"os"
	"time"

	"github.com/insolar/insolar/log"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
)

func StartDBMS() (pgURL string, cleaner func()) {
	env := os.Getenv("USE_COCKROACH_DB")
	if len(env) > 0 {
		log.Info("Starting CockroachDB...")
		return startCockroachDB()
	} else {
		log.Info("Starting PostgreSQL...")
		return startPostgreSQL()
	}
}

func waitForDBMS(pool *dockertest.Pool, resource *dockertest.Resource, connString string) (url string, cleaner func()) {
	attempt := 0
	ok := false
	for attempt < 20 {
		attempt++
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			log.Infof("[waitForDBMS] pgx.Connect failed: %v, waiting... (attempt %d)", err, attempt)
			time.Sleep(1 * time.Second)
			continue
		}

		_ = conn.Close(context.Background())
		ok = true
		break
	}

	if !ok {
		_ = pool.Purge(resource)
		log.Panicf("[waitForDBMS] couldn't connect to CockroachDB")
	}

	cleanerFunc := func() {
		// purge the container
		err := pool.Purge(resource)
		if err != nil {
			log.Panicf("[waitForDBMS] pool.Purge failed: %v", err)
		}
	}

	return connString, cleanerFunc
}

func startPostgreSQL() (url string, cleaner func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("[startPostgreSQL] dockertest.NewPool failed: %v", err)
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
		log.Panicf("[startPostgreSQL] pool.Run failed: %v", err)
	}

	connString := "postgres://postgres:s3cr3t@" + resource.GetHostPort("5432/tcp") + "/insolar?sslmode=disable"
	return waitForDBMS(pool, resource, connString)
}

func startCockroachDB() (url string, cleaner func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("[startCockroachDB] dockertest.NewPool failed: %v", err)
	}

	opts := &dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        "v19.2.3",
		Cmd:        []string{"start-single-node", "--insecure"},
	}
	resource, err := pool.RunWithOptions(opts)

	if err != nil {
		log.Panicf("[startCockroachDB] pool.Run failed: %v", err)
	}

	connString := "postgres://root@" + resource.GetHostPort("26257/tcp") + "/postgres?sslmode=disable"
	return waitForDBMS(pool, resource, connString)
}
