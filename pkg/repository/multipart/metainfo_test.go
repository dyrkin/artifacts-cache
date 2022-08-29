package multipart

import "testing"

func Test_makeMetaInfo(t *testing.T) {
	contentDescriptors := []*contentDescriptor{
		{
			content: nil,
			path:    "path/file1",
			size:    1,
		},
		{
			content: nil,
			path:    "path/file2",
			size:    2,
		},
	}

	metaInfo := makeMetaInfo(contentDescriptors)
	if metaInfo != "path/file1:1;path/file2:2\n" {
		t.Errorf("Expected metaInfo to be 'path/file1:1;path/file2:2\n', but got '%s'", metaInfo)
	}
}
