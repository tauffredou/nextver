package provider

import (
	"github.com/tauffredou/nextver/model"
)

type Provider interface {
	GetLatestRelease() model.Release
}
