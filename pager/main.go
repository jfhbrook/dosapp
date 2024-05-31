/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package pager

import (
	"os"
	"os/exec"
)

func Show(file string) error {
	pager := os.Getenv("Pager")

	if pager == "" {
		pager = "cat"
	}

	cmd := exec.Command(pager, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
