package redwing

import "testing"

func TestGetContents(t *testing.T) {
	s, err := fileContents("README.MD")
	if err != nil || s == "" {
		t.Fatalf("Expected to read content from a file")
	}
}

func TestMissingFile(t *testing.T) {
	_, err := fileContents("foo.txt")
	if err == nil {
		t.Fatalf("Expected to get an error for a missing file")
	}
}

func TestFindPath(t *testing.T) {
	err := checkPathExists(".")
	if err != nil {
		t.Fatalf("Unable to find path that exists")
	}
}

func TestNotFindMissingPath(t *testing.T) {
	err := checkPathExists("missing/")
	if err == nil {
		t.Fatalf("Should get an error for a missing path")
	}
}
