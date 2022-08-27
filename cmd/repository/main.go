package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/database"
	"gitlab-cache/pkg/repository"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/server"
)

const (
	DbName     = "cache"
	DbLogin    = "cache"
	DbPassword = "pwd"
	DbAddress  = "docker.home:5432"
	BaseDir    = "/Users/unkind/java/projects/gitlab-cache/filedir"
)

func main() {
	log.Level(zerolog.DebugLevel)
	bd := basedir.MustNewBaseDir(BaseDir)
	db := database.NewDatabase(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbLogin, DbPassword, DbAddress, DbName))
	err := db.Connect()
	if err != nil {
		log.Fatal().Msgf("can't connect to db. error: %s", err)
	}
	err = db.Migrate()
	if err != nil {
		log.Fatal().Msgf("can't migrate database. error: %s", err)
	}
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
