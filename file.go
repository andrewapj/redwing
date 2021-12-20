package redwing

import (
	"io/fs"
)

func fileContents(fileSystem fs.FS, file string) (string, error) {
	contents, err := fs.ReadFile(fileSystem, file)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
