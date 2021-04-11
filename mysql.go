package redwing

import "database/sql"

var sqlCreate = `
create table if not exists redwing_migrate
(
	id int not null,
	modified timestamp null,
	constraint redwing_migrate_pk primary key (id)
);`
var sqlInsert = `INSERT INTO redwing_migrate (id, modified) values (?, NOW())`

type mySqlProcessor struct{}

func (m *mySqlProcessor) createMigrationTable(db *sql.DB) error {
	_, err := db.Exec(sqlCreate)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySqlProcessor) updateMigrationTable(tx *sql.Tx, fileNum int) error {

	_, err := tx.Exec(sqlInsert, fileNum)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySqlProcessor) getLastMigration(db *sql.DB) (int, error) {
	var maxId sql.NullInt32
	err := db.QueryRow("SELECT max(id) FROM redwing_migrate").Scan(&maxId)
	if err != nil {
		return 0, err
	}
	if maxId.Valid {
		return int(maxId.Int32), nil
	} else {
		return 0, nil
	}
}
