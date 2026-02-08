package catfile

import (
	"bytes"
	"os"
	"testing"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func TestReadObject(t *testing.T) {
	tmpDir := t.TempDir()

	utils.RunCmd(t, tmpDir, "git", "init")

	// Create a file with "hello world" and hash it into the object store.
	filePath := tmpDir + "/testfile"
	if err := os.WriteFile(filePath, []byte("hello world\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash := utils.RunCmd(t, tmpDir, "git", "hash-object", "-w", filePath)

	// chdir into the temp repo so readObject can find .git/objects.
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	data := readObject(hash)

	// Git stores: "blob <size>\0<content>"
	expected := "blob 12\x00hello world\n"
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("got %q, want %q", data, expected)
	}
}
