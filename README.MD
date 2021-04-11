# Redwing

Simple database migrations for Go

## Databases

* MySQL

## Example usage

```go
func main() {
    db, _ := sql.Open("mysql", "user:pass@tcp(host:port)/schema")
    processed, err := redwing.Migrate(db, redwing.MySQL, "migrations/")

    if err != nil {
        fmt.Printf("Error processing migrations: %v", err)
    }

    fmt.Printf("Processed the following migrations: %v", processed)
}

```

* The user is responsible for creating a ```sql.DB``` and closing it when appropriate.
* In this example the migrations are stored under ```migrations/```. Any other valid path can be used.
* The ```Migrate``` function will return a ```[]int```. This contains the numbers representing the migrations
  successfully processed.
* The ```Migrate``` function will also return an error if a migration fails.

## Migrations

Migrations should be stored in a single directory in the following format:

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