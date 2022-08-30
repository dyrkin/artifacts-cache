package multipart

import (
	"artifacts-cache/pkg/file"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"testing"
)

func TestBinaryStreamIn_Save(t *testing.T) {
	binaryStreamIn := NewBinaryStreamIn(".")
	compressedStr1 := compressString("1111111111")
	compressedStr2 := compressString("222222222222222")
	reader := bytes.NewBuffer(nil)
	reader.WriteString(fmt.Sprintf("path1:%d;path2:%d\n", len(compressedStr1), len(compressedStr2)))
	reader.Write(compressedStr1)
	reader.Write(compressedStr2)
	if err := binaryStreamIn.Save(reader); err != nil {
		t.Errorf("binaryStreamIn.Save() error = %v", err)
	}

	f, _ := os.ReadFile("path1")
	if string(f) != "1111111111" {
		t.Errorf("expected %s, got %s", "1111111111", string(f))
	}
	f, _ = os.ReadFile("path2")
	if string(f) != "222222222222222" {
		t.Errorf("expected %s, got %s", "222222222222222", string(f))
	}

	file.RemoveQuiet("path1")
	file.RemoveQuiet("path2")
}

func compressString(str string) []byte {
	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)
	w.Write([]byte(str))
	w.Close()
	return buf.Bytes()
}
