package partition

import (
	"errors"
	"fmt"
	"gitlab-cache/pkg/file"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"io"
	"os"
	"path"
)

var (
	CantCreatePartitionDirError      = errors.New("can't create directory for partition")
	CantCreatePartitionFileError     = errors.New("can't create partition")
	CantCreateTransactionError       = errors.New("can't create transaction")
	CantOpenPartitionFileError       = errors.New("partition exists in index, but we can't open partition file")
	CantAddPartitionToIndexError     = errors.New("can't add partition to index")
	CantAddContentKeyToIndexError    = errors.New("can't add content key to index")
	CantWriteContentToPartitionError = errors.New("can't write content to partition")
)

const (
	MaxPartitionSize = 100 * 1024 * 1024 //100MB
)

type Partition interface {
	Open() error
	WriteContent(key string, content io.WriterTo) error
	Close() error
	IsClosed() bool
}

type partition struct {
	offset      int64
	uuid        string
	baseDir     basedir.BaseDir
	index       index.Index
	transaction index.Transaction
	file        *os.File
	partitionId int64
	closed      bool
}

func NewPartition(uuid string, baseDir basedir.BaseDir, index index.Index) *partition {
	return &partition{
		offset:  0,
		uuid:    uuid,
		baseDir: baseDir,
		index:   index,
	}
}

func (p *partition) Open() error {
	dir, err := p.baseDir.MakeSubdirByUUID(p.uuid)
	if err != nil {
		return fmt.Errorf("%w. %s", CantCreatePartitionDirError, err)
	}
	location := path.Join(dir, p.uuid)
	if existingPartitionId, ok := p.index.PartitionExists(location); ok {
		partitionFile, err := file.Open(location)
		if err != nil {
			return fmt.Errorf("%w. %s", CantOpenPartitionFileError, err)
		}
		p.file = partitionFile
		p.partitionId = existingPartitionId
	} else {
		partitionFile, err := file.CreateEmpty(location)
		if err != nil {
			return fmt.Errorf("%w. %s", CantCreatePartitionFileError, err)
		}
		p.transaction, err = p.index.CreateTransaction()
		if err != nil {
			file.CloseQuiet(partitionFile)
			file.RemoveQuiet(location)
			return fmt.Errorf("%w. %s", CantCreateTransactionError, err)
		}
		partitionId, err := p.transaction.AddPartition(p.uuid)
		if err != nil {
			file.CloseQuiet(partitionFile)
			file.RemoveQuiet(location)
			return fmt.Errorf("%w. %s", CantAddPartitionToIndexError, err)
		}
		p.file = partitionFile
		p.partitionId = partitionId
	}
	return nil
}

func (p *partition) WriteContent(key string, content io.WriterTo) error {
	begin := p.offset
	size, err := content.WriteTo(p.file)
	p.offset += size
	if err != nil {
		return fmt.Errorf("%w. %s", CantWriteContentToPartitionError, err)
	}
	err = p.transaction.AddFileToPartition(key, p.partitionId, begin, size)
	if err != nil {
		return fmt.Errorf("%w. %s", CantAddContentKeyToIndexError, err)
	}
	if p.offset > MaxPartitionSize {
		return p.Close()
	}
	return nil
}

func (p *partition) Close() (err error) {
	if p.file != nil {
		file.CloseQuiet(p.file)
	}
	if p.transaction != nil {
		err = p.transaction.Commit()
	}
	p.closed = true
	return err
}

func (p *partition) IsClosed() bool {
	return p.closed
}
