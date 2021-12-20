package redwing

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
)

//ErrDialectNotSupported is an error returned when a dialect is not supported
var ErrDialectNotSupported = errors.New("dialect not supported")

//Options allows the user to provide extra optional parameters for the migration
type Options struct {
	Logging bool
}

//Migrate starts a database migration
func Migrate(db *sql.DB, dialect Dialect, f fs.FS, options *Options) ([]int, error) {

	processor, err := setProcessor(dialect)
	if err != nil {
		return []int{}, err
	}

	if err := processor.createMigrationTable(db); err != nil {
		PrintLog("Unable to create migration table", options)
		return []int{}, err
	}

	fileNum, err := processor.getLastMigration(db)
	if err != nil {
		PrintLog("Unable to retrieve last migration", options)
		return []int{}, err
	}
	PrintLog(fmt.Sprintf("Found %d previous migrations", fileNum), options)

	PrintLog(fmt.Sprintf("Processing any valid migrations"), options)
	processed := make([]int, 0)
	for {
		fileNum++
		fileName := strconv.Itoa(fileNum) + ".sql"

		fileContent, err := fileContents(f, fileName)
		if err != nil {
			PrintLog(err.Error(), options)
			break
		}
		err = executeMigration(db, fileContent, processor, fileNum)
		if err != nil {
			return processed, err
		}
		processed = append(processed, fileNum)
	}

	PrintLog(fmt.Sprintf("Processed %d new migrations", len(processed)), options)
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
