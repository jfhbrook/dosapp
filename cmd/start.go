/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start the application",
	Long: `Start the DOS application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: pull app name from args")
		fmt.Println("TODO: require that the app is installed")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: init app configuration")
		fmt.Println("TODO: refresh app configuration")
		fmt.Println("TODO: run start task")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
