package repository

import (
	"artifacts-cache/pkg/multipart"
	"artifacts-cache/pkg/repository/basedir"
	"artifacts-cache/pkg/repository/index"
)

type Factory interface {
	Create() Repository
}

type factory struct {
	baseDir             basedir.BaseDir
	index               index.Index
	binaryStreamFactory multipart.BinaryStreamOutFactory
}

func NewRepositoryFactory(baseDir basedir.BaseDir, index index.Index, binaryStreamFactory multipart.BinaryStreamOutFactory) *factory {
	return &factory{
		baseDir:             baseDir,
		index:               index,
		binaryStreamFactory: binaryStreamFactory,
	}
}

func (f *factory) Create() Repository {
	return NewRepository(f.baseDir, f.index, f.binaryStreamFactory)
}
