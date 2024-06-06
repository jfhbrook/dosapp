/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	application "github.com/jfhbrook/dosapp/app"
	"github.com/jfhbrook/dosapp/config"
)

var installCmd = &cobra.Command{
	Use:   "install [app]",
	Short: "Install an application",
	Long: `Install an application. Download the package from the registry,
set up configuration for the app, and run its installer.

To fetch the package without installing, run 'dosapp fetch [app]'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: We can't actually just run this command directly, because we
		// also need to do the update/force flag song/dance here. Both of these
		// commands should call a function that takes the boolean flags and
		// returns new boolean flags.
		fetchCmd.Run(cmd, args)

		appName := args[0]
		conf := config.NewConfig()

		// updateFlag, _ := cmd.Flags().GetBool("update")
		forceFlag, _ := cmd.Flags().GetBool("force")
		editFlag, _ := cmd.Flags().GetBool("edit")
		editFlagChanged := cmd.Flags().Changed("edit")
		docsFlag, _ := cmd.Flags().GetBool("docs")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		refreshFlag, _ := cmd.Flags().GetBool("refresh")

		refreshFlag = refreshFlag || forceFlag

		if overwriteFlag {
			refreshFlag = true
		}

		app := application.NewApp(appName, conf)

		if err := app.Mkdir(); err != nil {
			log.Panic().Err(err).Msg("Failed to create app directory")
		}

		if !overwriteFlag && app.EnvFileExists() {
			log.Warn().Msgf("Environment file already exists at %s", app.EnvFilePath())
			log.Warn().Msgf("To overwrite and refresh the configuration, run 'dosapp install %s --overwrite'", appName)
		} else {
			refreshFlag = true
			if err := app.WriteEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to write env file")
			}
		}

		if !refreshFlag && app.TaskFileExists() {
			log.Warn().Msgf("Taskfile already exists at %s", app.TaskFilePath())
			log.Warn().Msgf("To refresh the app configuration, run 'dosapp install %s --refresh'", appName)
		} else {
			refreshFlag = true
			// Will write task file on refresh step
		}

		// If we're not refreshing, then we only want to edit the file if
		// explicitly asked
		var shouldEdit bool

		if refreshFlag {
			shouldEdit = editFlag
		} else {
			shouldEdit = editFlag && editFlagChanged
		}

		if shouldEdit {
			if err := app.EditEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to edit application config file")
			}
		}

		if refreshFlag {
			if err := app.Refresh(); err != nil {
				log.Panic().Err(err).Msg("Failed to refresh application config")
			}
		}

		if docsFlag {
			if err := app.ShowDocs(); err != nil {
				log.Error().Err(err).Msg("Failed to show docs")
			}
		}

		if err := app.Run("install"); err != nil {
			log.Panic().Err(err).Msg("Failed to install application")
		}
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
	installCmd.Flags().BoolP("update", "U", false, "Update the package")
}
