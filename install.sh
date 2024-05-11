#!/usr/bin/env bash

mkdir -p ~/.local/bin

if [ -d ~/.local/share/dosapp ]; then
  (cd ~/.local/share/dosapp && git pull)
else
  git clone git@github.com:jfhbrook/dosapp.git ~/.local/share/dosapp
fi

ln -s ~/.local/share/dosapp/bin/dosapp ~/.local/bin/dosapp

if ! command -v dosapp &> /dev/null; then
  echo "Add the following to your PATH, in your shell profile:

      export PATH=\$PATH:~/.local/bin"
fi
