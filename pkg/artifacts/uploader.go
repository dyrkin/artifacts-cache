package artifacts

import (
	"math/rand"
	"sync"
	"time"
)

const (
	treadsPerStorage = 3
)

type uploader struct {
	storages []string
	treads   int
}

func NewUploader(storages []string) *uploader {
	shuffleStorages(storages)
	return &uploader{storages, len(storages) * treadsPerStorage}
}

func shuffleStorages(storages []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(storages), func(i, j int) {
		storages[i], storages[j] = storages[j], storages[i]
	})
}

func (u *uploader) Upload(cwd string, files []string) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(files))
}

func (u *uploader) upload(cwd string, file string) error {
	return nil
}
