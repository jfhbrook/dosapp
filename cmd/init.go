/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/jfhbrook/dosapp/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dosapp's configuration",
	Long:  `Initialize dosapp's main configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()

		// TODO: Does env file exist?
		// TODO: Is overwrite flag set?

		log.Info().Msg("TODO: Create env file")

		if edit, err := cmd.Flags().GetBool("edit"); err != nil && edit {
		  editor := os.Getenv("EDITOR")
		  file := filepath.Join(conf.ConfigHome, "dosapp.env")

			if err := config.EditConfig(&editor, &file); err != nil {
				log.Panic().Err(err).Msg("Failed to edit config file")
			}
		}

		conf = config.LoadConfig()

		log.Info().Msg("TODO: Reload config")

		if err := config.RefreshMain(&conf); err != nil {
			log.Panic().Err(err).Msg("Failed to reload config")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("edit", "e", true, "Edit environment files")
	initCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
