package repository

import (
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
)

type Factory interface {
	Create() Repository
}

type factory struct {
	baseDir basedir.BaseDir
	index   index.Index
}

func NewFactory(baseDir basedir.BaseDir, index index.Index) *factory {
	return &factory{
		baseDir: baseDir,
		index:   index,
	}
}

func (f *factory) Create() Repository {
	return NewRepository(f.baseDir, f.index)
}
