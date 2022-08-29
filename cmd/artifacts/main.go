package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	log.Level(zerolog.DebugLevel)
	scope := os.Getenv("ARTIFACTS_SCOPE_ID")
	if scope == "" {
		log.Fatal().Msg("ARTIFACTS_SCOPE_ID is not set")
	}
	repositoriesCommaSeparated := os.Getenv("ARTIFACTS_REPOSITORIES")
	if repositoriesCommaSeparated == "" {
		log.Fatal().Msg("ARTIFACTS_REPOSITORIES is not set")
	}
	repositories := strings.Split(repositoriesCommaSeparated, ",")
	if len(os.Args) == 3 && (os.Args[1] == "pull" || os.Args[1] == "push") && os.Args[2] != "" {
		cmd := os.Args[1]
		pattern := os.Args[2]
		if strings.Contains(pattern, "..") {
			log.Fatal().Msgf("pattern can't contain '..'")
		}
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Msgf("can't get current directory. error: %s", err)
		}
		switch cmd {
		case "pull":
			pull(cwd, pattern)
		case "push":
			push(cwd, pattern, repositories)
		}
	} else {
		log.Fatal().Msgf("usage: %s <pull|push> <path|--all>", path.Base(os.Args[0]))
	}
}

func push(cwd string, pattern string, repositories []string) {
	switch pattern {
	case "--all":
	default:
		files, err := filepath.Glob(path.Join(cwd, pattern))
		if err != nil {
			log.Fatal().Msgf("can't find files using pattern [%s]. error: %s", pattern, err)
		}
		log.Info().Msgf("%s", files)
	}
}

func pull(cwd string, pattern string) {

}
