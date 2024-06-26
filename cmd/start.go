/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	application "github.com/jfhbrook/dosapp/app"
	"github.com/jfhbrook/dosapp/config"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start an application",
	Long: `Start a DOS application.

The --refresh flag will generate a fresh Taskfile.yml and DOXBox .conf files
based on the application configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		conf := config.NewConfig()

		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		app := application.NewApp(appName, conf)

		if !app.Exists() {
			log.Fatal().Msgf("%s not found. Did you install it?", appName)
		}

		if refreshFlag {
			if err := app.Refresh(); err != nil {
				log.Panic().Err(err).Msg("Failed to refresh application config")
			}
		}

		app.Run("start")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
