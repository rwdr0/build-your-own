package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func TestTreeShaFromCommit(t *testing.T) {
	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	if err := os.WriteFile(filepath.Join(tmpDir, "hello.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	utils.RunCmd(t, tmpDir, "git", "add", ".")
	utils.RunCmd(t, tmpDir, "git", "-c", "user.email=t@t", "-c", "user.name=t", "commit", "-m", "init")

	commitHash := utils.RunCmd(t, tmpDir, "git", "rev-parse", "HEAD")
	expectedTree := utils.RunCmd(t, tmpDir, "git", "rev-parse", "HEAD^{tree}")

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	got, err := treeShaFromCommit(commitHash)
	if err != nil {
		t.Fatalf("treeShaFromCommit: %v", err)
	}
	if got != expectedTree {
		t.Errorf("tree sha mismatch: got %s, want %s", got, expectedTree)
	}
}
