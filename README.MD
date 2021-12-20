# Redwing

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewapj/redwing)](https://goreportcard.com/report/github.com/andrewapj/redwing)

Simple database migrations for Go

## Installation

```go get github.com/andrewapj/redwing```

## Databases

* MySQL

## Example usage

```go
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

```

* The user is responsible for creating a ```sql.DB``` and closing it when appropriate.
* In this example the migrations are stored under ```migrations/```. Any other valid path can be used.
* The ```Migrate``` function will return a ```[]int```. This contains the numbers representing the migrations
  successfully processed.
* The ```Migrate``` function will also return an error if a migration fails.

## Migrations

Migrations should be stored in the following format:

```
1.sql
2.sql
3.sql
etc...
```

* Each migration will be processed sequentially in turn. Any migration that is successfully processed will be returned
  within the ```[]int```

* If the system encounters a migration that can not be processed it will return an error in addition to the processed
  migrations, if any.

* Migrations that have already been processed are stored in the migration table and will not be repeated.

## Migrations table

* Database migrations are stored within the ```redwing_migrate``` table.
* This contains two columns:
    * An id related to the sequential number of the migration.
    * A timestamp of when the migration was applied.