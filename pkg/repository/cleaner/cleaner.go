package cleaner

import (
	"artifacts-cache/pkg/file"
	"artifacts-cache/pkg/repository/basedir"
	"artifacts-cache/pkg/repository/index"
	"artifacts-cache/pkg/repository/system"
	"github.com/rs/zerolog/log"
	"path"
	"sync"
	"time"
)

type Cleaner struct {
	baseDir                   basedir.BaseDir
	index                     index.Index
	diskAllocation            system.DiskAllocation
	targetFreeSpaceInPercents int
	mutex                     *sync.Mutex
}

func NewCleaner(baseDir basedir.BaseDir, index index.Index, diskAllocation system.DiskAllocation, targetFreeSpaceInPercents int) *Cleaner {
	return &Cleaner{
		baseDir:                   baseDir,
		index:                     index,
		diskAllocation:            diskAllocation,
		targetFreeSpaceInPercents: targetFreeSpaceInPercents,
		mutex:                     &sync.Mutex{},
	}
}

func (c *Cleaner) Schedule() {
	go func() {
		for range time.Tick(time.Minute * 5) {
			c.clean()
		}
	}()
	c.clean()
}

func (c *Cleaner) clean() {
	availableSpaceInPercents := c.diskAllocation.GetFreeSpaceInPercents()
	if availableSpaceInPercents > float64(c.targetFreeSpaceInPercents) {
		log.Info().Msgf("cleanup is not required. available space in percents [%f%%]. expected [%d%%]", availableSpaceInPercents, c.targetFreeSpaceInPercents)
		return
	}
	log.Info().Msgf("available space in percents [%f%%]. expected [%d%%]. starting cleanup", availableSpaceInPercents, c.targetFreeSpaceInPercents)
	uuids, err := c.index.GetEldestUuids(100)
	if err != nil {
		log.Error().Msgf("can't get eldest uuids. %s", err)
		return
	}

	for _, uuid := range uuids {
		partition := path.Join(c.baseDir.GetSubdirByUUID(uuid), uuid)
		err := file.Remove(partition)
		if err != nil {
			log.Error().Msgf("can't delete partition. %s", err)
		}
		err = c.index.DeletePartition(uuid)
		if err != nil {
			log.Error().Msgf("can't delete artifact from index. %s", err)
			continue
		}
		if c.diskAllocation.GetFreeSpaceInPercents() > float64(c.targetFreeSpaceInPercents) {
			break
		}
	}
	if c.diskAllocation.GetFreeSpaceInPercents() < float64(c.targetFreeSpaceInPercents) {
		c.clean()
	}
}
