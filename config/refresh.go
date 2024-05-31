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

func (conf Config) Refresh() error {
	err := os.WriteFile(filepath.Join(conf.ConfigHome, "Taskfile.yml"), taskfile, 0644)
	if err != nil {
		return err
	}
	log.Info().Msg("run task init")

	return nil
}
