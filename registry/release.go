package registry

import (
	"github.com/Masterminds/semver/v3"
)

// TODO: This struct is specific to github. Perhaps the github registry
// should be in its own file?
type Release struct {
	Name           string
	Version        *semver.Version
	ReleaseVersion *semver.Version
	Url            string
}

func NewRelease(name string, version string, releaseVersion string, url string) (*Release, error) {
	ver, err := semver.NewVersion(version)

	if err != nil {
		return nil, err
	}

	var relVer *semver.Version
	relVer, err = semver.NewVersion(releaseVersion)

	if err != nil {
		return nil, err
	}

	return &Release{
		Name:           name,
		Version:        ver,
		ReleaseVersion: relVer,
		Url:            url,
	}, nil
}

// TODO: This is the url for the HTML list of assets, not the actual asset
// we want. That's going to be on Assets[0].Url.
// See: https://github.com/google/go-github/blob/446a2968d80b0adc88c54a1c1a17a05b25e83ea7/github/repos_releases.go#L68
func (release *Release) URL() string {
	return release.Url
}
