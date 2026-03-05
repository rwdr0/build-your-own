package main

import (
	"bytes"
	"fmt"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func CatFile(objectHash string) {
	data := utils.ReadObject(objectHash)
	_, content, _ := bytes.Cut(data, []byte{0})
	fmt.Print(string(content))
}
