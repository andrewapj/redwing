package test

import (
	"database/sql"
	"github.com/andrewapj/redwing"
	"os"
	"reflect"
	"testing"
)

//Database is an interface that contains helper methods that should be implemented by a struct representing a
//database under test.
type Database interface {
	OpenDB() (db *sql.DB, err error)
	GetDialect() redwing.Dialect
	CreateMigrationTableWithData(db *sql.DB) error
	GetLastMigration(db *sql.DB) (int, error)
	CleanupDB(db *sql.DB)
}

//Migrate checks that a migration works
func Migrate(t *testing.T, td Database) {

	// Given
	db, err := td.OpenDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	td.CleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, td.GetDialect(), os.DirFS("test_migrations/mysql/valid"), &redwing.Options{})
	if err != nil {
		t.Fatalf("Failed to perform migration: %+v", err)
	}

	// Then
	if !reflect.DeepEqual(processed, []int{1, 2}) {
		t.Fatalf("Expected 2 migrations, got %d", len(processed))
	}
	maxId, err := td.GetLastMigration(db)
	if err != nil {
		t.Fatalf("Could not retrieve last migration from table: %v", err)
	}
	if maxId != 2 {
		t.Fatalf("Expected the migration table to contain 2 migrations, got %d", maxId)
	}
}

//MigrateAfterPreviouslyCompleted checks that a new migration works and starts off after the previous migration.
func MigrateAfterPreviouslyCompleted(t *testing.T, td Database) {
	// Given
	db, err := td.OpenDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	td.CleanupDB(db)

	// And: An existing migration for 1.sql
	err = td.CreateMigrationTableWithData(db)
	if err != nil {
		t.Fatalf("Unable to create a migration table with an existing migration: %v", err)
	}

	// When
	processed, err := redwing.Migrate(db, td.GetDialect(), os.DirFS("test_migrations/mysql/valid"), &redwing.Options{})

	// Then: Only 1 migration (2.sql) should be processed
	if err != nil {
		t.Fatalf("Unexpected error processing migration: %v", err)
	}
	if !reflect.DeepEqual(processed, []int{2}) {
		t.Fatalf("Expected only one migration for 2.sql, got %v", processed)
	}
}

//FirstMigrationBroken checks that if the first migration breaks that the correct data is returned and that the
//migrations table is correct
func FirstMigrationBroken(t *testing.T, td Database) {

	// Given
	db, err := td.OpenDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	td.CleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, td.GetDialect(), os.DirFS("test_migrations/mysql/first_broken"), &redwing.Options{})

	// Then
	if err == nil {
		t.Fatal("Expected an error in processing")
	}
	if len(processed) != 0 {
		t.Fatalf("Expected zero migrations, got %d", len(processed))
	}
	maxId, _ := td.GetLastMigration(db)
	if maxId != 0 {
		t.Fatalf("Expected zero migrations in the migrations table, got %d", maxId)
	}
}

//SecondMigrationBroken checks that if the second migration fails that the correct data is returned and that the
//migrations table is correct
func SecondMigrationBroken(t *testing.T, td Database) {

	// Given
	db, err := td.OpenDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	td.CleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, td.GetDialect(), os.DirFS("test_migrations/mysql/second_broken"), &redwing.Options{})

	// Then
	if err == nil {
		t.Fatal("Expected an error when processing the second migration")
	}
	if !reflect.DeepEqual(processed, []int{1}) {
		t.Fatalf("Expected one migration, got %d", len(processed))
	}
	maxId, _ := td.GetLastMigration(db)
	if maxId != 1 {
		t.Fatalf("Expected one migration in the migrations table, got %d", maxId)
	}
}

//NoMigrationsInPath checks that no action is taken if there are no migrations to process
func NoMigrationsInPath(t *testing.T, td Database) {

	// Given
	db, err := td.OpenDB()
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	td.CleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, td.GetDialect(), os.DirFS("test_migrations/mysql/empty"), &redwing.Options{})
	if err != nil {
		t.Fatalf("Failed to perform migration: %+v", err)
	}

	// Then
	if len(processed) != 0 {
		t.Fatalf("Expected zero migrations, got %d", len(processed))
	}
}
