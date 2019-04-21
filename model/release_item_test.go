package model

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var (
	testDate = MustParse(time.RFC3339, "2010-01-02T12:34:00Z")
)

func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func TestNewReleaseItem_withScope(t *testing.T) {
	ri := NewReleaseItem("tauf", testDate, "feat(scope): pouet")

	expected := ReleaseItem{
		Kind:   "feat",
		Scope:  "scope",
		Detail: "",
		Title:  "pouet",
		Level:  MINOR,
		Author: "tauf",
		Date:   testDate,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withSpecialChars(t *testing.T) {
	ri := NewReleaseItem("tauf", testDate, "some-kind(scope): pouet")

	expected := ReleaseItem{
		Kind:   "some-kind",
		Scope:  "scope",
		Detail: "",
		Title:  "pouet",
		Level:  UNDEFINED,
		Author: "tauf",
		Date:   testDate,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withoutScope(t *testing.T) {
	ri := NewReleaseItem("tauf", testDate, "feat: pouet")

	expected := ReleaseItem{
		Kind:   "feat",
		Scope:  "",
		Detail: "",
		Title:  "pouet",
		Level:  MINOR,
		Author: "tauf",
		Date:   testDate,
	}

	assert.Equal(t, ri, expected)
}

func TestNewReleaseItem_withText(t *testing.T) {
	message := `feat(feature-1234): commit message

This do that
`
	ri := NewReleaseItem("tauf", testDate, message)

	expected := ReleaseItem{
		Kind:   "feat",
		Title:  "commit message",
		Detail: "This do that",
		Level:  MINOR,
		Scope:  "feature-1234",
		Date:   testDate,
		Author: "tauf",
	}
	assert.Equal(t, ri, expected)

}
