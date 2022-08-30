package artifacts

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/artifacts/client"
	"net/http"
	"sync"
)

type Puller interface {
	Pull(cwd string, subset string, filter string)
}

type puller struct {
	repositoryClientFactory client.PullRepositoryClientFactory
	repositories            []string
}

func NewPuller(repositoryClientFactory client.PullRepositoryClientFactory, repositories []string) *puller {
	return &puller{repositoryClientFactory, repositories}
}

func (p *puller) Pull(cwd string, subset string, filter string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(p.repositories))
	for _, repository := range p.repositories {
		repositoryClient := p.repositoryClientFactory.Create(repository)
		p.pull(repositoryClient, cwd, subset, filter, 3, wg)
	}
	wg.Wait()
}

func (p *puller) pull(repositoryClient client.PullRepositoryClient, cwd string, subset string, filter string, attempts int, wg *sync.WaitGroup) {
	go func() {
		if err := repositoryClient.Pull(cwd, subset, filter); err != nil {
			if attempts > 0 && (errors.Is(err, http.ErrServerClosed) || errors.Is(err, http.ErrHandlerTimeout)) {
				p.pull(repositoryClient, cwd, subset, filter, attempts-1, wg)
			} else {
				log.Error().Msgf("can't get files from repository %s. error %s", filter, err)
			}
		}
		wg.Done()
	}()
}
