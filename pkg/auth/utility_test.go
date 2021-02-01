package auth

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestShouldRecordBody(t *testing.T) {
  headers := make(map[string]string)
  headers["content-type"] = "application/json"
  assert.True(t, shouldRecordBody(headers))
  headers["content-type"] = "aPPlication/json"
  assert.True(t, shouldRecordBody(headers))
  headers["content-type"] = "text"
  assert.False(t, shouldRecordBody(headers))
}
