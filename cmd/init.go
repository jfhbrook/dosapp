/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dosapp's configuration",
	Long:  `Initialize dosapp's main configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: Create env file")
		log.Info().Msg("TODO: Edit env file")
		log.Info().Msg("TODO: Reload config")
		log.Info().Msg("TODO: Refresh main config")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
