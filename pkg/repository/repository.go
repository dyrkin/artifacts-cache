package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/repository/partition"
	"io"
	"sync"
)

var (
	CantRecreatePartitionError = errors.New("can't recreate partition")
)

type Repository interface {
	partition.ContentWriter
}

type repository struct {
	baseDir   basedir.BaseDir
	index     index.Index
	partition partition.Partition
	mutex     *sync.Mutex
}

func NewRepository(baseDir basedir.BaseDir, index index.Index) *repository {
	return &repository{
		baseDir: baseDir,
		index:   index,
		mutex:   &sync.Mutex{},
	}
}

func (r *repository) WriteContent(subset, name string, content io.Reader) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	err := r.recreatePartition()
	if err != nil {
		return fmt.Errorf("%w. %s", CantRecreatePartitionError, err)
	}
	return r.partition.WriteContent(subset, name, content)
}

func (r *repository) recreatePartition() error {
	if r.partition == nil || r.partition.IsClosed() {
		uuId := uuid.New().String()
		r.partition = partition.NewPartition(uuId, r.baseDir, r.index)
		return r.partition.Open()
	}
	return nil
}
