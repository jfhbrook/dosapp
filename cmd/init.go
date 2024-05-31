/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/jfhbrook/dosapp/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dosapp's configuration",
	Long:  `Initialize dosapp's main configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()

		editFlag, _ := cmd.Flags().GetBool("edit")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		if conf.EnvFileExists() && !overwriteFlag {
			log.Warn().Msgf("Environment file already exists at %s/dosapp.env", conf.ConfigHome)
			log.Warn().Msg("To overwrite and refresh the configuration, run 'dosapp init --overwrite'")
		} else {
			refreshFlag = true
			if err := conf.WriteEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to write env file")
			}
		}

		if editFlag {
			if err := conf.EditEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to edit config file")
			}
		}

		conf = config.LoadConfig()

		if !refreshFlag && conf.TaskFileExists() {
			log.Warn().Msgf("Taskfile already exists at %s", conf.TaskFilePath())
			log.Warn().Msg("To refresh the configuration, run 'dosapp init --refresh'")
		} else {
			refreshFlag = true
		}

		if refreshFlag {
			if err := conf.Refresh(); err != nil {
				log.Panic().Err(err).Msg("Failed to refresh config")
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
