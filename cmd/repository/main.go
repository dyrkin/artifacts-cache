package main

import (
	"artifacts-cache/pkg/multipart"
	"artifacts-cache/pkg/repository"
	"artifacts-cache/pkg/repository/basedir"
	"artifacts-cache/pkg/repository/database"
	"artifacts-cache/pkg/repository/index"
	"artifacts-cache/pkg/repository/server"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	binaryStreamFactory := multipart.NewBinaryStreamOutFactory(bd)
	repositoryFactory := repository.NewRepositoryFactory(bd, idx, binaryStreamFactory)
	multiRepository := repository.NewMultiRepository(repositoryFactory, 5)
	server.NewServer(8080, multiRepository).Serve()
}
