package client

import (
	"artifacts-cache/pkg/artifacts/client/url"
	"artifacts-cache/pkg/compression"
	"artifacts-cache/pkg/file"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	CantOpenFileError                  = errors.New("can't open file")
	CantCreatePostRequestError         = errors.New("can't create post request")
	InvalidResponseFromRepositoryError = errors.New("invalid response from repository")
)

type PushRepositoryClient interface {
	Push(cwd string, subset string, path string) error
}

type pushRepositoryClient struct {
	repositoryUrlRotator url.Rotator
}

func NewPushRepositoryClient(repositoryUrlRotator url.Rotator) *pushRepositoryClient {
	return &pushRepositoryClient{
		repositoryUrlRotator: repositoryUrlRotator,
	}
}

func (c *pushRepositoryClient) Push(cwd string, subset string, path string) error {
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
