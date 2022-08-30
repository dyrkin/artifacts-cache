package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/artifacts"
	"gitlab-cache/pkg/artifacts/client"
	"gitlab-cache/pkg/artifacts/client/url"
	"gitlab-cache/pkg/file"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	log.Level(zerolog.DebugLevel)
	subset := os.Getenv("ARTIFACTS_SUBSET_ID")
	if subset == "" {
		log.Fatal().Msg("ARTIFACTS_SUBSET_ID is not set")
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
			pull(cwd, subset, pattern, repositories)
		case "push":
			push(cwd, subset, pattern, repositories)
		}
	} else {
		log.Fatal().Msgf("usage: %s <pull|push> <path|--all>", path.Base(os.Args[0]))
	}
}

func push(cwd string, subset string, pattern string, repositories []string) {
	switch pattern {
	case "--all":
	default:
		paths, err := findFiles(cwd, pattern)
		if err != nil {
			log.Fatal().Msgf("can't find paths using pattern [%s]. error: %s", pattern, err)
		}
		rotator := url.NewUrlRotator(repositories)
		limitThreads := len(repositories) * 3
		artifacts.NewPusher(client.NewPushRepositoryClient(rotator), limitThreads).Push(cwd, subset, paths)
		log.Info().Msgf("uploaded %d files", len(paths))
	}
}

func pull(cwd string, subset string, pattern string, repositories []string) {
	switch pattern {
	case "--all":
		artifacts.NewPuller(client.NewPullRepositoryClientFactory(), repositories).Pull(cwd, subset, "*")
	default:
		artifacts.NewPuller(client.NewPullRepositoryClientFactory(), repositories).Pull(cwd, subset, pattern)
	}
}

func findFiles(cwd string, pattern string) ([]string, error) {
	p := path.Join(cwd, pattern)
	if file.IsDir(p) {
		return file.FindFilesInDir(p)
	}
	return filepath.Glob(p)
}
