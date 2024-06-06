/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	application "github.com/jfhbrook/dosapp/app"
	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/pacman"
)

type FetchResult struct {
	// True when we called --update but there was no update
	NoUpdate bool
	// Signals a refresh to the installer
	Refresh bool
}

func Fetch(cmd *cobra.Command, args []string) *FetchResult {
	appName := args[0]
	conf := config.NewConfig()

	forceFlag, _ := cmd.Flags().GetBool("force")
	updateFlag, _ := cmd.Flags().GetBool("update")
	refreshFlag, _ := cmd.Flags().GetBool("refresh")

	noUpdate := false

	reg := pacman.NewRegistry(conf)
	pkg := reg.FindPackage(appName)

	if !pkg.RemotePackageExists() {
		log.Fatal().Str("package", appName).Msg("Package not found.")
	}

	updateDirections := false
	forceDirections := false

	// Update forces install only if there's an update
	if pkg.HasUpdate() {
		if updateFlag {
			forceFlag = true
		} else {
			updateDirections = true
		}
	} else if updateFlag {
		noUpdate = true
	}

	// If no local package (or forcing), fetch/unpack and signal refresh
	if !pkg.LocalPackageExists() || forceFlag {
		refreshFlag = true

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
	} else if noUpdate {
		log.Info().Str(
			"package", appName,
		).Msgf("No update available for %s", appName)
	}

	return &FetchResult{
		NoUpdate: noUpdate,
		Refresh:  refreshFlag,
	}
}

var installCmd = &cobra.Command{
	Use:   "install [app]",
	Short: "Install an application",
	Long: `Install an application. Download the package from the registry,
set up configuration for the app, and run its installer.

The --update flag will update the package if a new version is available.

The --force flag will force the package to update, even if a new version is
not available.

The --overwrite flag will overwrite the application configuration, if it
exists. By default, package updates do not overwrite the application
configuration.

The --refresh flag will generate a fresh Taskfile.yml and DOXBox .conf files
based on the application configuration. This is implied by --force, --update
and --overwrite.

To skip running the DOS installer for the application, set --installer=false.

By default, dosapp will open the configuration file in an editor when
installing a new application or executing --refresh. To skip opening the
application's configuration, set --edit=false. To edit the configuration even
when not refreshing, set --edit=true.

By default, dosapp will display the package README in a pager before running
the installer. To skip opening the README, set --docs=false. To display the
README even when --installer=false, set --docs=true.

To fetch the package without installing, run 'dosapp fetch [app]'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch the package, if appropriate
		result := Fetch(cmd, args)

		editFlag, _ := cmd.Flags().GetBool("edit")
		editFlagChanged := cmd.Flags().Changed("edit")
		docsFlag, _ := cmd.Flags().GetBool("docs")
		docsFlagChanged := cmd.Flags().Changed("docs")
		overwriteFlag, _ := cmd.Flags().GetBool("overwrite")
		installerFlag, _ := cmd.Flags().GetBool("installer")

		// Overwrite implies refresh, otherwise we can trust Fetch to give us
		// the right value
		refreshFlag := overwriteFlag || result.Refresh

		// If there wasn't an update and we're not explicitly refreshing, we're
		// actually done.
		if result.NoUpdate && !refreshFlag {
			return
		}

		appName := args[0]
		conf := config.NewConfig()
		app := application.NewApp(appName, conf)

		if err := app.Mkdir(); err != nil {
			log.Panic().Err(err).Msg("Failed to create app directory")
		}

		// Only overwrite the env file if we're explicitly told to, even in the
		// cases of --update, --force and --refresh
		if !overwriteFlag && app.EnvFileExists() {
			log.Warn().Msgf("Environment file already exists at %s", app.EnvFilePath())
			log.Warn().Msgf("To overwrite and refresh the configuration, run 'dosapp install %s --overwrite'", appName)
		} else {
			// Overwrite implies refresh
			refreshFlag = true
			if err := app.WriteEnvFile(); err != nil {
				log.Panic().Err(err).Msg("Failed to write env file")
			}
		}

		if !refreshFlag && app.TaskFileExists() {
			log.Warn().Msgf("Taskfile already exists at %s", app.TaskFilePath())
			log.Warn().Msgf("To refresh the app configuration, run 'dosapp install %s --refresh'", appName)
		} else {
			// If no task file, we have to refresh anyway
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

		// Refresh the app config if:
		// - We passed --refresh
		// - We passed --update and there was an update
		// - We passed --force
		// - We passed --overwrite
		if refreshFlag {
			if err := app.Refresh(); err != nil {
				log.Panic().Err(err).Msg("Failed to refresh application config")
			}
		}

		// If we're not running the installer, then we only want to show the
		// README if explicitly asked
		var shouldShow bool

		if installerFlag {
			shouldShow = docsFlag
		} else {
			shouldShow = docsFlag && docsFlagChanged
		}

		if shouldShow {
			if err := app.ShowDocs(); err != nil {
				log.Error().Err(err).Msg("Failed to show docs")
			}
		}

		if installerFlag {
			if err := app.Run("install"); err != nil {
				log.Panic().Err(err).Msg("Failed to install application")
			}
		}
	},
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a package from the registry",
	Long: `Fetch a package from the registry without installing it. This step
will run automatically on 'dosapp install'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Fetch(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().BoolP("force", "f", false, "Force package install")
	fetchCmd.Flags().BoolP("update", "U", false, "Update the package")

	installCmd.Flags().BoolP("force", "f", false, "Force package install")
	installCmd.Flags().BoolP("update", "U", false, "Update the package")
	installCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	installCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files (implied by force and update)")

	installCmd.Flags().BoolP("edit", "e", true, "Edit environment files")
	installCmd.Flags().BoolP("docs", "d", true, "Display the README")
	installCmd.Flags().BoolP("installer", "I", true, "Run the DOS installer")
}
