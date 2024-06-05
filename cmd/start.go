/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/application"
	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/registry"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start the application",
	Long:  `Start the DOS application.`,
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		conf := config.NewConfig()

		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		reg := registry.NewRegistry(conf)
		pkg := reg.FindPackage(appName)
		app := application.NewApp(appName, pkg, conf)

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
