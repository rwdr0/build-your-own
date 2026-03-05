package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func LsTree(objectHash string) {
	objectBinary := utils.ReadObject(objectHash)
	zeroByteIdx := bytes.IndexByte(objectBinary, 0)
	content := objectBinary[zeroByteIdx+1:]

	names := make([]string, 0)
	var name strings.Builder

	for _, byte := range content {
		switch byte {
		case ' ':
			name.Reset()
		case 0:
			names = append(names, name.String())
		default:
			name.WriteByte(byte)
		}
	}

	for _, name := range names {
		fmt.Println(name)
	}
}
