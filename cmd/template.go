/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a config or package template",
	Long:  `Render a config or package template.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO: load config")
		log.Info().Msg("TODO: generate env object")
		log.Info().Msg("TODO: if config template, pull from map of templates")
		log.Info().Msg("TODO: if package template, load template from package")
		log.Info().Msg("TODO: if config template, destination in ~/.config/dosapp")
		log.Info().Msg("TODO: if package template, destination in ~/.config/dosapp/apps/{app}")
		log.Info().Msg("TODO: render template to destination")
	},
}

func init() {
	taskCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
