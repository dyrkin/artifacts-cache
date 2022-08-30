package basedir

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
)

type BaseDir interface {
	Dir() string
	MakeSubdirByUUID(uuid string) (string, error)
}

type baseDir struct {
	dir            string
	dataSubdir     string
	databaseSubdir string
}

func MustNewBaseDir(dir string, dataSubdir string, databaseSubdir string) *baseDir {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal().Msgf("can't create base dir [%s]. error: %s", dir, err)
	}
	return &baseDir{dir: dir, dataSubdir: dataSubdir, databaseSubdir: databaseSubdir}
}

func (d *baseDir) MakeSubdirByUUID(uuid string) (string, error) {
	parts := strings.Split(uuid, "-")
	subdir := path.Join(d.dir, d.dataSubdir, parts[1], parts[2])
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		return "", fmt.Errorf("can't create subdir [%s]. error: %w", subdir, err)
	}
	return subdir, nil
}

func (d *baseDir) MakeDatabaseSubdir() error {
	subdir := d.DatabaseSubdir()
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		return fmt.Errorf("can't create subdir [%s]. error: %w", subdir, err)
	}
	return nil
}

func (d *baseDir) Dir() string {
	return d.dir
}

func (d *baseDir) DatabaseSubdir() string {
	return path.Join(d.dir, d.databaseSubdir)
}

func (d *baseDir) DataSubdir() string {
	return path.Join(d.dir, d.dataSubdir)
}
