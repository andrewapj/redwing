package redwing

import (
	"database/sql"
	"embed"
	"testing"
)

func TestDialectNotSupported(t *testing.T) {
	_, err := Migrate(&sql.DB{}, Dialect(999), &embed.FS{}, &Options{})
	if err != ErrDialectNotSupported {
		t.Fatalf("Expected a ErrDialectNotSupported error")
	}
}
