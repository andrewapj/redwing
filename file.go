package redwing

import "os"

func fileContents(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func checkPathExists(path string) error {
	_, err := os.ReadDir(path)
	if err != nil {
		return ErrPathNotFound
	}
	return nil
}
