/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/application"
	"github.com/jfhbrook/dosapp/config"
)

var installCmd = &cobra.Command{
	Use:   "install [app]",
	Short: "Install an application",
	Long:  `Set up the package and configuration for the app, and run its installer.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		conf := config.LoadConfig()
		editFlag, _ := cmd.Flags().GetBool("edit")
		docsFlag, _ := cmd.Flags().GetBool("docs")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		if overwriteFlag {
			refreshFlag = true
		}

		if refreshFlag {
			if err := conf.Refresh(); err != nil {
				log.Panic().Err(err).Msg("Failed to refresh config")
			}
		}

		app := application.LoadApp(&conf, appName)
		if overwriteFlag || !app.EnvFileExists() {
			app.WriteEnvFile()
			if editFlag {
				app.EditEnvFile()
			}
		}

		if !refreshFlag && app.TaskFileExists() {
			log.Warn().Msgf("Taskfile already exists at %s", app.TaskFilePath)
			log.Warn().Msgf("To refresh the app configuration, run 'dosapp install %s --refresh'", appName)
		} else {
			refreshFlag = true
		}

		if refreshFlag {
			app.Refresh()
		}

		if docsFlag {
			app.ShowDocs()
		}

		app.Run("install")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// NOTE: Bash dosapp uses --no-edit and --no-docs here. In this case, you
	// would need to call --edit=false or -e=false. It's not my favorite, but
	// seems to play well with go idioms.
	installCmd.Flags().BoolP("edit", "e", true, "Edit environment files")
	installCmd.Flags().BoolP("docs", "d", true, "Display the README")
	installCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	installCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
