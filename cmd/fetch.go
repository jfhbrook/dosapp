/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/registry"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		conf := config.NewConfig()

		forceFlag, _ := cmd.Flags().GetBool("force")

		reg := registry.NewRegistry(conf)
		pkg := reg.FindPackage(appName)

		forceDirections := false

		if !pkg.LocalPackageExists() || forceFlag {
			if !pkg.LocalArtifactExists() || forceFlag {
				if err := pkg.Fetch(); err != nil {
					log.Fatal().Err(err).Msg("Failed to fetch package from registry")
				}
			} else {
				forceDirections = true
			}

			if err := pkg.Unpack(); err != nil {
				log.Fatal().Err(err).Msg("Failed to unpack package artifact")
			}
		} else {
			forceDirections = true
		}

		if forceDirections {
			log.Info().Msg("To fetch the latest package, use the --force flag")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().BoolP("force", "f", false, "Force installation of package")
}
