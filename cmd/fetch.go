/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	// "github.com/jfhbrook/dosapp/application"
	// "github.com/jfhbrook/dosapp/config"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: pull github releases")
		log.Info().Msg("TODO: filter by app name")
		log.Info().Msg("TODO: choose latest release (by sort order)")
		log.Info().Msg("TODO: snag artifact url")
		log.Info().Msg("TODO: download artifact to ~/.cache/dosapp/packages")
		log.Info().Msg("TODO: extract artifact to ~/.local/share/dosapp/packages")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
