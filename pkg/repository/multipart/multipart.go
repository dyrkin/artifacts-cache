package multipart

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab-cache/pkg/repository/basedir"
	"gitlab-cache/pkg/repository/index"
	"gitlab-cache/pkg/repository/partition"
	"io"
)

var (
	CantOpenPartitionError = errors.New("can't open partition")
)

type BinaryStreamFactory struct {
	baseDir basedir.BaseDir
}

type binaryStream struct {
	contentDescriptors           []*contentDescriptor
	contentDescriptorsForClosing []*contentDescriptor
	metaPushed                   bool
	meta                         io.Reader
}

type contentDescriptor struct {
	content io.ReadCloser
	path    string
	size    int64
}

func NewBinaryStreamFactory(baseDir basedir.BaseDir) *BinaryStreamFactory {
	return &BinaryStreamFactory{
		baseDir: baseDir,
	}
}

func (c *BinaryStreamFactory) Create(contentEmplacement *index.ContentEmplacement) (io.ReadCloser, error) {
	var contentDescriptors []*contentDescriptor
	for _, emplacement := range contentEmplacement.Emplacements {
		p := partition.NewReadOnlyPartition(emplacement.Partition, c.baseDir)
		err := p.Open()
		if err != nil {
			return nil, fmt.Errorf("%w. %s", CantOpenPartitionError, err)
		}
		content := p.OpenContent(emplacement.Offset, emplacement.Size)
		descriptor := &contentDescriptor{content: content, path: emplacement.Path, size: emplacement.Size}
		contentDescriptors = append(contentDescriptors, descriptor)
	}
	metaInfo := makeMetaInfo(contentDescriptors)
	return &binaryStream{
		contentDescriptors:           contentDescriptors,
		contentDescriptorsForClosing: contentDescriptors,
		meta:                         bytes.NewBufferString(metaInfo),
	}, nil
}

func (b *binaryStream) Read(p []byte) (n int, err error) {
	n, err = b.meta.Read(p)
	if err != io.EOF {
		return
	}
	for i, contentDescriptor := range b.contentDescriptors {
		n, err = contentDescriptor.content.Read(p)
		if err == io.EOF && i < len(b.contentDescriptors)-1 {
			err = nil
			b.contentDescriptors = b.contentDescriptors[i+1:]
			continue
		}
		if err != nil {
			return
		}
		return
	}
	return 0, nil
}

func (b *binaryStream) Close() error {
	for _, descriptor := range b.contentDescriptorsForClosing {
		descriptor.content.Close()
	}
	return nil
}
