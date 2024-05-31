/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package editor

import (
	"errors"
	"os"
	"os/exec"
)

type Editor struct {
	Bin string
}

func NewEditor(bin string) *Editor {
	return &Editor{
		Bin: bin,
	}
}

func (editor *Editor) Edit(file string) error {
	if editor.Bin == "" {
		return errors.New("No editor specified.")
	}

	cmd := exec.Command(editor.Bin, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
