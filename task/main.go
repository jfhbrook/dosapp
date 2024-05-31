package task

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Run(path string, env []string, args ...string) error {
	cmd := exec.Command("task", args...)

	cmd.Dir = filepath.Dir(path)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
