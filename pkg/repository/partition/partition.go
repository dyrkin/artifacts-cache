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
	CantCheckPartitionExistenceError = errors.New("can't check partition existence")
	CantCreatePartitionFileError     = errors.New("can't create partition")
	CantCreateTransactionError       = errors.New("can't create transaction")
	CantOpenPartitionFileError       = errors.New("partition exists in index, but partition file can not be found")
	CantAddPartitionToIndexError     = errors.New("can't add partition to index")
	CantAddContentKeyToIndexError    = errors.New("can't add content key to index")
	CantWriteContentToPartitionError = errors.New("can't write content to partition")
	ReadOnlyPartitionError           = errors.New("partition is read only")
)

const (
	MaxPartitionSize = 500 * 1024 * 1024 //500MB
)

type ContentWriter interface {
	WriteContent(subset string, name string, content io.Reader) error
}

type ContentOpener interface {
	OpenContent(offset int64, size int64) io.ReadCloser
}

type Partition interface {
	ContentWriter
	Open() error
	Close() error
	IsClosed() bool
}

type partition struct {
	offset      int64
	uuid        string
	baseDir     basedir.BaseDir
	index       index.Index
	file        *os.File
	partitionId int64
	closed      bool
	readOnly    bool
}

type readCloser struct {
	closer io.Closer
	reader io.Reader
}

func NewPartition(uuid string, baseDir basedir.BaseDir, index index.Index) *partition {
	return &partition{
		offset:  0,
		uuid:    uuid,
		baseDir: baseDir,
		index:   index,
	}
}

func NewReadOnlyPartition(uuid string, baseDir basedir.BaseDir) *partition {
	return &partition{
		offset:   0,
		uuid:     uuid,
		baseDir:  baseDir,
		readOnly: true,
	}
}

func (p *partition) Open() error {
	dir, err := p.baseDir.MakeSubdirByUUID(p.uuid)
	if err != nil {
		return fmt.Errorf("%w. %s", CantCreatePartitionDirError, err)
	}
	location := path.Join(dir, p.uuid)
	if !p.readOnly {
		existingPartitionId, ok, err := p.index.PartitionExists(location)
		if err != nil {
			return fmt.Errorf("%w. %s", CantCheckPartitionExistenceError, err)
		}
		if ok {
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
			if err != nil {
				file.CloseQuiet(partitionFile)
				file.RemoveQuiet(location)
				return fmt.Errorf("%w. %s", CantCreateTransactionError, err)
			}
			partitionId, err := p.index.AddPartition(p.uuid)
			if err != nil {
				file.CloseQuiet(partitionFile)
				file.RemoveQuiet(location)
				return fmt.Errorf("%w. %s", CantAddPartitionToIndexError, err)
			}
			p.file = partitionFile
			p.partitionId = partitionId
		}
	} else {
		partitionFile, err := file.Open(location)
		if err != nil {
			return fmt.Errorf("%w. %s", CantOpenPartitionFileError, err)
		}
		p.file = partitionFile
	}
	return nil
}

func (p *partition) WriteContent(subset string, name string, content io.Reader) error {
	if !p.readOnly {
		offset := p.offset
		size, err := io.Copy(p.file, content)
		p.offset += size
		if err != nil {
			return fmt.Errorf("%w. %s", CantWriteContentToPartitionError, err)
		}
		err = p.index.AddFileToPartition(p.partitionId, subset, name, offset, size)
		if err != nil {
			return fmt.Errorf("%w. %s", CantAddContentKeyToIndexError, err)
		}
		if p.offset > MaxPartitionSize {
			return p.Close()
		}
	} else {
		return ReadOnlyPartitionError
	}
	return nil
}

func (p *partition) OpenContent(offset int64, size int64) io.ReadCloser {
	return &readCloser{p, io.NewSectionReader(p.file, offset, size)}
}

func (p *partition) Close() (err error) {
	if p.file != nil {
		file.CloseQuiet(p.file)
	}
	p.closed = true
	return err
}

func (p *partition) IsClosed() bool {
	return p.closed
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	return rc.reader.Read(p)
}

func (rc *readCloser) Close() error {
	return rc.closer.Close()
}
