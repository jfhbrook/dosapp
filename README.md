# dosapp

Manage DOSBox apps.

## Install

Installing `dosapp` is currently a matter of cloning this repo and linking
`./bin/dosapp` to a location in your `$PATH`. To YOLO that install, you can
try running:

```sh
bash <(curl -sSfL https://raw.githubusercontent.com/jfhbrook/dosapp/main/install.sh)
```

This script should also update `dosapp` if it's already installed.

## Getting Started

To get started, run:

```sh
dosapp init
```

This will generate the initial configuration.

### Environment Variables

`dosapp` is configured through dotenv files and environment variables. The
following environment variables are respected:

#### Mount Locations

- `DOSAPP_DISK_HOME`: The base location for your dosapp disks. Some apps may
  create custom disks under this location.
- `DOSAPP_DISK_A`: The standard location for the A drive.
- `DOSAPP_DISK_B`: The standard location for the B drive.
- `DOSAPP_DISK_C`: The standard location for the C drive. Many apps will
  require this mount to function.

#### External Binaries

- `DOSAPP_DOSBOX_BIN`: The location of the DOSBox binary. Defaults to
  `dosbox-x`, which has more features than DOSBox. However, this may be
  overwritten in cases where vanilla DOSBox is more reliable, or is what's
  installed.
- `DOSAPP_7Z_BIN`: The location of the 7z/7zz binary. The default is `7zz`.
- `EDITOR`: The editor to use for editing dotenv files. This is commonly set
  to `vi`.
- `PAGER`: The pager to use for viewing READMEs. Defaults to `cat`, but is
  commonly set to `less`.

#### Runtime Directory Locations

These directories are typically configured to match the XDG spec, but may be
overwritten.

- `DOSAPP_DATA_HOME`: The location to store app data. Defaults to
  `~/.local/share/dosapp`.
- `DOSAPP_STATE_HOME`: The location to store non-critical app state. Defaults
  to `~/.local/state/dosapp`.
- `DOSAPP_CACHE_HOME`: The location to store cached app data. Defaults to
  `~/.cache/dosapp`.
- `DOSAPP_LINK_HOME`: The location to place bin scripts for apps. Defaults to
  `~/.local/bin`.
- `DOSAPP_DOWNLOAD_HOME`: The location to store downloads. Defaults to
  `~/.cache/dosapp/downloads`, and respects `DOSAPP_CACHE_HOME`.

#### Other

- `DEBUG`: Set to `1` to enable debug logging.

## Usage

### Start DOSBox

To start DOSBox with the standard configuration and mounts, run:

```sh
dosapp
```

### Install an App

To install an app, run:

```sh
dosapp install [APP_NAME]
```

### Start an App

To start an installed app with app-specific configuration and mounts, run:

```sh
dosapp start [APP_NAME]
```

### Create a Bin Script for an App

To create scripts for your installed app, run `dosapp link [APP_NAME]`.
For example, to create a `wp` script that launches WordPerfect:

```sh
dosapp link wordperfect
```

## Packages

Here's a list of packages:

- `wordperfect`
- **TODO**: `turbo-pascal`
- **TODO**: `dostodon`
