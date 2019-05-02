package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadHubToken(t *testing.T) {
	r, err := readHubToken("fixtures/config_hub.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "xxxxxx", r)
}

func TestReadHubToken_FileNotExits(t *testing.T) {
	_, err := readHubToken("dummy.yaml")
	assert.Error(t, err, "")
}
