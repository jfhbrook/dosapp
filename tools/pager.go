/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package tools

import (
	"os"
	"os/exec"
)

type Pager struct {
	Bin string
}

func NewPager(bin string) *Pager {
	return &Pager{
		Bin: bin,
	}
}

func (pager *Pager) Show(file string) error {
	bin := pager.Bin
	// TODO: This doesn't seem to be respecting dotenv. Pass the config in here
	// explicitly, instead of depending on the environment.
	if bin == "" {
		bin = "cat"
	}

	cmd := exec.Command(bin, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
