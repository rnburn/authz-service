package auth

import (
	"strings"
)

// contentTypeAllowList is the list of allowed content types in lowercase
var contentTypeAllowListLowerCase = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

func shouldRecordBody(headers map[string]string) bool {
	contentType := strings.ToLower(headers["content-type"])
	for _, recordableContentType := range contentTypeAllowListLowerCase {
		if strings.Contains(contentType, recordableContentType) {
			return true
		}
	}
	return false
}
