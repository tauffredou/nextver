package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGithubProvider_emptyConfig(t *testing.T) {
	_, err := NewGithubProvider("owner", "repo", "token", nil)
	if assert.Error(t, err) {
		assert.Equal(t, &ConfigurationError{}, err)
	}

}
