package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rwdr0/build-your-own/git/app/packfile"
	"github.com/rwdr0/build-your-own/git/app/utils"
)

// Clone implements the "git clone" command, fetching objects from the given
// remote URL into the given target directory using the git smart HTTP protocol.
func Clone() {
	args := utils.GetArgumentsForStage(2, 3)
	if len(args) < 2 {
		log.Fatalf("usage: clone <url> <directory>")
	}
	url, dir := args[0], args[1]

	for _, sub := range []string{".git/objects", ".git/refs"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			log.Fatalf("mkdir %s: %v", sub, err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, ".git/HEAD"), []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		log.Fatalf("write HEAD: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		log.Fatalf("chdir %s: %v", dir, err)
	}

	refs, err := packfile.FetchRefs(url)
	if err != nil {
		log.Fatalf("Failed to fetch packfile refs %v", err)
	}

	packFile, err := packfile.FetchPackfile(url, refs)
	if err != nil {
		log.Fatalf("Failed to fetch packfile %v", err)
	}

	packfile.Unpack(packFile)
	headRef := refs[0]
	treeSha, err := treeShaFromCommit(headRef)
	if err != nil {
		log.Fatalf("Error retrieving treeSha from commit object %v", err)
	}
	checkout(treeSha, ".")
}

// checkout materializes the tree identified by treeSha into dir, recursing
// into subtrees. dir must already exist.
func checkout(treeSha, dir string) {
	data := utils.ReadObject(treeSha)
	_, body, ok := bytes.Cut(data, []byte{0})
	if !ok {
		log.Fatalf("tree %s: missing header terminator", treeSha)
	}

	for len(body) > 0 {
		sp := bytes.IndexByte(body, ' ')
		mode := string(body[:sp])
		body = body[sp+1:]

		nul := bytes.IndexByte(body, 0)
		name := string(body[:nul])
		body = body[nul+1:]

		sha := hex.EncodeToString(body[:20])
		body = body[20:]

		path := filepath.Join(dir, name)
		if mode == "40000" {
			if err := os.MkdirAll(path, 0o755); err != nil {
				log.Fatalf("mkdir %s: %v", path, err)
			}
			checkout(sha, path)
			continue
		}

		blob := utils.ReadObject(sha)
		_, content, ok := bytes.Cut(blob, []byte{0})
		if !ok {
			log.Fatalf("blob %s: missing header terminator", sha)
		}

		perm := os.FileMode(0o644)
		if mode == "100755" {
			perm = 0o755
		}
		if err := os.WriteFile(path, content, perm); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
	}
}

// treeShaFromCommit reads the loose commit object at the given hash and returns
// the tree sha it points to.
func treeShaFromCommit(commitHash string) (string, error) {
	data := utils.ReadObject(commitHash)

	_, body, ok := bytes.Cut(data, []byte{0})
	if !ok {
		return "", fmt.Errorf("commit %s: missing header terminator", commitHash)
	}

	const prefix = "tree "
	if !bytes.HasPrefix(body, []byte(prefix)) {
		return "", fmt.Errorf("commit %s: missing tree line", commitHash)
	}
	line := body[len(prefix):]
	nl := bytes.IndexByte(line, '\n')
	if nl < 0 {
		return "", fmt.Errorf("commit %s: malformed tree line", commitHash)
	}
	return string(line[:nl]), nil
}
