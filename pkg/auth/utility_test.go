package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
