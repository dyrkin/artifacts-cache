package compression

import (
	"compress/gzip"
	"errors"
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

func CompressingReader(uncompressedDataReader io.ReadCloser) io.ReadCloser {
	r, w := io.Pipe()
	go func() {

		gzw := gzip.NewWriter(w)
		_, err := io.Copy(gzw, uncompressedDataReader)

		if err != nil {
			w.CloseWithError(err)
			return
		}

		w.CloseWithError(gzw.Close())
	}()
	return r
}

func DecompressingReader(reader io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(reader)
}
