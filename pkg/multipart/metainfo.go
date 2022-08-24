package multipart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	InvalidMetaInfoError = errors.New("invalid meta info")
)

func MakeMetaInfo(contentDescriptors []*contentDescriptor) string {
	metaInfo := ""
	for idx, descriptor := range contentDescriptors {
		if idx < len(contentDescriptors)-1 {
			metaInfo += fmt.Sprintf("%s:%d;", descriptor.path, descriptor.size)
		} else {
			metaInfo += fmt.Sprintf("%s:%d", descriptor.path, descriptor.size)
		}
	}
	return metaInfo
}

func ParseMetaInfo(metaInfo string) ([]*contentDescriptor, error) {
	var contentDescriptors []*contentDescriptor
	for _, description := range strings.Split(metaInfo, ";") {
		parts := strings.Split(description, ":")
		size, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w. %s", InvalidMetaInfoError, err)
		}
		contentDescriptors = append(contentDescriptors, &contentDescriptor{path: parts[0], size: size})
	}
	return contentDescriptors, nil
}
