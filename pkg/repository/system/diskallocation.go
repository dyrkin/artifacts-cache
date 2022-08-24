package system

import (
	"artifacts-cache/pkg/repository/basedir"
	"golang.org/x/sys/unix"
	"time"
)

type DiskAllocation interface {
	Start()
	GetFreeSpaceInPercents() float64
}

type diskAllocation struct {
	baseDir basedir.BaseDir
	stats   *unix.Statfs_t
}

func NewDiskAllocation(baseDir basedir.BaseDir) *diskAllocation {
	return &diskAllocation{
		baseDir: baseDir,
		stats:   &unix.Statfs_t{},
	}
}

func (d *diskAllocation) Start() {
	_ = unix.Statfs(d.baseDir.Dir(), d.stats)
	go func() {
		for range time.Tick(time.Microsecond * 200) {
			_ = unix.Statfs(d.baseDir.Dir(), d.stats)
		}
	}()
}

func (d *diskAllocation) GetFreeSpaceInPercents() float64 {
	return 100 / float64(d.stats.Blocks) * float64(d.stats.Bfree)
}
