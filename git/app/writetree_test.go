package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func computeTreeHash(content []byte) string {
	header := fmt.Sprintf("tree %d\x00", len(content))
	full := append([]byte(header), content...)
	hash := sha1.Sum(full)
	return fmt.Sprintf("%x", hash)
}

func TestWriteTree_RegularFile(t *testing.T) {
	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	if err := os.WriteFile(filepath.Join(tmpDir, "hello.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	utils.RunCmd(t, tmpDir, "git", "add", ".")
	expectedHash := utils.RunCmd(t, tmpDir, "git", "write-tree")

	result := WriteTree(tmpDir, false)
	ourHash := computeTreeHash(result)

	if ourHash != expectedHash {
		t.Errorf("hash mismatch: got %s, want %s", ourHash, expectedHash)
	}
}

func TestWriteTree_ExecutableFile(t *testing.T) {
	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	if err := os.WriteFile(filepath.Join(tmpDir, "run.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	utils.RunCmd(t, tmpDir, "git", "add", ".")
	expectedHash := utils.RunCmd(t, tmpDir, "git", "write-tree")

	result := WriteTree(tmpDir, false)
	ourHash := computeTreeHash(result)

	if ourHash != expectedHash {
		t.Errorf("hash mismatch: got %s, want %s", ourHash, expectedHash)
	}
}

func TestWriteTree_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	files := map[string]struct {
		content string
		mode    os.FileMode
	}{
		"a.txt": {"file a\n", 0o644},
		"b.sh":  {"#!/bin/sh\n", 0o755},
		"c.txt": {"file c\n", 0o644},
	}
	for name, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(f.content), f.mode); err != nil {
			t.Fatal(err)
		}
	}

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	utils.RunCmd(t, tmpDir, "git", "add", ".")
	expectedHash := utils.RunCmd(t, tmpDir, "git", "write-tree")

	result := WriteTree(tmpDir, false)
	ourHash := computeTreeHash(result)

	if ourHash != expectedHash {
		t.Errorf("hash mismatch: got %s, want %s", ourHash, expectedHash)
	}
}
