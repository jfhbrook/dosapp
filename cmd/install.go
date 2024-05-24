/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an application",
	Long: `Set up the package and configuration for the app, and run its installer.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: pull app name from args")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: init app configuration")
		fmt.Println("TODO: refresh app configuration")
		fmt.Println("TODO: show readme")
		fmt.Println("TODO: run install task")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolP("edit", "e", true, "Edit environment files")
	installCmd.Flags().BoolP("overwrite", "o", false, "Overwrite existing configuration")
	installCmd.Flags().BoolP("docs", "d", true, "Display the README")
	installCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
