package application

import (
	"os"
	"path/filepath"

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
	return os.MkdirAll(app.Path(), os.ModeDir)
}

func (app *App) Exists() bool {
	_, err := os.Stat(app.Path())
	return err == nil
}

func (app *App) Environ() []string {
	// TODO: Env file loading currently depends on the Taskfile referencing
	// the file. But it would be nice to read/parse the env file directly.
	return app.Config.Environ()
}

func (app *App) EnvFilePath() string {
	return filepath.Join(app.Path(), "dosapp.env")
}

func (app *App) EnvFileTemplatePath() string {
	return filepath.Join(app.Config.PackageHome, app.Name, "dosapp.env.tmpl")
}

// TODO: Read the template from the package directory, and template it out
// with an .Env object.
//
// The truth is I didn't want to bite off templating right now. I might
// exec a call to gomplate just to get this unblocked.
func (app *App) ReadEnvFileTemplate() ([]byte, error) {
	panic("ReadEnvFileTemplate not implemented")
}

func (app *App) WriteEnvFile() error {
	envPath := app.EnvFilePath()
	envFile, err := app.ReadEnvFileTemplate()
	if err != nil {
		return err
	}
	return os.WriteFile(envPath, envFile, 0644)
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

func (app *App) TaskFileTemplatePath() string {
	return filepath.Join(app.Config.PackageHome, app.Name, "Taskfile.yml.tmpl")
}

func (app *App) TaskFileExists() bool {
	taskPath := app.TaskFilePath()
	_, err := os.Stat(taskPath)
	return err == nil
}

// TODO: Same deal as the env file template
func (app *App) ReadTaskFileTemplate() ([]byte, error) {
	panic("ReadTaskFileTemplate not implemented")
}

func (app *App) WriteTaskFile() error {
	taskPath := app.TaskFilePath()
	taskFile, err := app.ReadTaskFileTemplate()
	if err != nil {
		return err
	}
	return os.WriteFile(taskPath, taskFile, 0644)
}

func (app *App) DocsPath() string {
	return filepath.Join(app.Path(), "README.md")
}

func (app *App) CopyDocs() error {
	panic("CopyDocs not implemented")
}

func (app *App) ShowDocs() error {
	docsPath := app.DocsPath()
	return app.Config.Pager.Show(docsPath)
}

func (app *App) Refresh() error {
	if err := app.WriteTaskFile(); err != nil {
		return err
	}

	if err := app.CopyDocs(); err != nil {
		return err
	}

	return app.Run("init")
}

func (app *App) Run(args ...string) error {
	taskPath := app.TaskFilePath()
	return task.Run(taskPath, app.Environ(), args...)
}
