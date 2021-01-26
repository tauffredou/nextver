package sorter_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/tauffredou/nextver/model"
	"github.com/tauffredou/nextver/sorter"
	"sort"
	"testing"
)

func TestSortBySemver(t *testing.T) {

	releases := []model.Release{
		model.Release{CurrentVersion: "v1.4.0"},
		model.Release{CurrentVersion: "v1.2.0"},
		model.Release{CurrentVersion: "v1.11.0"},
		model.Release{CurrentVersion: "v1.7.0"},
	}

	expected := []model.Release{
		model.Release{CurrentVersion: "v1.11.0"},
		model.Release{CurrentVersion: "v1.7.0"},
		model.Release{CurrentVersion: "v1.4.0"},
		model.Release{CurrentVersion: "v1.2.0"},
	}

	sort.Sort(sorter.BySemver(releases))

	assert.Equal(t, expected, releases)

}

func TestSortByDate(t *testing.T) {

	releases := []model.Release{
		model.Release{CurrentVersion: "2006-01-02-150405"},
		model.Release{CurrentVersion: "2008-01-02-150405"},
		model.Release{CurrentVersion: "2002-01-02-150405"},
		model.Release{CurrentVersion: "2001-01-02-150405"},
	}

	expected := []model.Release{
		model.Release{CurrentVersion: "2008-01-02-150405"},
		model.Release{CurrentVersion: "2006-01-02-150405"},
		model.Release{CurrentVersion: "2002-01-02-150405"},
		model.Release{CurrentVersion: "2001-01-02-150405"},
	}

	sort.Sort(sorter.BySemver(releases))

	assert.Equal(t, expected, releases)

}
