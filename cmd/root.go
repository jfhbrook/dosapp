/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/jfhbrook/dosapp/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "dosapp",
	Version: "2.0.0",
	Short:   "Manage DOSBox applications",
	Long:    `Install, run and link DOSBox applications using task and go templates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.LoadConfig()
		level := cmd.Flag("log-level").Value.String()
		if level == "" {
			level = conf.LogLevel
		}
		if level == "" {
			level = "info"
		}
		switch level {
		case "trace":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			log.Fatal().Str("DOSAPP_LOG_LEVEL", level).Msg("Invalid log level")
		}

		log.Debug().Str(
			"DOSAPP_LOG_LEVEL", conf.LogLevel,
		).Str(
			"DOSAPP_DOSBOX_BIN", conf.DosBoxBin,
		).Str(
			"DOSAPP_7Z_BIN", conf.SevenZipBin,
		).Str(
			"DOSAPP_DATA_HOME", conf.DataHome,
		).Str(
			"DOSAPP_STATE_HOME", conf.StateHome,
		).Str(
			"DOSAPP_CACHE_HOME", conf.CacheHome,
		).Str(
			"DOSAPP_LINK_HOME", conf.LinkHome,
		).Str(
			"DOSAPP_PACKAGE_HOME", conf.PackageHome,
		).Str(
			"DOSAPP_DOWNLOAD_HOME", conf.DownloadHome,
		).Str(
			"PAGER", conf.Pager,
		).Str(
			"DOSBOX_DISK_A", conf.DiskA,
		).Str(
			"DOSBOX_DISK_B", conf.DiskB,
		).Str(
			"DOSBOX_DISK_C", conf.DiskC,
		).Msg("Loaded config")

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
	rootCmd.PersistentFlags().String("log-level", "", "Logging level")
	rootCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
