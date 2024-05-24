/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start the application",
	Long:  `Start the DOS application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: pull app name from args")
		log.Info().Msg("TODO: require that the app is installed")
		log.Info().Msg("TODO: refresh main configuration")
		log.Info().Msg("TODO: refresh main configuration")
		log.Info().Msg("TODO: init app configuration")
		log.Info().Msg("TODO: refresh app configuration")
		log.Info().Msg("TODO: run start task")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
