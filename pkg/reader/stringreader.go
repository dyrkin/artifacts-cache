package reader

import "io"

type StringReader struct {
	reader io.Reader
}

func NewStringReader(reader io.Reader) *StringReader {
	return &StringReader{reader: reader}
}

func (r *StringReader) Readline() (string, error) {
	line := make([]byte, 0, 100)
	for {
		buf := make([]byte, 1)
		n, err := r.reader.Read(buf)
		if n > 0 {
			c := buf[0]
			if c == '\n' {
				break
			}
			line = append(line, c)
		}
		if err != nil {
			return "", err
		}
	}
	return string(line), nil
}
