#!/usr/bin/env bash

mkdir -p ~/.local/bin

if [ -d ~/.local/share/dosapp ]; then
  (cd ~/.local/share/dosapp && git pull)
else
  git clone git@github.com:jfhbrook/dosapp.git ~/.local/share/dosapp
fi

(cd ~/.local/share/dosapp && task install)

if ! command -v dosapp &> /dev/null; then
  echo "Add the following to your PATH, in your shell profile:

      export PATH=\$PATH:~/.local/bin"
fi
