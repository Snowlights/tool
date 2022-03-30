package parse

import (
	"errors"
	"strings"
)

const (
	JsonTag           = "json"
	PropertiesTagName = "properties"

	dot                 = "."
	comma               = ","
	singleHorizontalBar = "-"
)

var (
	InvalidUnmarshalError = errors.New("v must be a non-nil struct pointer")
	UnsupportedTypeError  = errors.New("unsupported type")
)

func parseTag(tag string) string {
	if idx := strings.Index(tag, comma); idx != -1 {
		return tag[:idx]
	}
	return tag
}
