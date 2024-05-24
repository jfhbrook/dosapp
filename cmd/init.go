/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dosapp's configuration",
	Long: `Initialize dosapp's main configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Create env file")
		fmt.Println("TODO: Edit env file")
		fmt.Println("TODO: Reload config")
		fmt.Println("TODO: Refresh main config")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("refresh", "r", false, "Generate new task and conf files")
}
