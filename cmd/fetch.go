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

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a package from the registry",
	Long: `Fetch a package from the registry without installing it. This step
will run automatically on 'dosapp install'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		conf := config.NewConfig()

		forceFlag, _ := cmd.Flags().GetBool("force")
		updateFlag, _ := cmd.Flags().GetBool("update")

		reg := registry.NewRegistry(conf)
		pkg := reg.FindPackage(appName)

		if !pkg.RemotePackageExists() {
			log.Fatal().Str("package", appName).Msg("Package not found.")
		}

		updateDirections := false
		forceDirections := false

		if pkg.HasUpdate() {
			if updateFlag {
				forceFlag = true
			} else {
				updateDirections = true
			}
		}

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

		if updateDirections {
			log.Warn().Str(
				"package", appName,
			).Msgf("To update %s, use the --update flag", appName)
		} else if forceDirections {
			log.Warn().Str(
				"package", appName,
			).Msgf("To force fetching %s, use the --force flag", appName)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().BoolP("force", "f", false, "Force installation of package")
	fetchCmd.Flags().BoolP("update", "U", false, "Update the package")
}
