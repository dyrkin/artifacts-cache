package url

import (
	"artifacts-cache/pkg/slice"
	"sync"
)

type Rotator interface {
	Next() string
}

type urlRotator struct {
	urls  []string
	mutex *sync.Mutex
}

func NewUrlRotator(urls []string) *urlRotator {
	shuffled := make([]string, len(urls))
	copy(shuffled, urls)
	slice.Shuffle(shuffled)
	return &urlRotator{shuffled, &sync.Mutex{}}
}

func (u *urlRotator) Next() string {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	url := u.urls[0]
	u.urls = append(u.urls[1:], url)
	return url
}
