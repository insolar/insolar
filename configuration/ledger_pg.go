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

package configuration

type PostgreSQL struct {
	URL           string // postgresql:// connection string
	MigrationPath string // path to the directory with migration scripts
}

// Ledger holds configuration for ledger.
type LedgerPg struct {
	// PostgreSQL defines configuration related to PostgreSQL.
	PostgreSQL PostgreSQL

	// IsPostgresBase indicates that heavy uses Postgres as a database
	IsPostgresBase bool
}

// NewLedger creates new default Ledger configuration.
func NewLedgerPg() LedgerPg {
	return LedgerPg{
		PostgreSQL: PostgreSQL{
			URL:           "postgres://postgres@localhost/postgres?sslmode=disable",
			MigrationPath: "migrations",
		},
		IsPostgresBase: false,
	}
}
