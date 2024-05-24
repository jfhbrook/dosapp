/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [app]",
	Short: "Create a script that starts the app",
	Long: `Create a script in DOSAPP_LINK_HOME that starts the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: pull app name from args")
		fmt.Println("TODO: require that the app is installed")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: init app configuration")
		fmt.Println("TODO: refresh app configuration")
		fmt.Println("TODO: run link task")
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
