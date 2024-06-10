# dosapp

Manage DOSBox apps.

## Install

First, make sure you have all the dependencies installed:

- [DosBox-X](https://dosbox-x.com/) (or DOSBox with configuration)
- [go-task](https://taskfile.dev)
- `bash` - `dosapp` executes tasks using a shell
- `curl`
- `7zz` - Packages currently only work with .7z files, but having `unzip` and
  `tar` are probably good ideas as well

Then, install `dosapp`:

```sh
go install github.com/jfhbrook/dosapp@latest
```

## Getting Started

To get started, run:

```sh
dosapp config
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

### Remove the App

To remove an app from your install, clean up any links and remove the
configuration, run:

```sh
dosapp remove [APP_NAME]
```

## Packages

Here's a list of packages:

- `wordperfect`
- `turbo-pascal`
- **TODO**: `dostodon`
- **TODO**: `freebasic`
- **TODO**: [As-Easy-As](http://www.triusinc.com/forums/viewtopic.php?t=10)
- **TODO**: [VGAPaint](https://www.bttr-software.de/products/vp386/)

## Development

This project contains a `Taskfile`, with two tasks:

- `task check` - Run `shellcheck` against the shell scripts in the project.
- `task install` - Symlink `./bin/dosapp` to `~/.local/bin/dosapp`.

## License

I'm releasing this under an MIT license. See [LICENSE](./LICENSE) for details.
