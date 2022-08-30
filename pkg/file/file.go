package file

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return info != nil && !errors.Is(err, fs.ErrNotExist) && info.IsDir()
}

func FindFilesInDir(dir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func MkdirAllForFile(path string) error {
	return os.MkdirAll(filepath.Dir(path), os.ModePerm)
}
