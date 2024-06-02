/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	// "github.com/jfhbrook/dosapp/application"
	// "github.com/jfhbrook/dosapp/config"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a config or package template",
	Long:  `Render a config or package template.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// templateName := args[0]
		// conf := config.NewConfig()

		configFlag, _ := cmd.Flags().GetBool("config")
		packageName := cmd.Flag("package").Value.String()

		// var dest string

		if packageName != "" && configFlag {
			log.Fatal().Msg("Cannot specify both a package and a config template")
		} else if packageName != "" {
			log.Info().Msg("TODO: load app object")
			log.Info().Msg("TODO: generate app env object")
			log.Info().Msg("TODO: load template from package")
			log.Info().Msg("TODO: set destination in ~/.config/dosapp/apps/{app}")
		} else if configFlag {
			log.Info().Msg("TODO: generate config env object")
			log.Info().Msg("TODO: pull from map of templates")
			log.Info().Msg("TODO: set destination in ~/.config/dosapp")
		} else {
			log.Fatal().Msg("Must specify either a package or a config template")
		}

		log.Info().Msg("TODO: render template to destination")
	},
}

func init() {
	taskCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringP("package", "p", "", "Package that contains template")
	templateCmd.Flags().BoolP("config", "C", false, "Use a config template")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
