package compression

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

var (
	CantCompressFileError = errors.New("can't compress file")
)

type compressor struct {
	err                    error
	compressedDataReader   io.ReadCloser
	uncompressedDataWriter io.WriteCloser
}

func CompressingReader(uncompressedDataReader io.Reader) io.ReadCloser {
	compressedDataReader, compressingDataWriter := io.Pipe()

	uncompressedDataWriter, _ := gzip.NewWriterLevel(compressingDataWriter, gzip.BestSpeed)

	compressor := &compressor{err: nil, compressedDataReader: compressedDataReader, uncompressedDataWriter: uncompressedDataWriter}

	go func() {
		_, compressor.err = io.Copy(uncompressedDataWriter, uncompressedDataReader)
		if compressor.err != nil {
			compressor.err = fmt.Errorf("%w. %s", CantCompressFileError, compressor.err)
			return
		}
	}()

	return compressor
}

func (c *compressor) Read(p []byte) (n int, err error) {
	if c.err != nil {
		return 0, c.err
	}
	return c.compressedDataReader.Read(p)
}

func (c *compressor) Close() error {
	_ = c.uncompressedDataWriter.Close()
	_ = c.compressedDataReader.Close()
	return nil
}

func DecompressingReader(reader io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(reader)
}
