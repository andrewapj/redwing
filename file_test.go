package redwing

import (
	"os"
	"testing"
)

func TestGetContents(t *testing.T) {
	fs := os.DirFS(".")
	s, err := fileContents(fs, "README.MD")
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
