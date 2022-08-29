package client

import (
	"errors"
	"fmt"
	"gitlab-cache/pkg/artifacts/compression"
	"gitlab-cache/pkg/repository/file"
	"io"
	"net/http"
	"os"
)

var (
	CantOpenFileError                  = errors.New("can't open file")
	CantCreatePostRequestError         = errors.New("can't create post request")
	CantSendFileToRepositoryError      = errors.New("can't send file to repository")
	InvalidResponseFromRepositoryError = errors.New("invalid response from repository")
)

type repositoryClient struct {
	repository string
	semaphore  chan bool
}

func NewRepositoryClient(repository string, limitThreads int) *repositoryClient {
	return &repositoryClient{
		repository: repository,
		semaphore:  make(chan bool, limitThreads),
	}
}

func (c *repositoryClient) Push(subset string, cwd string, path string) error {
	c.semaphore <- true
	defer func() {
		<-c.semaphore
	}()
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w. %s", CantOpenFileError, err)
	}
	defer f.Close()
	compressingReader := compression.CompressingReader(f)
	defer compressingReader.Close()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/push", c.repository), compressingReader)
	if err != nil {
		return fmt.Errorf("%w. %s", CantCreatePostRequestError, err)
	}
	req.Header.Set("subset", subset)
	req.Header.Set("path", file.RemoveCwd(cwd, path))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w. %s", CantSendFileToRepositoryError, err)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return fmt.Errorf("%w. code: %s. body: %s", InvalidResponseFromRepositoryError, response.Status, body)
	}
	return nil
}
