/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/application"
	"github.com/jfhbrook/dosapp/config"
)

var linkCmd = &cobra.Command{
	Use:   "link [app]",
	Short: "Create a script that starts the app",
	Long:  `Create a script in DOSAPP_LINK_HOME that starts the application.`,
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

		if err := app.Run("link"); err != nil {
			log.Panic().Err(err).Msg("Failed to link application")
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
