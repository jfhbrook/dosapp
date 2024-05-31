package application

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/task"
)

type App struct {
	Name   string
	Config *config.Config
}

func NewApp(conf *config.Config, name string) *App {
	app := App{
		Name:   name,
		Config: conf,
	}

	return &app
}

func (app *App) Path() string {
	return filepath.Join(app.Config.ConfigHome, "apps", app.Name)
}

func (app *App) Mkdir() error {
	return os.MkdirAll(app.Path(), 0755)
}

func (app *App) Exists() bool {
	_, err := os.Stat(app.Path())
	return err == nil
}

func (app *App) Env() map[string]string {
	// TODO: I don't like this call, nor the one in the next function. They
	// should be fine given a command only ever creates one App, but mutating
	// global state for a non-global config makes me unhappy.
	godotenv.Overload(app.EnvFilePath())

	// godotenv mutated its global state, so we can do this. It's fine.
	// No really.
	return app.Config.Env()
}

func (app *App) Environ() []string {
	godotenv.Overload(app.EnvFilePath())

	return app.Config.Environ()
}

func (app *App) EnvFilePath() string {
	return filepath.Join(app.Path(), "dosapp.env")
}

func (app *App) EnvFileTemplatePath() string {
	return filepath.Join(app.Config.PackageHome, app.Name, "dosapp.env.tmpl")
}

func (app *App) WriteEnvFile() error {
	tmplPath := app.EnvFileTemplatePath()
	log.Warn().Msgf("Using template at %s", tmplPath)
	tmpl, err := template.New("dosapp.env.tmpl").ParseFiles(tmplPath)

	if err != nil {
		return err
	}

	envFilePath := app.EnvFilePath()
	var f *os.File
	f, err = os.Create(envFilePath)
	defer f.Close()

	if err != nil {
		return err
	}

	data := map[string]map[string]string{
		"Env": app.Env(),
	}

	return tmpl.Execute(f, data)
}

func (app *App) EditEnvFile() error {
	envPath := app.EnvFilePath()
	return app.Config.Editor.Edit(envPath)
}

func (app *App) EnvFileExists() bool {
	envPath := app.EnvFilePath()
	_, err := os.Stat(envPath)
	return err == nil
}

func (app *App) TaskFilePath() string {
	return filepath.Join(app.Path(), "Taskfile.yml")
}

func (app *App) TaskFileExists() bool {
	taskPath := app.TaskFilePath()
	_, err := os.Stat(taskPath)
	return err == nil
}

func (app *App) WriteTaskFile() error {
	tmplPath := filepath.Join(app.Config.PackageHome, app.Name, "Taskfile.yml")

	taskFile, err := os.ReadFile(tmplPath)

	if err != nil {
		return err
	}

	taskPath := app.TaskFilePath()
	return os.WriteFile(taskPath, taskFile, 0644)
}

func (app *App) DocsPath() string {
	return filepath.Join(app.Path(), "README.md")
}

func (app *App) WriteDocs() error {
	tmplPath := filepath.Join(app.Config.PackageHome, app.Name, "README.md")

	docs, err := os.ReadFile(tmplPath)

	if err != nil {
		return err
	}

	docsPath := app.DocsPath()
	return os.WriteFile(docsPath, docs, 0644)
}

func (app *App) ShowDocs() error {
	docsPath := app.DocsPath()
	return app.Config.Pager.Show(docsPath)
}

func (app *App) Refresh() error {
	if err := app.WriteTaskFile(); err != nil {
		return err
	}

	if err := app.WriteDocs(); err != nil {
		return err
	}

	return app.Run("init")
}

func (app *App) Run(args ...string) error {
	taskPath := app.TaskFilePath()
	return task.Run(taskPath, app.Environ(), args...)
}
