/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"os"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dosapp",
	Short: "Manage DOSBox applications",
	Long: `Install, run and link DOSBox applications using task and go templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: refresh main")
		log.Info().Msg("TODO: start main")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
