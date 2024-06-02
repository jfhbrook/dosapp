/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package config

import (
	_ "embed"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/jfhbrook/dosapp/editor"
	"github.com/jfhbrook/dosapp/pager"
	"github.com/jfhbrook/dosapp/task"
)

//go:embed dosapp.env
var envFile []byte

//go:embed Taskfile.yml
var taskFile []byte

//go:embed main.conf.tmpl
var mainConfTmpl []byte

//go:embed start.conf.tmpl
var startConfTmpl []byte

var Templates = map[string]string{
	"main.conf.tmpl":  string(mainConfTmpl),
	"start.conf.tmpl": string(startConfTmpl),
}

type Config struct {
	// TODO: Root is only necessary until I move templating into dosapp
	Root         string
	ConfigHome   string
	LogLevel     string
	DosBoxBin    string
	SevenZipBin  string
	DataHome     string
	StateHome    string
	CacheHome    string
	DiskHome     string
	LinkHome     string
	PackageHome  string
	DownloadHome string
	DiskA        string
	DiskB        string
	DiskC        string
	Editor       *editor.Editor
	Pager        *pager.Pager
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func expandUser(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	usr, err := user.Current()

	if err != nil {
		return path, err
	}

	if path == "~" {
		return usr.HomeDir, nil
	}

	return filepath.Join(usr.HomeDir, path[2:]), nil
}

func mustExpandUser(path string) string {
	expanded, err := expandUser(path)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to expand user")
	}
	return expanded
}

func NewConfig() *Config {
	configHome := os.Getenv("DOSAPP_CONFIG_HOME")
	if configHome == "" {
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configHome = filepath.Join(xdgConfigHome, "dosapp")
		} else {
			configHome = filepath.Join(os.Getenv("HOME"), ".config", "dosapp")
		}
	}

	// TODO: This call overwrites any env vars that are set from outside!
	// This may be a sign to move away from dotenv and towards yaml.
	godotenv.Overload(filepath.Join(configHome, "dosapp.env"))

	root := getEnv("DOSAPP_ROOT", "")

	logLevel := getEnv("DOSAPP_LOG_LEVEL", "")
	dosBoxBin := getEnv("DOSAPP_DOSBOX_BIN", "dosbox-x")
	sevenZipBin := getEnv("DOSAPP_7Z_BIN", "7zz")

	dataHome := getEnv("DOSAPP_DATA_HOME", "")
	if dataHome == "" {
		if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
			dataHome = filepath.Join(xdgDataHome, "dosapp")
		} else {
			dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share", "dosapp")
		}
	}

	stateHome := os.Getenv("DOSAPP_STATE_HOME")
	if stateHome == "" {
		if xdgStateHome := os.Getenv("XDG_STATE_HOME"); xdgStateHome != "" {
			stateHome = filepath.Join(xdgStateHome, "dosapp")
		} else {
			stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state", "dosapp")
		}
	}

	cacheHome := os.Getenv("DOSAPP_CACHE_HOME")
	if cacheHome == "" {
		if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
			cacheHome = filepath.Join(xdgCacheHome, "dosapp")
		} else {
			cacheHome = filepath.Join(os.Getenv("HOME"), ".cache", "dosapp")
		}
	}

	diskHome := getEnv("DOSAPP_DISK_HOME", filepath.Join(os.Getenv("HOME"), "dosapp"))
	linkHome := getEnv("DOSAPP_LINK_HOME", filepath.Join(os.Getenv("HOME"), ".local", "bin"))

	// TODO: This needs to be manually overridden in the config to point to this
	// repo. Eventually, I'll implement package downloads from github releases.
	packageHome := getEnv("DOSAPP_PACKAGE_HOME", filepath.Join(stateHome, "packages"))

	downloadHome := getEnv("DOSAPP_DOWNLOAD_HOME", filepath.Join(cacheHome, "downloads"))
	diskA := getEnv("DOSAPP_DISK_A", filepath.Join(os.Getenv("HOME"), "Documents"))
	diskB := getEnv("DOSAPP_DISK_B", "")
	diskC := getEnv("DOSAPP_DISK_C", filepath.Join(os.Getenv("HOME"), "dosapp", "c"))

	edBin := os.Getenv("EDITOR")
	pgBin := os.Getenv("PAGER")

	ed := editor.NewEditor(os.Getenv("EDITOR"))
	pg := pager.NewPager(os.Getenv("PAGER"))

	conf := Config{
		mustExpandUser(root),
		mustExpandUser(configHome),
		logLevel,
		mustExpandUser(dosBoxBin),
		mustExpandUser(sevenZipBin),
		mustExpandUser(dataHome),
		mustExpandUser(stateHome),
		mustExpandUser(cacheHome),
		mustExpandUser(diskHome),
		mustExpandUser(linkHome),
		mustExpandUser(packageHome),
		mustExpandUser(downloadHome),
		mustExpandUser(diskA),
		mustExpandUser(diskB),
		mustExpandUser(diskC),
		ed,
		pg,
	}

	log.Debug().Str(
		"DOSAPP_LOG_LEVEL", conf.LogLevel,
	).Str(
		"DOSAPP_DOSBOX_BIN", conf.DosBoxBin,
	).Str(
		"DOSAPP_7Z_BIN", conf.SevenZipBin,
	).Str(
		"DOSAPP_DATA_HOME", conf.DataHome,
	).Str(
		"DOSAPP_STATE_HOME", conf.StateHome,
	).Str(
		"DOSAPP_CACHE_HOME", conf.CacheHome,
	).Str(
		"DOSAPP_LINK_HOME", conf.LinkHome,
	).Str(
		"DOSAPP_PACKAGE_HOME", conf.PackageHome,
	).Str(
		"DOSAPP_DOWNLOAD_HOME", conf.DownloadHome,
	).Str(
		"DOSBOX_DISK_A", conf.DiskA,
	).Str(
		"DOSBOX_DISK_B", conf.DiskB,
	).Str(
		"DOSBOX_DISK_C", conf.DiskC,
	).Str(
		"EDITOR", edBin,
	).Str(
		"PAGER", pgBin,
	).Msg("Loaded config")

	return &conf
}

func (conf *Config) Env() map[string]string {
	env := map[string]string{
		"DOSAPP_ROOT":          conf.Root,
		"DOSAPP_CONFIG_HOME":   conf.ConfigHome,
		"DOSAPP_LOG_LEVEL":     conf.LogLevel,
		"DOSAPP_DOSBOX_BIN":    conf.DosBoxBin,
		"DOSAPP_7Z_BIN":        conf.SevenZipBin,
		"DOSAPP_DATA_HOME":     conf.DataHome,
		"DOSAPP_STATE_HOME":    conf.StateHome,
		"DOSAPP_CACHE_HOME":    conf.CacheHome,
		"DOSAPP_DISK_HOME":     conf.DiskHome,
		"DOSAPP_LINK_HOME":     conf.LinkHome,
		"DOSAPP_PACKAGE_HOME":  conf.PackageHome,
		"DOSAPP_DOWNLOAD_HOME": conf.DownloadHome,
		"DOSAPP_DISK_A":        conf.DiskA,
		"DOSAPP_DISK_B":        conf.DiskB,
		"DOSAPP_DISK_C":        conf.DiskC,
		"EDITOR":               conf.Editor.Bin,
		"PAGER":                conf.Pager.Bin,
	}

	return env
}

func (conf *Config) Environ() []string {
	env := []string{
		"DOSAPP_ROOT=" + conf.Root,
		"DOSAPP_CONFIG_HOME=" + conf.ConfigHome,
		"DOSAPP_LOG_LEVEL=" + conf.LogLevel,
		"DOSAPP_DOSBOX_BIN=" + conf.DosBoxBin,
		"DOSAPP_7Z_BIN=" + conf.SevenZipBin,
		"DOSAPP_DATA_HOME=" + conf.DataHome,
		"DOSAPP_STATE_HOME=" + conf.StateHome,
		"DOSAPP_CACHE_HOME=" + conf.CacheHome,
		"DOSAPP_DISK_HOME=" + conf.DiskHome,
		"DOSAPP_LINK_HOME=" + conf.LinkHome,
		"DOSAPP_PACKAGE_HOME=" + conf.PackageHome,
		"DOSAPP_DOWNLOAD_HOME=" + conf.DownloadHome,
		"DOSAPP_DISK_A=" + conf.DiskA,
		"DOSAPP_DISK_B=" + conf.DiskB,
		"DOSAPP_DISK_C=" + conf.DiskC,
		"EDITOR=" + conf.Editor.Bin,
		"PAGER=" + conf.Pager.Bin,
	}

	return append(os.Environ(), env...)
}

func (conf *Config) EnvFilePath() string {
	return filepath.Join(conf.ConfigHome, "dosapp.env")
}

func (conf *Config) WriteEnvFile() error {
	envPath := conf.EnvFilePath()
	return os.WriteFile(envPath, envFile, 0644)
}

func (conf *Config) EditEnvFile() error {
	envPath := conf.EnvFilePath()
	return conf.Editor.Edit(envPath)
}

func (conf *Config) EnvFileExists() bool {
	envPath := conf.EnvFilePath()
	_, err := os.Stat(envPath)
	return err == nil
}

func (conf *Config) TaskFilePath() string {
	return filepath.Join(conf.ConfigHome, "Taskfile.yml")
}

func (conf *Config) WriteTaskFile() error {
	taskPath := conf.TaskFilePath()
	return os.WriteFile(taskPath, taskFile, 0644)
}

func (conf *Config) TaskFileExists() bool {
	taskPath := conf.TaskFilePath()
	_, err := os.Stat(taskPath)
	return err == nil
}

func (conf *Config) Refresh() error {
	err := conf.WriteTaskFile()
	if err != nil {
		return err
	}

	return conf.Run("init")
}

func (conf *Config) Run(args ...string) error {
	taskPath := conf.TaskFilePath()
	return task.Run(taskPath, conf.Environ(), args...)
}
