package multipart

import (
	"artifacts-cache/pkg/compression"
	"artifacts-cache/pkg/file"
	reader2 "artifacts-cache/pkg/reader"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"path"
)

var (
	CantReadMetaInfoError       = errors.New("can't read meta info")
	CantCreateDecompressorError = errors.New("can't create decompressor")
	CantCreateFileError         = errors.New("can't create file")
	CantSaveDataToFileError     = errors.New("can't save data to file")
	CantCreateDirForFileError   = errors.New("can't create directory for file")
)

type binaryStreamIn struct {
	cwd string
}

func NewBinaryStreamIn(cwd string) *binaryStreamIn {
	return &binaryStreamIn{cwd: cwd}
}

func (b *binaryStreamIn) Save(reader io.Reader) error {
	metaInfo, err := reader2.NewStringReader(reader).Readline()
	if err != nil {
		return fmt.Errorf("%w. %s", CantReadMetaInfoError, err)
	}
	contentDescriptors, err := ParseMetaInfo(metaInfo)
	if err != nil {
		return err
	}
	for _, descriptor := range contentDescriptors {
		fileReader := io.LimitReader(reader, descriptor.size)
		uncompressedReader, err := compression.DecompressingReader(fileReader)
		if err != nil {
			return fmt.Errorf("%w. %s", CantCreateDecompressorError, err)
		}
		err = file.MkdirAllForFile(path.Join(b.cwd, descriptor.path))
		if err != nil {
			return fmt.Errorf("%w. %s", CantCreateDirForFileError, err)
		}
		f, err := file.CreateEmpty(path.Join(b.cwd, descriptor.path))
		if err != nil {
			return fmt.Errorf("%w. %s", CantCreateFileError, err)
		}
		_, err = io.Copy(f, uncompressedReader)
		if err != nil {
			return fmt.Errorf("%w. %s", CantSaveDataToFileError, err)
		}
	}
	log.Printf("OK. downloaded %d files", len(contentDescriptors))
	return nil
}
