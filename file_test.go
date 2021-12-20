package redwing

import (
	"embed"
	"os"
	"testing"
)

var (
	//go:embed db/mysql
	migrations embed.FS
)

func TestGetContents(t *testing.T) {
	fs := os.DirFS("db/mysql")
	s, err := fileContents(fs, "docker-compose.yml")
	if err != nil || s == "" {
		t.Fatalf("Expected to read content from a file")
	}
}

func TestGetEmbedContents(t *testing.T) {
	s, err := fileContents(migrations, "docker-compose.yml")
	if err != nil || s == "" {
		t.Fatalf("Expected to read content from a file")
	}
}

func TestMissingFile(t *testing.T) {
	fs := os.DirFS(".")
	_, err := fileContents(fs, "Foo")
	if err == nil {
		t.Fatalf("Expected to get an error for a missing file")
	}
}

func TestMissingPath(t *testing.T) {
	fs := os.DirFS("/some/missing/path")
	_, err := fileContents(fs, "Foo")
	if err == nil {
		t.Fatalf("Expected to get an error for a missing path")
	}
}
