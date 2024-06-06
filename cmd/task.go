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

	application "github.com/jfhbrook/dosapp/app"
	"github.com/jfhbrook/dosapp/config"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Run a build task",
	Long: `Run a build task. This command is intended for use in package
Taskfiles.`,
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render a config or package template",
	Long:  `Render a config or package template.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templatePath := args[0]
		conf := config.NewConfig()

		configFlag, _ := cmd.Flags().GetBool("config")
		linkFlag, _ := cmd.Flags().GetBool("link")
		packageName := cmd.Flag("package").Value.String()

		templateName := filepath.Base(templatePath)
		fileName := strings.TrimSuffix(templateName, filepath.Ext(templateName))

		var env map[string]string
		var tmpl *template.Template
		var dest string

		if packageName != "" && configFlag {
			log.Fatal().Msg("Cannot specify both a package and a config template")
		} else if packageName != "" {
			var err error

			app := application.NewApp(packageName, conf)
			src := filepath.Join(app.Package.LocalPackagePath(), templatePath)

			env = app.Env()

			if linkFlag {
				var linkPath string
				linkPath, err = filepath.Rel("bin", templatePath)

				if err != nil {
					log.Panic().Err(err).Msg("Failed to get relative path")
				}

				linkPath = strings.TrimSuffix(linkPath, filepath.Ext(linkPath))

				dest = filepath.Join(conf.LinkHome, linkPath)
			} else {
				dest = filepath.Join(app.Path(), fileName)
			}

			log.Debug().Str(
				"name", templateName,
			).Str(
				"src", src,
			).Msg("Parsing package template")

			tmpl, err = template.New(templateName).ParseFiles(src)

			if err != nil {
				log.Fatal().Err(err).Msg("Failed to parse package template")
			}
		} else if configFlag {
			if linkFlag {
				log.Fatal().Msg("Cannot link a config template")
			}
			env = conf.Env()
			dest = filepath.Join(conf.ConfigHome, templatePath)

			log.Debug().Str(
				"name", templateName,
			).Str(
				"src", templatePath,
			).Msg("Parsing config template")

			var err error
			tmpl, err = template.New(templateName).Parse(config.Templates[templatePath])

			if err != nil {
				log.Panic().Err(err).Msg("Failed to parse config template")
			}
		} else {
			log.Fatal().Msg("Must specify either a package or a config template")
		}

		log.Debug().Str(
			"destination", dest,
		).Msg("Rendering template")

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
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringP("package", "p", "", "Package that contains template")
	templateCmd.Flags().BoolP("config", "C", false, "Use a config template")
	templateCmd.Flags().BoolP("link", "L", false, "Template is a link")
}
