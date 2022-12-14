package client

import (
	"artifacts-cache/pkg/multipart"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	CantSaveResponseToDiskError = errors.New("can't save response to disk")
	FilesNotFoundError          = errors.New("file not found")
)

type PullRepositoryClientFactory interface {
	Create(repository string) PullRepositoryClient
}

type PullRepositoryClient interface {
	Pull(cwd string, subset string, path string) error
}

type pullRepositoryClient struct {
	repository string
}

type pullRepositoryClientFactory struct {
}

func NewPullRepositoryClient(repository string) *pullRepositoryClient {
	return &pullRepositoryClient{
		repository: repository,
	}
}

func NewPullRepositoryClientFactory() PullRepositoryClientFactory {
	return &pullRepositoryClientFactory{}
}

func (f *pullRepositoryClientFactory) Create(repository string) PullRepositoryClient {
	return &pullRepositoryClient{
		repository: repository,
	}
}

func (c *pullRepositoryClient) Pull(cwd string, subset string, filter string) error {
	response, err := http.Get(fmt.Sprintf("%s/pull?subset=%s&filter=%s", c.repository, subset, filter))
	if err != nil {
		return fmt.Errorf("can't send request. error: %w", err)
	}
	if response.StatusCode == http.StatusNotFound {
		return FilesNotFoundError
	} else if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return fmt.Errorf("%w. code: %s. body: %s", InvalidResponseFromRepositoryError, response.Status, body)
	}
	binaryStreamIn := multipart.NewBinaryStreamIn(cwd)
	err = binaryStreamIn.Save(response.Body)
	if err != nil {
		return fmt.Errorf("%w. %s", CantSaveResponseToDiskError, err)
	}
	return nil
}
