package multipart

import "fmt"

func makeMetaInfo(contentDescriptors []*contentDescriptor) string {
	metaInfo := ""
	for idx, descriptor := range contentDescriptors {
		if idx < len(contentDescriptors)-1 {
			metaInfo += fmt.Sprintf("%s:%d;", descriptor.path, descriptor.size)
		} else {
			metaInfo += fmt.Sprintf("%s:%d", descriptor.path, descriptor.size)
		}
	}
	metaInfo += "\n"
	return metaInfo
}
