package provider

import (
	"github.com/tauffredou/nextver/model"
)

type Provider interface {
	GetLatestRelease() *model.Release
	GetRelease(name string) *model.Release
}
