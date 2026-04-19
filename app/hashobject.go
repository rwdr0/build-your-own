package main

import (
	"log"
	"os"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

// HashObject implements the "git hash-object -w" command, reading a file from
// disk, writing it as a blob object, and printing its SHA-1 hash to stdout.
func HashObject() {
	sourcePath := utils.GetArgumentsForStage(3)[0]
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal("HashObject could not read file")
	}
	utils.HashObject(content, "blob", utils.HashOptions{PrintHash: true, Write: true})
}
