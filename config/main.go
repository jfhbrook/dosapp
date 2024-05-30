/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
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
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() Config {
	config_home := os.Getenv("DOSAPP_CONFIG_HOME")
	if config_home == "" {
		if xdg_config_home, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
			config_home = filepath.Join(xdg_config_home, "dosapp")
		} else {
			config_home = filepath.Join(os.Getenv("HOME"), ".config", "dosapp")
		}
	}

	godotenv.Load(filepath.Join(config_home, "dosapp.env"))

	logLevel := getEnv("DOSAPP_LOG_LEVEL", "")
	dosBoxBin := getEnv("DOSAPP_DOSBOX_BIN", "dosbox-x")
	sevenZipBin := getEnv("DOSAPP_7Z_BIN", "7zz")

	dataHome := getEnv("DOSAPP_DATA_HOME", "")
	if dataHome == "" {
		if xdg_data_home, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
			dataHome = filepath.Join(xdg_data_home, "dosapp")
		} else {
			dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share", "dosapp")
		}
	}

	stateHome := os.Getenv("DOSAPP_STATE_HOME")
	if stateHome == "" {
		if xdg_state_home, ok := os.LookupEnv("XDG_STATE_HOME"); ok {
			stateHome = filepath.Join(xdg_state_home, "dosapp")
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

	// TODO: Populate and return a config struct
}
