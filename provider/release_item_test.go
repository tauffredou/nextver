package provider

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewReleaseItem_withScope(t *testing.T) {
	ri := NewReleaseItem("feat(scope): pouet")

	expected := ReleaseItem{
		Kind:   "feat",
		Scope:  "scope",
		Detail: "",
		Title:  "pouet",
		Level:  MINOR,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withSpecialChars(t *testing.T) {
	ri := NewReleaseItem("some-kind(scope): pouet")

	expected := ReleaseItem{
		Kind:   "some-kind",
		Scope:  "scope",
		Detail: "",
		Title:  "pouet",
		Level:  UNDEFINED,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withoutScope(t *testing.T) {
	ri := NewReleaseItem("feat: pouet")

	expected := ReleaseItem{
		Kind:   "feat",
		Scope:  "",
		Detail: "",
		Title:  "pouet",
		Level:  MINOR,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withText(t *testing.T) {
	message := `feat(feature-1234): commit message

This do that
`
	ri := NewReleaseItem(message)

	expected := ReleaseItem{
		Kind:   "feat",
		Title:  "commit message",
		Detail: "This do that",
		Level:  MINOR,
		Scope:  "feature-1234",
	}
	assert.Equal(t, ri, expected)

}
