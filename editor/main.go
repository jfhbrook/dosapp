/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package editor

import (
	"errors"
	"os"
	"os/exec"
)

func Edit(file string) error {
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
