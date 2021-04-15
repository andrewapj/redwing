package redwing

import "testing"

func TestDialectNotSupported(t *testing.T) {
	r := redwing{dialect: 99}
	_, err := r.Migrate()
	if err != ErrDialectNotSupported {
		t.Fatalf("Expected a ErrDialectNotSupported error")
	}
}

func TestPathNotFound(t *testing.T) {
	r := redwing{dialect: MySQL, path: "missing/"}
	_, err := r.Migrate()

	if err != ErrPathNotFound {
		t.Fatalf("Expected a path not found error for a missing path")
	}
}
