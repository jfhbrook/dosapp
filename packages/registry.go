package packages

import (
	"github.com/jfhbrook/dosapp/config"
)

// TODO: what is the interface for a registry?
type Registry struct {
	Name   string
	Config *config.Config
}

func NewRegistry(conf *config.Config, name string) *Registry {
	return &Registry{
		Name:   name,
		Config: conf,
	}
}
