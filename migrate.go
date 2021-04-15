package redwing

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strconv"
	"strings"
)

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

//ErrDialectNotSupported is an error returned when a dialect is not supported
var ErrDialectNotSupported = errors.New("dialect not supported")

//ErrPathNotFound is an error returned when the migrations path could not be found
var ErrPathNotFound = errors.New("path not found")

//Migrate starts a database migration given the valid sql.DB,
// a database Dialect and a path containing the migrations.
func (r *redwing) Migrate() ([]int, error) {

	processor, err := setProcessor(r.dialect)
	if err != nil {
		return []int{}, err
	}

	err = checkPathExists(r.path)
	if err != nil {
		return nil, err
	}

	if err := processor.createMigrationTable(r.db); err != nil {
		return []int{}, err
	}

	fileNum, err := processor.getLastMigration(r.db)
	if err != nil {
		return []int{}, err
	}

	processed := make([]int, 0)
	for {
		fileNum++
		fileName := strconv.Itoa(fileNum) + ".sql"

		fileContent, err := fileContents(filepath.Clean(r.path + "/" + fileName))
		if err != nil {
			break
		}
		err = executeMigration(r.db, fileContent, processor, fileNum)
		if err != nil {
			return processed, err
		}
		processed = append(processed, fileNum)
	}

	return processed, nil
}

func executeMigration(db *sql.DB, content string, processor sqlProcessor, fileNum int) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := strings.Split(content, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		if _, err := tx.Exec(statement); err != nil {
			return err
		}
		if err := processor.updateMigrationTable(tx, fileNum); err != nil {
			return err
		}
	}

	return tx.Commit()
}
