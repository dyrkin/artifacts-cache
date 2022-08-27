package repository

import (
	"io"
)

type multiRepository struct {
	workers chan Repository
}

func NewMultiRepository(repositoryFactory Factory, limitThreads int) *multiRepository {
	workers := make(chan Repository, limitThreads+5)

	for i := 0; i < limitThreads; i++ {
		workers <- repositoryFactory.Create()
	}
	return &multiRepository{
		workers: workers,
	}
}

func (r *multiRepository) WriteContent(subset string, name string, content io.Reader) error {
	worker := <-r.workers
	err := worker.WriteContent(subset, name, content)
	r.workers <- worker
	return err
}

func (r *multiRepository) FindContent(subset, filter string) (io.ReadCloser, error) {
	worker := <-r.workers
	reader, err := worker.FindContent(subset, filter)
	r.workers <- worker
	return reader, err
}
