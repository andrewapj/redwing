package redwing

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

func fileContents(fileSystem fs.FS, file string) (result string, err error) {
	info, err := fs.Stat(fileSystem, ".")
	if info == nil || err != nil {
		return "", errors.New(fmt.Sprintf("Unable to find file %s. Invalid filesystem", file))
	}

	_ = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() && filepath.Base(filepath.Clean(path)) == file {
			content, _ := fs.ReadFile(fileSystem, path)
			result = string(content)
		}
		return nil
	})

	if result != "" {
		return result, nil
	} else {
		return "", errors.New(fmt.Sprintf("Unable to find file %s", file))
	}
}
