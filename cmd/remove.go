/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an application",
	Long: `Unlink the app and remove its configuration.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: pull app name from args")
		fmt.Println("TODO: require that the app is installed")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: init app configuration")
		fmt.Println("TODO: refresh app configuration")
		fmt.Println("TODO: run remove-link task")
		fmt.Println("TODO: run remove task")
		fmt.Println("TODO: rm -rf the app")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	removeCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
