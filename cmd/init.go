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
		conf := config.NewConfig()

		editFlag, _ := cmd.Flags().GetBool("edit")
		editFlagChanged := cmd.Flags().Changed("edit")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		if !overwriteFlag && conf.EnvFileExists() {
			log.Warn().Msgf("Environment file already exists at %s", conf.EnvFilePath())
			log.Warn().Msg("To overwrite and refresh the configuration, run 'dosapp init --overwrite'")
		} else {
			refreshFlag = true
			if err := conf.WriteEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to write env file")
			}
		}

		// It's technically possible for the env file to configure a different
		// location for the TaskFile. But that's based on the same config home
		// as with the env file. In other words, this is a sensible assumption,
		// and setting config home this way is unsupported.
		if !refreshFlag && conf.TaskFileExists() {
			log.Warn().Msgf("Taskfile already exists at %s", conf.TaskFilePath())
			log.Warn().Msg("To refresh the configuration, run 'dosapp init --refresh'")
		} else {
			refreshFlag = true
			// Will write task file on refresh step
		}

		// If we're not refreshing, then we only want to edit the file if
		// explicitly asked
		var shouldEdit bool

		if refreshFlag {
			shouldEdit = editFlag
		} else {
			shouldEdit = editFlag && editFlagChanged
		}

		if shouldEdit {
			if err := conf.EditEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to edit config file")
			}
		}

		// Pick up any changes made from editing
		conf = config.NewConfig()

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
