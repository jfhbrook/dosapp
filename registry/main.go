package registry

import (
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/config"
)

type Registry interface {
	Stage() *Stage
	FindPackage(name string) *Package
}

func NewRegistry(conf *config.Config) Registry {
	u, err := url.Parse(conf.Registry)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse registry URL.")
	}

	if u.Scheme != "github" {
		log.Fatal().Str("forge", u.Scheme).Msgf("Forge is unsupported.")
	}

	return newGitHubRegistry(conf, u)
}
