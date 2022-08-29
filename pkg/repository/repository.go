package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/repository/multipart"
	"gitlab-cache/pkg/repository/partition"
	"io"
	"sync"
)

var (
	CantRecreatePartitionError = errors.New("can't recreate partition")
	CantFindContentError       = errors.New("can't find content")
	NoContentOnServerError     = errors.New("server has no content")
)

type Repository interface {
	partition.ContentWriter
	FindContent(subset, filter string) (io.ReadCloser, error)
}

type repository struct {
	baseDir             basedir.BaseDir
	index               index.Index
	binaryStreamFactory *multipart.BinaryStreamFactory
	partition           partition.Partition
	mutex               *sync.Mutex
}

func NewRepository(baseDir basedir.BaseDir, index index.Index, binaryStreamFactory *multipart.BinaryStreamFactory) *repository {
	return &repository{
		baseDir:             baseDir,
		index:               index,
		binaryStreamFactory: binaryStreamFactory,
		mutex:               &sync.Mutex{},
	}
}

func (r *repository) WriteContent(subset, path string, content io.Reader) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	err := r.recreatePartition()
	if err != nil {
		return fmt.Errorf("%w. %s", CantRecreatePartitionError, err)
	}
	return r.partition.WriteContent(subset, path, content)
}

func (r *repository) FindContent(subset, filter string) (io.ReadCloser, error) {
	contentEmplacement, err := r.index.FindContentEmplacement(subset, filter)
	if err != nil {
		return nil, fmt.Errorf("%w. %s", CantFindContentError, err)
	}
	if len(contentEmplacement.Emplacements) == 0 {
		return nil, NoContentOnServerError
	}
	return r.binaryStreamFactory.Create(contentEmplacement)
}

func (r *repository) recreatePartition() error {
	if r.partition == nil || r.partition.IsClosed() {
		uuId := uuid.New().String()
		r.partition = partition.NewPartition(uuId, r.baseDir, r.index)
		return r.partition.Open()
	}
	return nil
}
