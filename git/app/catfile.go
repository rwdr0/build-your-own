package main

import (
	"bytes"
	"fmt"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

// CatFile implements the "git cat-file -p" command, reading a git object by its
// hash and printing its content (without the header) to stdout.
func CatFile() {
	objectHash := utils.GetArgumentsForStage(3)[0]
	data := utils.ReadObject(objectHash)
	_, content, _ := bytes.Cut(data, []byte{0})
	fmt.Print(string(content))
}
