package main

import (
	"artifacts-cache/pkg/multipart"
	"artifacts-cache/pkg/repository"
	"artifacts-cache/pkg/repository/basedir"
	"artifacts-cache/pkg/repository/cleaner"
	"artifacts-cache/pkg/repository/database"
	"artifacts-cache/pkg/repository/index"
	"artifacts-cache/pkg/repository/server"
	"artifacts-cache/pkg/repository/system"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

const (
	DatabaseSubDir                 = "database"
	DataSubdirDir                  = "data"
	TargetAvailableSpaceInPercents = 10
)

func main() {
	base := ""
	if len(os.Args) == 2 && os.Args[1] != "" {
		base = os.Args[1]
	} else {
		log.Fatal().Msgf("usage: %s <base_dir>", path.Base(os.Args[0]))
	}
	log.Info().Msgf("base dir: %s", base)
	log.Level(zerolog.DebugLevel)
	bd := basedir.MustNewBaseDir(base, DataSubdirDir, DatabaseSubDir)
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
	diskAllocation := system.NewDiskAllocation(bd)
	diskAllocation.Start()
	cleaner := cleaner.NewCleaner(bd, idx, diskAllocation, TargetAvailableSpaceInPercents)
	cleaner.Schedule()
	binaryStreamFactory := multipart.NewBinaryStreamOutFactory(bd)
	repositoryFactory := repository.NewRepositoryFactory(bd, idx, binaryStreamFactory)
	multiRepository := repository.NewMultiRepository(repositoryFactory, 5)
	server.NewServer(8080, multiRepository).Serve()
}
