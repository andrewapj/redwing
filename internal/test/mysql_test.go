package test

import (
	"database/sql"
	"github.com/andrewapj/redwing"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

//Implementation of a Database for MySQL
type testMySQLDatabase struct{}

func (m *testMySQLDatabase) OpenDB() (db *sql.DB, err error) {
	return sql.Open("mysql", "redwing:redwing@tcp(127.0.0.1:3306)/redwing")
}

func (m *testMySQLDatabase) GetDialect() redwing.Dialect {
	return redwing.MySQL
}

func (m *testMySQLDatabase) CreateMigrationTableWithData(db *sql.DB) error {
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

func (m *testMySQLDatabase) GetLastMigration(db *sql.DB) (int, error) {
	var maxId int
	err := db.QueryRow("SELECT max(id) FROM redwing_migrate").Scan(&maxId)
	if err != nil {
		return 0, err
	}
	return maxId, nil
}

func (m *testMySQLDatabase) CleanupDB(db *sql.DB) {
	_, _ = db.Exec("DROP TABLE redwing_migrate")
	_, _ = db.Exec("DROP TABLE table1")
	_, _ = db.Exec("DROP TABLE table2")
}

func TestMySQLMigrate(t *testing.T) {
	Migrate(t, &testMySQLDatabase{})
}

func TestMySQLMigrateAfterPreviouslyCompleted(t *testing.T) {
	MigrateAfterPreviouslyCompleted(t, &testMySQLDatabase{})
}

func TestMySQLFirstMigrationBroken(t *testing.T) {
	FirstMigrationBroken(t, &testMySQLDatabase{})
}

func TestMySQLSecondMigrationBroken(t *testing.T) {
	SecondMigrationBroken(t, &testMySQLDatabase{})
}

func TestMySQLNoMigrationsInPath(t *testing.T) {
	NoMigrationsInPath(t, &testMySQLDatabase{})
}
