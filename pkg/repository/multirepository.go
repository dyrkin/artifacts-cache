package repository

import (
	"io"
)

type multiRepository struct {
	writeRepositories chan Repository
	readRepository    Repository
}

func NewMultiRepository(repositoryFactory Factory, limitThreads int) *multiRepository {
	writeRepositories := make(chan Repository, limitThreads+1)

	for i := 0; i < limitThreads; i++ {
		writeRepositories <- repositoryFactory.Create()
	}
	readRepository := repositoryFactory.Create()
	return &multiRepository{
		writeRepositories: writeRepositories,
		readRepository:    readRepository,
	}
}

func (r *multiRepository) WriteContent(subset string, path string, content io.Reader) error {
	repository := <-r.writeRepositories
	err := repository.WriteContent(subset, path, content)
	r.writeRepositories <- repository
	return err
}

func (r *multiRepository) FindContent(subset, filter string) (io.ReadCloser, error) {
	return r.readRepository.FindContent(subset, filter)
}
