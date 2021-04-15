package redwing

import (
	"database/sql"
)

type sqlProcessor interface {
	createMigrationTable(db *sql.DB) error
	updateMigrationTable(tx *sql.Tx, fileNum int) error
	getLastMigration(db *sql.DB) (int, error)
}

func setProcessor(dbType Dialect) (sqlProcessor, error) {
	switch dbType {
	case MySQL:
		return &mySqlProcessor{}, nil
	default:
		return nil, ErrDialectNotSupported
	}
}
