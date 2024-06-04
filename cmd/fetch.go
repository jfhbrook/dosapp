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
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.NewConfig()
		reg := registry.NewRegistry(conf)

		// if !pkg.Fetched() || forceFlag {
		//   if err := pkg.Fetch(); err != nil {
		//		 log.Fatal().Err(err).Msg("failed to fetch package")
		//   }
		// }
		//
		// if !pkg.Extracted() || forceFlag {
		//   if err := pkg.Extract(); err != nil {
		// 	   log.Fatal().Err(err).Msg("failed to download package")
		//   }
		// }

		pkg := reg.FindPackage("wordperfect")

		/*
			if !pkg.LocalPackageExists() {
				log.Warn().Msg("No local version found")
			} else {
				log.Warn().Msg(pkg.LocalVersion.String())
				log.Warn().Msg(pkg.LocalReleaseVersion.String())
			}
		*/

		if !pkg.RemotePackageExists() {
			log.Fatal().Msg("Package not found")
		}

		if !pkg.LocalArtifactExists() {
			if err := pkg.Fetch(); err != nil {
				log.Fatal().Err(err).Msg("Failed to fetch package")
			}
		} else {
			log.Info().Msg("To download the package again, use the --force flag")
		}

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
