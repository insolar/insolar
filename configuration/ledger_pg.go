// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

type PostgreSQL struct {
	URL           string // postgresql:// connection string
	MigrationPath string // path to the directory with migration scripts
}

// Ledger holds configuration for ledger.
type LedgerPg struct {
	// PostgreSQL defines configuration related to PostgreSQL.
	PostgreSQL PostgreSQL
}

// NewLedger creates new default Ledger configuration.
func NewLedgerPg() LedgerPg {
	return LedgerPg{
		PostgreSQL: PostgreSQL{
			URL:           "postgres://postgres@localhost/postgres?sslmode=disable",
			MigrationPath: "migrations",
		},
	}
}
