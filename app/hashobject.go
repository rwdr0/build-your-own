package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func HashObject() {
	inputFileName := utils.GetArgumentsForStage(3)[0]

	content, err := os.ReadFile(inputFileName)
	if err != nil {
		log.Fatal("HashObject could not read file")
	}

	formattedContent := fmt.Sprintf("blob %d\x00%s", len(content), content)

	hash := sha1.Sum([]byte(formattedContent))
	hexHash := fmt.Sprintf("%x", hash)
	path := fmt.Sprintf(".git/objects/%s/%s", hexHash[:2], hexHash[2:])

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write([]byte(formattedContent))
	w.Close()

	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, buf.Bytes(), 0o644)

	fmt.Println(hexHash)
}
