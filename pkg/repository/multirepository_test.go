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
		go mr.WriteContent(fmt.Sprintf("hello-%d", i), nil)
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

func (m *mockRepository) WriteContent(key string, content io.WriterTo) error {
	log.Info().Msgf("processing: %s by %d", key, m.id)
	time.Sleep(jobDuration)
	m.wg.Done()
	return nil
}
