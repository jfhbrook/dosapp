package registry

import (
	"fmt"
	"net/url"

	"github.com/jfhbrook/dosapp/config"
)

type RegistryError struct {
	message string
}

func (e *RegistryError) Error() string {
	return e.message
}

type Registry interface {
	FindPackage(name string) *Package
}

func NewRegistry(conf *config.Config) (Registry, error) {
	u, err := url.Parse(conf.Registry)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "github" {
		return nil, &RegistryError{fmt.Sprintf("Forge %s is unsupported.", u.Scheme)}
	}

	return newGitHubRegistry(conf, u)
}
