package redwing

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrDBNotSupported = errors.New("database not supported")
var ErrMigrationTableCreate = errors.New("can not create the migration table")

type Dialect int

const (
	MySQL Dialect = iota
)

type sqlProcessor interface {
	createMigrationTable(db *sql.DB) error
	updateMigrationTable(tx *sql.Tx, fileNum int) error
	getLastMigration(db *sql.DB) (int, error)
}

func Migrate(db *sql.DB, dbType Dialect, path string) ([]int, error) {

	processor, err := setProcessor(dbType)
	if err != nil {
		return []int{}, err
	}

	if err := processor.createMigrationTable(db); err != nil {
		return []int{}, ErrMigrationTableCreate
	}

	fileNum, err := processor.getLastMigration(db)
	if err != nil {
		return []int{}, ErrMigrationTableCreate
	}

	processed := make([]int, 0)
	for {
		fileNum++
		fileName := strconv.Itoa(fileNum) + ".sql"

		fileContent, err := fileContents(filepath.Clean(path + "/" + fileName))
		if err != nil {
			break
		}
		err = executeMigration(db, fileContent, processor, fileNum)
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

func setProcessor(dbType Dialect) (sqlProcessor, error) {

	switch dbType {
	case MySQL:
		return &mySqlProcessor{}, nil
	default:
		return nil, ErrDBNotSupported
	}
}

func fileContents(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}