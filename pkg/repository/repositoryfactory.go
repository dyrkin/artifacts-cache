package repository

import (
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/repository/multipart"
)

type Factory interface {
	Create() Repository
}

type factory struct {
	baseDir             basedir.BaseDir
	index               index.Index
	binaryStreamFactory *multipart.BinaryStreamFactory
}

func NewRepositoryFactory(baseDir basedir.BaseDir, index index.Index, binaryStreamFactory *multipart.BinaryStreamFactory) *factory {
	return &factory{
		baseDir:             baseDir,
		index:               index,
		binaryStreamFactory: binaryStreamFactory,
	}
}

func (f *factory) Create() Repository {
	return NewRepository(f.baseDir, f.index, f.binaryStreamFactory)
}
