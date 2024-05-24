/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "dosapp",
	Short: "Manage DOSBox applications",
	Long:  `Install, run and link DOSBox applications using task and go templates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	},
	Run: func(cmd *cobra.Command, args []string) {
		level := cmd.Flag("log-level").Value.String()
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
			log.Fatal().Str("log-level", level).Msg("Invalid log level")
		}

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
	rootCmd.PersistentFlags().String("log-level", "info", "Logging level")
	rootCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
