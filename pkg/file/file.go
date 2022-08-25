package file

import (
	"errors"
	"io/fs"
	"os"
)

func CreateEmpty(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func RemoveQuiet(path string) {
	_ = os.RemoveAll(path)
}

func CloseQuiet(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}

func Open(path string) (*os.File, error) {
	return os.Open(path)
}
