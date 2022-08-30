package multipart

import (
	"artifacts-cache/pkg/repository/basedir"
	"artifacts-cache/pkg/repository/index"
	"artifacts-cache/pkg/repository/partition"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	CantOpenPartitionError = errors.New("can't open partition")
)

type BinaryStreamOutFactory interface {
	Create(contentEmplacement *index.ContentEmplacement) (io.ReadCloser, error)
}

type binaryStreamOutFactory struct {
	baseDir basedir.BaseDir
}

type binaryStreamOut struct {
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

func NewBinaryStreamOutFactory(baseDir basedir.BaseDir) *binaryStreamOutFactory {
	return &binaryStreamOutFactory{
		baseDir: baseDir,
	}
}

func (c *binaryStreamOutFactory) Create(contentEmplacement *index.ContentEmplacement) (io.ReadCloser, error) {
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
	metaInfo := MakeMetaInfo(contentDescriptors)
	return &binaryStreamOut{
		contentDescriptors:           contentDescriptors,
		contentDescriptorsForClosing: contentDescriptors,
		meta:                         bytes.NewBufferString(metaInfo + "\n"),
	}, nil
}

func (b *binaryStreamOut) Read(p []byte) (n int, err error) {
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

func (b *binaryStreamOut) Close() error {
	for _, descriptor := range b.contentDescriptorsForClosing {
		descriptor.content.Close()
	}
	return nil
}
