package main

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/andrewapj/redwing"
	_ "github.com/go-sql-driver/mysql"
)

var (
	//go:embed migrations
	migrations embed.FS
)

func main() {
	db, err := sql.Open("mysql", "redwing:redwing@tcp(127.0.0.1:3306)/redwing")
	if err != nil {
		panic("Can't connect to the DB")
	}
	defer db.Close()

	processed, err := redwing.Migrate(db, redwing.MySQL, migrations, &redwing.Options{Logging: true})

	if err != nil {
		fmt.Printf("Error processing migrations: %v\n", err)
	}

	fmt.Printf("Processed the following migrations: %v", processed)
}
