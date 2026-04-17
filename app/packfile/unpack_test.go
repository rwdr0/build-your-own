package packfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func TestUnpack_WritesNonDeltaObjects(t *testing.T) {
	const testURL = "https://github.com/codecrafters-io/git-sample-1"
	const wantHash = "47b37f1a82bfe85f6d8df52b6258b75e4343b7fd"

	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	packData, err := fetchPackfile(testURL, []string{wantHash})
	if err != nil {
		t.Fatalf("fetchPackfile failed: %v", err)
	}
	if len(packData) == 0 {
		t.Fatal("empty packfile returned")
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	Unpack(packData)

	objPath := filepath.Join(".git", "objects", wantHash[:2], wantHash[2:])
	if _, err := os.Stat(objPath); err != nil {
		t.Errorf("expected non-delta commit %s to be written: %v", wantHash, err)
	}
}
