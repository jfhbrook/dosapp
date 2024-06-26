/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/config"
)

var rootCmd = &cobra.Command{
	Use:     "dosapp",
	Version: "2.0.0",
	Short:   "Manage DOSBox applications",
	Long: `Install, start and link DOSBox applications.

When no subcommand is given, dosapp will start DOSBox with the default
configuration.

To refresh the default configuration, run 'dosapp config'.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		level := cmd.Flag("log-level").Value.String()

		if level == "" {
			level = os.Getenv("DOSAPP_LOG_LEVEL")
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.NewConfig()

		if err := conf.Run("start"); err != nil {
			log.Panic().Err(err).Msg("Failed to start")
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("toodles!")
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
