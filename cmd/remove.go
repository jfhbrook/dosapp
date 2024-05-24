/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [app]",
	Short: "Remove an application",
	Long:  `Unlink the app and remove its configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: pull app name from args")
		log.Info().Msg("TODO: require that the app is installed")
		log.Info().Msg("TODO: refresh main configuration")
		log.Info().Msg("TODO: init app configuration")
		log.Info().Msg("TODO: refresh app configuration")
		log.Info().Msg("TODO: run remove-link task")
		log.Info().Msg("TODO: run remove task")
		log.Info().Msg("TODO: rm -rf the app")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	removeCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
