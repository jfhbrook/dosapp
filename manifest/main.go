package manifest

import (
	"os"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Name           string          `yaml:"name"`
	Version        *semver.Version `yaml:"version"`
	ReleaseVersion *semver.Version `yaml:"release_version"`
}

func FromFile(path string) (*Manifest, error) {
	var file []byte
	var err error

	file, err = os.ReadFile(path)

	m := &Manifest{}

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, m)

	return m, err
}
