package application

import (
	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/task"
)

type App struct {
	name   string
	config *config.Config
}

func LoadApp(conf *config.Config, name string) App {
	return App{
		name,
		conf,
	}
}

func (app *App) Environ() []string {
	// TODO: Env file loading currently depends on the Taskfile referencing
	// the file. But it would be nice to read/parse the env file directly.
	return app.config.Environ()
}

func (app *App) EnvFilePath() string {
	return ""
}

func (app *App) WriteEnvFile() {
}

func (app *App) EditEnvFile() {
}

func (app *App) EnvFileExists() bool {
	return false
}

func (app *App) TaskFilePath() string {
	panic("TaskFilePath not implemented")
	return ""
}

func (app *App) TaskFileExists() bool {
	return false
}

func (app *App) Refresh() {
}

func (app *App) ShowDocs() {
}

func (app *App) Run(args ...string) error {
	taskPath := app.TaskFilePath()
	return task.Run(taskPath, app.Environ(), args...)
}
