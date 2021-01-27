package sorter_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tauffredou/nextver/model"
	"github.com/tauffredou/nextver/sorter"
)

func TestSortBySemver(t *testing.T) {

	releases := []model.Release{
		{CurrentVersion: "v0.1.0"},
		{CurrentVersion: "v0.10.0"},
		{CurrentVersion: "v0.3.0"},
		{CurrentVersion: "v0.5.0"},
		{CurrentVersion: "v0.5.1"},
		{CurrentVersion: "v0.6.0"},
	}
	expected := []model.Release{
		{CurrentVersion: "v0.10.0"},
		{CurrentVersion: "v0.6.0"},
		{CurrentVersion: "v0.5.1"},
		{CurrentVersion: "v0.5.0"},
		{CurrentVersion: "v0.3.0"},
		{CurrentVersion: "v0.1.0"},
	}

	sort.Sort(sorter.BySemver(releases))

	assert.Equal(t, expected, releases)

}

func TestSortByDate(t *testing.T) {

	releases := []model.Release{
		{CurrentVersion: "2006-01-02-150405"},
		{CurrentVersion: "2008-01-02-150405"},
		{CurrentVersion: "2002-01-02-150405"},
		{CurrentVersion: "2001-01-02-150405"},
	}

	expected := []model.Release{
		{CurrentVersion: "2008-01-02-150405"},
		{CurrentVersion: "2006-01-02-150405"},
		{CurrentVersion: "2002-01-02-150405"},
		{CurrentVersion: "2001-01-02-150405"},
	}

	sort.Sort(sorter.BySemver(releases))

	assert.Equal(t, expected, releases)

}
