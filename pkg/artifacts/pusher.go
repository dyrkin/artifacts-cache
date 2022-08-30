package artifacts

import (
	"artifacts-cache/pkg/artifacts/client"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

type Pusher interface {
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

func (p *pusher) Push(cwd string, subset string, paths []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(paths))
	for _, path := range paths {
		p.push(cwd, subset, path, 3, wg)
	}
	wg.Wait()
}

func (p *pusher) push(cwd string, subset string, path string, attempts int, wg *sync.WaitGroup) {
	p.semaphore <- true
	go func(path string) {
		if err := p.repositoryClient.Push(cwd, subset, path); err != nil {
			if attempts > 0 && (errors.Is(err, http.ErrServerClosed) || errors.Is(err, http.ErrHandlerTimeout)) {
				<-p.semaphore
				p.push(cwd, subset, path, attempts-1, wg)
			} else {
				log.Error().Msgf("can't push path [%s] to repository. error [%s]", path, err)
			}
		}
		wg.Done()
	}(path)
	<-p.semaphore
}
