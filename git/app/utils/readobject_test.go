package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestReadObject(t *testing.T) {
	tmpDir := t.TempDir()

	RunCmd(t, tmpDir, "git", "init")

	filePath := tmpDir + "/testfile"
	if err := os.WriteFile(filePath, []byte("hello world\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash := RunCmd(t, tmpDir, "git", "hash-object", "-w", filePath)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	data := ReadObject(hash)

	expected := "blob 12\x00hello world\n"
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("got %q, want %q", data, expected)
	}
}
