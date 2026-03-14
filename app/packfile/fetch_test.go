package packfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func TestFetch_VerifyPack(t *testing.T) {
	const testURL = "https://github.com/codecrafters-io/git-sample-1"
	refs := []string{"47b37f1a82bfe85f6d8df52b6258b75e4343b7fd"}

	tmpDir := t.TempDir()
	utils.RunCmd(t, tmpDir, "git", "init")

	packData, err := fetchPackfile(testURL, refs)
	if err != nil {
		t.Fatalf("fetchPackfile failed: %v", err)
	}
	if len(packData) == 0 {
		t.Fatal("empty packfile returned")
	}

	packDir := filepath.Join(tmpDir, ".git", "objects", "pack")
	if err := os.MkdirAll(packDir, 0o755); err != nil {
		t.Fatal(err)
	}

	packPath := filepath.Join(packDir, "test.pack")
	if err := os.WriteFile(packPath, packData, 0o644); err != nil {
		t.Fatal(err)
	}

	utils.RunCmd(t, tmpDir, "git", "index-pack", packPath)

	output := utils.RunCmd(t, tmpDir, "git", "verify-pack", "-v", packPath)

	hasObjects := strings.Contains(output, "commit") ||
		strings.Contains(output, "blob") ||
		strings.Contains(output, "tree")
	if !hasObjects {
		t.Errorf("verify-pack output contains no known object types:\n%s", output)
	}
}
