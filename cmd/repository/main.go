package main

import (
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/database"
	"gitlab-cache/pkg/repository"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/server"
	"io"
	"sync"
)

type mockRepository struct {
	id int
	wg *sync.WaitGroup
}

func (m *mockRepository) WriteContent(key string, content io.Reader) error {
	log.Info().Msgf("processing: %s by %d", key, m.id)
	m.wg.Done()
	return nil
}

func main() {
	bd := basedir.MustNewBaseDir("/Users/unkind/java/projects/gitlab-cache/filedir")
	db := database.NewDatabase("postgres://pqgotest:password@localhost/pqgotest")
	err := db.Connect()
	if err != nil {
		log.Fatal().Msgf("can't connect to db. error: %s", err)
	}
	idx := index.NewIndex(db)
	err = idx.Init()
	if err != nil {
		log.Fatal().Msgf("can't initialize index. error: %s", err)
	}
	repositoryFactory := repository.NewRepositoryFactory(bd, idx)
	server.NewServer(8080, repository.NewMultiRepository(repositoryFactory, 5)).Serve()
}
