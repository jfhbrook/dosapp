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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dosapp's configuration",
	Long:  `Initialize dosapp's main configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()
		envFile := filepath.Join(conf.ConfigHome, "dosapp.env")

		editFlag, _ := cmd.Flags().GetBool("edit")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		if exists(envFile) && !overwriteFlag {
			log.Warn().Msgf("Environment file already exists at %s/dosapp.env", conf.ConfigHome)
			log.Warn().Msg("To overwrite and refresh the configuration, run 'dosapp init --overwrite'")
		} else {
			refreshFlag = true
			log.Info().Msg("TODO: Create env file")
		}

		if editFlag {
			editor := os.Getenv("EDITOR")

			if err := config.EditConfig(&editor, &envFile); err != nil {
				log.Panic().Err(err).Msg("Failed to edit config file")
			}
		}

		conf = config.LoadConfig()

		if refreshFlag {
			if err := config.RefreshMain(&conf); err != nil {
				log.Panic().Err(err).Msg("Failed to reload config")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("edit", "e", true, "Edit environment files")
	initCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	initCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
