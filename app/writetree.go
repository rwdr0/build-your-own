package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

// WriteTree implements the "git write-tree" command
// optionally prints the root tree's SHA-1 hash.
func WriteTree(dir string, printHash bool) []byte {
	body, _, hexHash := writeTree(dir)
	if printHash {
		fmt.Println(hexHash)
	}
	return body
}

// writeTree recursively builds and writes a tree object for rootDirectory.
// It returns the tree body, its SHA-1 hash, and its hex-encoded hash string.
func writeTree(rootDirectory string) ([]byte, [20]byte, string) {
	directoryEntries, err := os.ReadDir(rootDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	treeBody := make([]byte, 0, 5000)

	for _, entry := range directoryEntries {
		var mode []byte
		var objectHash [20]byte
		entryName := []byte(entry.Name())
		entryPath := filepath.Join(rootDirectory, entry.Name())

		if entry.IsDir() {
			if bytes.Equal(entryName, []byte(".git")) {
				continue
			}
			mode = []byte("40000")
			_, objectHash, _ = writeTree(entryPath)
		} else {
			info, _ := entry.Info()

			mode = []byte("100644")
			if info.Mode().Perm()&0o111 != 0 {
				mode = []byte("100755")
			}
			fileContent, err := os.ReadFile(entryPath)
			if err != nil {
				log.Fatal("WriteTree could not read file: ", entryPath)
			}
			objectHash = utils.HashObject(fileContent, false)
		}

		treeBody = append(treeBody, mode...)
		treeBody = append(treeBody, ' ')
		treeBody = append(treeBody, entryName...)
		treeBody = append(treeBody, 0)
		treeBody = append(treeBody, objectHash[:]...)
	}

	header := fmt.Sprintf("tree %d\x00", len(treeBody))
	full := append([]byte(header), treeBody...)
	hash, hexHash := utils.WriteObject(full)

	return treeBody, hash, hexHash
}
