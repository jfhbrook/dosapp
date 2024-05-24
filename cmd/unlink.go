/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove a link",
	Long: `Remove a link created with 'dosapp link APP_NAME'.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unlink called")
  	fmt.Println("TODO: pull app name from args")
		fmt.Println("TODO: require that the app is installed")
		fmt.Println("TODO: refresh main configuration")
		fmt.Println("TODO: init app configuration")
		fmt.Println("TODO: refresh app configuration")
		fmt.Println("TODO: run remove-link task")
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)

	unlinkCmd.Flags().BoolP("refresh", "r", false, "Refresh task and conf files")
}
