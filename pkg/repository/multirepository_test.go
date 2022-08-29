package repository

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
	"testing"
	"time"
)

var jobDuration = time.Second

func TestMultiRepository_WriteContent(t *testing.T) {
	jobsCount := 30
	limitThreads := 10
	maxWorkDuration := time.Duration(int64(jobsCount/limitThreads+1) * int64(jobDuration))

	factory := &mockRepositoryFactory{wg: &sync.WaitGroup{}}
	factory.wg.Add(jobsCount)
	mr := NewMultiRepository(factory, limitThreads)
	start := time.Now()
	for i := 0; i < jobsCount; i++ {
		go mr.WriteContent("subset", fmt.Sprintf("name-%d", i), nil)
	}
	factory.wg.Wait()
	end := time.Now()
	workDuration := end.Sub(start)

	if workDuration > maxWorkDuration {
		t.Fatalf("work duration too long: %d. expected: %d", workDuration, maxWorkDuration)
	}
	log.Info().Msgf("took %s", workDuration)
}

type mockRepositoryFactory struct {
	wg *sync.WaitGroup
	id int
}

func (m *mockRepositoryFactory) Create() Repository {
	m.id++
	return &mockRepository{m.id, m.wg}
}

type mockRepository struct {
	id int
	wg *sync.WaitGroup
}

func (m *mockRepository) FindContent(subset, filter string) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepository) WriteContent(subset string, path string, content io.Reader) error {
	log.Info().Msgf("processing: %s from subset '%s' by %d", path, subset, m.id)
	time.Sleep(jobDuration)
	m.wg.Done()
	return nil
}
