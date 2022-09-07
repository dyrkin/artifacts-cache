package compression

import (
	"compress/gzip"
	"io"
)

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
