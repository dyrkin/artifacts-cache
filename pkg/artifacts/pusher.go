package artifacts

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/artifacts/client"
	"net/http"
	"sync"
)

type Uploader interface {
	Push(cwd string, subset string, paths []string)
}

type pusher struct {
	repositoryClient client.PushRepositoryClient
	semaphore        chan bool
}

func NewPusher(repositoryClient client.PushRepositoryClient, maxTreads int) *pusher {
	semaphore := make(chan bool, maxTreads)
	return &pusher{repositoryClient, semaphore}
}

func (u *pusher) Push(cwd string, subset string, paths []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	for _, path := range paths {
		u.push(cwd, subset, path, 3, wg)
	}
	wg.Wait()
}

func (u *pusher) push(cwd string, subset string, path string, attempts int, wg *sync.WaitGroup) {
	u.semaphore <- true
	go func(path string) {
		if err := u.repositoryClient.Push(cwd, subset, path); err != nil {
			if attempts > 0 && (errors.Is(err, http.ErrServerClosed) || errors.Is(err, http.ErrHandlerTimeout)) {
				<-u.semaphore
				u.push(cwd, subset, path, attempts-1, wg)
			} else {
				log.Error().Msgf("can't push path [%s] to repository. error [%s]", path, err)
			}
		}
		wg.Done()
	}(path)
	<-u.semaphore
}
