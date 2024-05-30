/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package config

import (
	_ "embed"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

//go:embed Taskfile.yml
var taskfile []byte

// NOTE: flags are checked in cmd
func RefreshMain(config *Config) error {
	err := os.WriteFile(filepath.Join(config.ConfigHome, "Taskfile.yml"), taskfile, 0644)
	if err != nil {
		return err
	}
	log.Info().Msg("run task init")

	return nil
}
