package client

import (
	"errors"
	"fmt"
	"gitlab-cache/pkg/artifacts/client/url"
	"gitlab-cache/pkg/artifacts/compression"
	"gitlab-cache/pkg/file"
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

type RepositoryClient interface {
	Push(cwd string, subset string, path string) error
}

type repositoryClient struct {
	repositoryUrlRotator url.Rotator
}

func NewRepositoryClient(repositoryUrlRotator url.Rotator) *repositoryClient {
	return &repositoryClient{
		repositoryUrlRotator: repositoryUrlRotator,
	}
}

func (c *repositoryClient) Push(cwd string, subset string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w. %s", CantOpenFileError, err)
	}
	defer f.Close()
	compressingReader := compression.CompressingReader(f)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/push", c.repositoryUrlRotator.Next()), compressingReader)
	if err != nil {
		return fmt.Errorf("%w. %s", CantCreatePostRequestError, err)
	}
	req.Header.Set("subset", subset)
	req.Header.Set("path", file.RemoveCwd(cwd, path))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't send request. error: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		defer response.Body.Close()
		return fmt.Errorf("%w. code: %s. body: %s", InvalidResponseFromRepositoryError, response.Status, body)
	}
	return nil
}
