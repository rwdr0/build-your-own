// Package utils helpers for unit testing
package utils

import (
	"os/exec"
	"strings"
	"testing"
)

// RunCmd runs an external command in dir and returns its trimmed stdout.
// It calls t.Fatal if the command exits with a non-zero status.
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
