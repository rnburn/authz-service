package auth

import (
  "strings"
)

// contentTypeAllowList is the list of allowed content types in lowercase
var contentTypeAllowListLowerCase = []string{
	"application/json",
	"application/x-www-form-urlencoded",
}

func shouldRecordBody(content_type string) bool {
  content_type = strings.ToLower(content_type)
  for _, recordableContentType := range contentTypeAllowListLowerCase {
    if strings.Contains(content_type, recordableContentType) {
      return true
    }
  }
  return false
}
