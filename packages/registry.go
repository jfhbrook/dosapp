package packages

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/config"
)

type RegistryError struct {
	message string
}

func (e *RegistryError) Error() string {
	return e.message
}

type Registry interface {
	PackageUrl(name string) (string, error)
}

type GitHubRegistry struct {
	Config     *config.Config
	user       string
	repository string
}

func newGitHubRegistry(conf *config.Config, u *url.URL) (Registry, error) {
	user := u.Host
	re := regexp.MustCompile(`^/`)
	repo := re.ReplaceAllString(u.Path, "")

	return &GitHubRegistry{
		Config:     conf,
		user:       user,
		repository: repo,
	}, nil
}

func (registry *GitHubRegistry) PackageUrl(name string) (string, error) {
	log.Info().Msg("TODO: pull github releases")
	log.Info().Msg("TODO: filter by app name")
	log.Info().Msg("TODO: choose latest release (by sort order)")
	log.Info().Msg("TODO: snag artifact url")
	return "", nil
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
