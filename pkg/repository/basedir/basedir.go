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
	dir string
}

func MustNewBaseDir(dir string) *baseDir {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal().Msgf("can't create base dir [%s]. error: %s", dir, err)
	}
	return &baseDir{dir}
}

func (d *baseDir) Dir() string {
	return d.dir
}

func (d *baseDir) MakeSubdirByUUID(uuid string) (string, error) {
	parts := strings.Split(uuid, "-")
	subdir := path.Join(d.dir, parts[1], parts[2])
	if err := os.MkdirAll(subdir, os.ModePerm); err != nil {
		return "", fmt.Errorf("can't create subdir [%s]. error: %w", subdir, err)
	}
	return subdir, nil
}
