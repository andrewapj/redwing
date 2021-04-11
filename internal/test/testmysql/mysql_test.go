package testmysql

import (
	"database/sql"
	"github.com/andrewapj/redwing"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
)

var migrationsBase = "../test_migrations/mysql/"
var connectionString = "redwing:redwing@tcp(127.0.0.1:3306)/redwing"

func TestMySQLMigrate(t *testing.T) {

	// Given
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	cleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, redwing.MySQL, migrationsBase+"valid")
	if err != nil {
		t.Fatalf("Failed to perform migration: %+v", err)
	}

	// Then
	if !reflect.DeepEqual(processed, []int{1, 2}) {
		t.Fatalf("Expected 2 migrations, got %d", len(processed))
	}
	maxId, err := getLastMigration(db)
	if err != nil {
		t.Fatalf("Could not retrieve last migration from table: %v", err)
	}
	if maxId != 2 {
		t.Fatalf("Expected the migration table to contain 2 migrations, got %d", maxId)
	}
}

func TestMigrateFromPreviouslyCompleted(t *testing.T) {
	// Given
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	cleanupDB(db)

	// And: An existing migration for 1.sql
	err = createMigrationTableWithData(db)
	if err != nil {
		t.Fatalf("Unable to create a migration table with an existing migration: %v", err)
	}

	// When
	processed, err := redwing.Migrate(db, redwing.MySQL, migrationsBase+"first_broken")

	// Then: Only 1 migration (2.sql) should be processed
	if err != nil {
		t.Fatalf("Unexpected error processing migration: %v", err)
	}
	if !reflect.DeepEqual(processed, []int{2}) {
		t.Fatalf("Expected only one migration for 2.sql, got %v", processed)
	}
}

func TestFirstMigrationBroken(t *testing.T) {

	// Given
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	cleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, redwing.MySQL, migrationsBase+"first_broken")

	// Then
	if err == nil {
		t.Fatal("Expected an error in processing")
	}
	if len(processed) != 0 {
		t.Fatalf("Expected zero migrations, got %d", len(processed))
	}
	maxId, _ := getLastMigration(db)
	if maxId != 0 {
		t.Fatalf("Expected zero migrations in the migrations table, got %d", maxId)
	}
}

func TestSecondMigrationBroken(t *testing.T) {

	// Given
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	cleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, redwing.MySQL, migrationsBase+"second_broken")

	// Then
	if err == nil {
		t.Fatal("Expected an error when processing the second migration")
	}
	if !reflect.DeepEqual(processed, []int{1}) {
		t.Fatalf("Expected one migration, got %d", len(processed))
	}
	maxId, _ := getLastMigration(db)
	if maxId != 1 {
		t.Fatalf("Expected one migration in the migrations table, got %d", maxId)
	}
}

func TestNoMigrationsInPath(t *testing.T) {

	// Given
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	cleanupDB(db)

	// When
	processed, err := redwing.Migrate(db, redwing.MySQL, migrationsBase+"empty")
	if err != nil {
		t.Fatalf("Failed to perform migration: %+v", err)
	}

	// Then
	if len(processed) != 0 {
		t.Fatalf("Expected zero migrations, got %d", len(processed))
	}
}

// Helper functions for MySQL

func createMigrationTableWithData(db *sql.DB) error {
	_, err := db.Exec(`create table if not exists redwing_migrate
(
	id int not null,
	modified timestamp null,
	constraint redwing_migrate_pk primary key (id)
);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO redwing_migrate (id, modified) values (1, NOW())`)
	if err != nil {
		return err
	}
	return nil
}

func getLastMigration(db *sql.DB) (int, error) {
	var maxId int
	err := db.QueryRow("SELECT max(id) FROM redwing_migrate").Scan(&maxId)
	if err != nil {
		return 0, err
	}
	return maxId, nil
}

func cleanupDB(db *sql.DB) {
	_, _ = db.Exec("TRUNCATE TABLE redwing_migrate")
	_, _ = db.Exec("DROP TABLE table1")
	_, _ = db.Exec("DROP TABLE table2")
}
