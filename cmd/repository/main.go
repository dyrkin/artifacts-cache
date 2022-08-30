package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/repository"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/database"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/repository/multipart"
	"gitlab-cache/pkg/repository/server"
)

const (
	BaseDir        = "/Users/unkind/java/projects/artifacts-cache/filedir"
	DatabaseSubDir = "database"
	DataSubdirDir  = "data"
)

func main() {
	log.Level(zerolog.DebugLevel)
	bd := basedir.MustNewBaseDir(BaseDir, DataSubdirDir, DatabaseSubDir)
	err := bd.MakeDatabaseSubdir()
	if err != nil {
		log.Fatal().Msgf("can't create database subdir. error: %s", err)
	}
	db := database.NewDatabase(fmt.Sprintf("file:%s/cache.db?_journal_mode=WAL", bd.DatabaseSubdir()))
	err = db.Connect()
	if err != nil {
		log.Fatal().Msgf("can't connect to db. error: %s", err)
	}
	err = db.Migrate()
	if err != nil {
		log.Fatal().Msgf("can't migrate database. error: %s", err)
	}
	idx := index.NewIndex(db)
	err = idx.Init()
	if err != nil {
		log.Fatal().Msgf("can't initialize index. error: %s", err)
	}
	binaryStreamFactory := multipart.NewBinaryStreamFactory(bd)
	repositoryFactory := repository.NewRepositoryFactory(bd, idx, binaryStreamFactory)
	multiRepository := repository.NewMultiRepository(repositoryFactory, 5)
	server.NewServer(8080, multiRepository).Serve()
}
