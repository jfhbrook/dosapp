/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink [app]",
	Short: "Remove a link",
	Long: `Remove a link created with 'dosapp link [app]'.`,
	Run: func(cmd *cobra.Command, args []string) {
  	log.Info().Msg("TODO: pull app name from args")
		log.Info().Msg("TODO: require that the app is installed")
		log.Info().Msg("TODO: refresh main configuration")
		log.Info().Msg("TODO: init app configuration")
		log.Info().Msg("TODO: refresh app configuration")
		log.Info().Msg("TODO: run remove-link task")
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)

	unlinkCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
