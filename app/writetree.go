package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func WriteTree() {
	writeTree(".", true) // internal recursive implementation
}

func writeTree(rootDirectory string, printHash bool) [20]byte {
	directoryEntries, err := os.ReadDir(rootDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	treeObject := make([]byte, 0, 5000)

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
			objectHash = writeTree(entryPath, false)
		} else {
			info, _ := entry.Info()

			mode = []byte("100644")
			if info.Mode().Perm()&0o111 != 0 {
				mode = []byte("100755")
			}
			objectHash = utils.HashObject(entryPath, false)
		}

		treeObject = append(treeObject, mode...)
		treeObject = append(treeObject, ' ')
		treeObject = append(treeObject, entryName...)
		treeObject = append(treeObject, 0)
		treeObject = append(treeObject, objectHash[:]...) // The slice operator [:] converts a fixed-size array to a slice.
	}

	header := fmt.Sprintf("tree %d\x00", len(treeObject))
	treeObject = append([]byte(header), treeObject...)

	hash, hexHash := utils.WriteObject(treeObject)

	if printHash {
		fmt.Println(hexHash)
	}

	return hash
}
