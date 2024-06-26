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

var removeCmd = &cobra.Command{
	Use:   "remove [app]",
	Short: "Remove an application",
	Long:  `Unlink the app and remove its configuration.`,
	Args:  cobra.ExactArgs(1),
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

		if err := app.Run("unlink"); err != nil {
			log.Panic().Err(err).Msg("Failed to unlink application")
		}

		if err := app.Run("remove"); err != nil {
			log.Panic().Err(err).Msg("Failed to remove application")
		}

		if err := app.Remove(); err != nil {
			log.Panic().Err(err).Msg("Failed to remove application")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	removeCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
