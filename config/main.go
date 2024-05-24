/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// TODO: Create a struct for the config

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() {
	config_home := os.Getenv("DOSAPP_CONFIG_HOME")
	if config_home == "" {
		if xdg_config_home, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
			config_home = filepath.Join(xdg_config_home, "dosapp")
		} else {
			config_home = filepath.Join(os.Getenv("HOME"), ".config", "dosapp")
		}
	}

	godotenv.Load(filepath.Join(config_home, "dosapp.env"))

	// TODO: Populate and return a config struct
}

/*
function load-config {
  export DOSAPP_ROOT="${DOSAPP_ROOT:-$(dirname "$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")")}"
  export DOSAPP_CONFIG_HOME="${DOSAPP_CONFIG_HOME:-${XDG_CONFIG_HOME:-${HOME}/.config}/dosapp}"

  if [ -f "${DOSAPP_CONFIG_HOME}/dosapp.env" ]; then
    # shellcheck disable=SC1091
    source "${DOSAPP_CONFIG_HOME}/dosapp.env"
  fi

  export DEBUG="${DEBUG:-}"
  export DOSAPP_DOSBOX_BIN="${DOSAPP_DOSBOX_BIN:-dosbox-x}"
  export DOSAPP_7Z_BIN="${DOSAPP_7Z_BIN:-7zz}"
  export DOSAPP_DATA_HOME="${DOSAPP_DATA_HOME:-${XDG_DATA_HOME:-${HOME}/.local/share}/dosapp}"
  export DOSAPP_STATE_HOME="${DOSAPP_STATE_HOME:-${XDG_STATE_HOME:-${HOME}/.local/state}/dosapp}"
  export DOSAPP_CACHE_HOME="${DOSAPP_CACHE_HOME:-${XDG_CACHE_HOME:-${HOME}/.cache}/dosapp}"
  export DOSAPP_DISK_HOME="${DOSAPP_DISK_HOME:-${HOME}/dosapp}"
  export DOSAPP_LINK_HOME="${DOSAPP_LINK_HOME:-${HOME}/.local/bin}"
  export DOSAPP_PACKAGE_HOME="${DOSAPP_PACKAGE_HOME:-${DOSAPP_ROOT}/packages}"
  export DOSAPP_DOWNLOAD_HOME="${DOSAPP_DOWNLOAD_HOME:-${DOSAPP_CACHE_HOME}/downloads}"
  export PAGER="${PAGER:-cat}"
  export DOSAPP_DISK_A="${DOSAPP_DISK_A:-${HOME}/Documents}"
  export DOSAPP_DISK_B="${DOSAPP_DISK_B:-}"
  export DOSAPP_DISK_C="${DOSAPP_DISK_C:-${HOME}/dosapp/c}"
}
*/
