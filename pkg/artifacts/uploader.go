package artifacts

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/artifacts/client"
	"net/http"
	"sync"
)

type Uploader interface {
	Upload(cwd string, subset string, paths []string)
}

type uploader struct {
	repositoryClient client.RepositoryClient
	semaphore        chan bool
}

func NewUploader(repositoryClient client.RepositoryClient, maxTreads int) *uploader {
	semaphore := make(chan bool, maxTreads)
	return &uploader{repositoryClient, semaphore}
}

func (u *uploader) Upload(cwd string, subset string, paths []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	for _, path := range paths {
		u.upload(cwd, subset, path, 3, wg)
	}
	wg.Wait()
}

func (u *uploader) upload(cwd string, subset string, path string, attempts int, wg *sync.WaitGroup) {
	u.semaphore <- true
	go func(path string) {
		if err := u.repositoryClient.Push(cwd, subset, path); err != nil {
			if attempts > 0 && (errors.Is(err, http.ErrServerClosed) || errors.Is(err, http.ErrHandlerTimeout)) {
				<-u.semaphore
				u.upload(cwd, subset, path, attempts-1, wg)
			} else {
				log.Error().Msgf("can't upload path [%s] to repository. error [%s]", path, err)
			}
		}
		wg.Done()
	}(path)
	<-u.semaphore
}
