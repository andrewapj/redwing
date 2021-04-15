package redwing

import "database/sql"

//redwing is the main struct for carrying out the migration
type redwing struct {
	db      *sql.DB
	path    string
	dialect Dialect
}

//New creates a redwing struct with mandatory parameters
func New(db *sql.DB, dialect Dialect, path string) *redwing {
	return &redwing{
		db:      db,
		dialect: dialect,
		path:    path,
	}
}
