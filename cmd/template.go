/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/application"
	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/packages"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a config or package template",
	Long:  `Render a config or package template.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		conf := config.NewConfig()

		configFlag, _ := cmd.Flags().GetBool("config")
		packageName := cmd.Flag("package").Value.String()

		basename := strings.TrimSuffix(templateName, filepath.Ext(templateName))

		var env map[string]string
		var tmpl *template.Template
		var dest string

		if packageName != "" && configFlag {
			log.Fatal().Msg("Cannot specify both a package and a config template")
		} else if packageName != "" {
			app := application.NewApp(conf, packageName)
			pkg := packages.NewPackage(conf, packageName)
			src := filepath.Join(pkg.Path(), templateName)

			env = app.Env()
			dest = filepath.Join(app.Path(), basename)

			var err error

			tmpl, err = template.New(templateName).ParseFiles(src)

			if err != nil {
				log.Fatal().Err(err).Msg("Failed to parse package template")
			}
		} else if configFlag {
			env = conf.Env()
			dest = filepath.Join(conf.ConfigHome, templateName)

			var err error
			tmpl, err = template.New(templateName).Parse(config.Templates[templateName])

			if err != nil {
				log.Panic().Err(err).Msg("Failed to parse config template")
			}
		} else {
			log.Fatal().Msg("Must specify either a package or a config template")
		}

		var f *os.File
		var err error

		f, err = os.Create(dest)
		defer f.Close()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open destination file")
		}

		data := map[string]map[string]string{
			"Env": env,
		}

		err = tmpl.Execute(f, data)

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to render template")
		}
	},
}

func init() {
	taskCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringP("package", "p", "", "Package that contains template")
	templateCmd.Flags().BoolP("config", "C", false, "Use a config template")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
