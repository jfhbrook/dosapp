package registry

import (
	"os"
	"path/filepath"

	"github.com/jfhbrook/dosapp/config"
)

type Cache struct {
	Config *config.Config
}

func newCache(conf *config.Config) *Cache {
	return &Cache{
		Config: conf,
	}
}

func (ch *Cache) CachedPackagePath(name string) string {
	return filepath.Join(ch.Config.PackageCacheHome, name)
}

func (ch *Cache) CachedPackageExists(name string) bool {
	_, err := os.Stat(ch.CachedPackagePath(name))
	return err == nil
}

func (ch *Cache) RemoveCachedPackage(name string) error {
	return os.RemoveAll(ch.CachedPackagePath(name))
}
