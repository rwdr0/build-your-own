// Package utils helpers for unit testing
package utils

import (
	"os/exec"
	"strings"
	"testing"
)

func RunCmd(t *testing.T, dir string, command ...string) string {
	t.Helper()

	name := command[0]
	args := command[1:]
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("%s %v failed: %v", name, args, err)
	}
	return strings.TrimSpace(string(out))
}
