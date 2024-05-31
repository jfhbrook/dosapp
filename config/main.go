/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
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
	Pager        string
	DiskA        string
	DiskB        string
	DiskC        string
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func LoadConfig() Config {
	configHome := os.Getenv("DOSAPP_CONFIG_HOME")
	if configHome == "" {
		if xdg_config_home, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
			configHome = filepath.Join(xdg_config_home, "dosapp")
		} else {
			configHome = filepath.Join(os.Getenv("HOME"), ".config", "dosapp")
		}
	}

	godotenv.Load(filepath.Join(configHome, "dosapp.env"))

	logLevel := getEnv("DOSAPP_LOG_LEVEL", "")
	dosBoxBin := getEnv("DOSAPP_DOSBOX_BIN", "dosbox-x")
	sevenZipBin := getEnv("DOSAPP_7Z_BIN", "7zz")

	dataHome := getEnv("DOSAPP_DATA_HOME", "")
	if dataHome == "" {
		if xdgDataHome, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
			dataHome = filepath.Join(xdgDataHome, "dosapp")
		} else {
			dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share", "dosapp")
		}
	}

	stateHome := os.Getenv("DOSAPP_STATE_HOME")
	if stateHome == "" {
		if xdgStateHome, ok := os.LookupEnv("XDG_STATE_HOME"); ok {
			stateHome = filepath.Join(xdgStateHome, "dosapp")
		} else {
			stateHome = filepath.Join(os.Getenv("HOME"), ".local", "state", "dosapp")
		}
	}

	cacheHome := os.Getenv("DOSAPP_CACHE_HOME")
	if cacheHome == "" {
		if xdg_cache_home, ok := os.LookupEnv("XDG_CACHE_HOME"); ok {
			cacheHome = filepath.Join(xdg_cache_home, "dosapp")
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
	pager := getEnv("PAGER", "cat")
	diskA := getEnv("DOSAPP_DISK_A", filepath.Join(os.Getenv("HOME"), "Documents"))
	diskB := getEnv("DOSAPP_DISK_B", "")
	diskC := getEnv("DOSAPP_DISK_C", filepath.Join(os.Getenv("HOME"), "dosapp", "c"))

	return Config{
		configHome,
		logLevel,
		dosBoxBin,
		sevenZipBin,
		dataHome,
		stateHome,
		cacheHome,
		diskHome,
		linkHome,
		packageHome,
		downloadHome,
		pager,
		diskA,
		diskB,
		diskC,
	}
}

func editConfig(file string) error {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return errors.New("No editor specified.")
	}

	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (conf Config) EditEnvFile() error {
	envFile := filepath.Join(conf.ConfigHome, "dosapp.env")
	return editConfig(envFile)
}

func (conf Config) EnvFileExists() bool {
	envFile := filepath.Join(conf.ConfigHome, "dosapp.env")
	_, err := os.Stat(envFile)
	return err == nil
}
