package model

const (
	DefaultPattern    = "vSEMVER"
	DefaultConfigFile = ".nextver/config.yml"
)

type Config struct {
	Version string
	Pattern string
}
