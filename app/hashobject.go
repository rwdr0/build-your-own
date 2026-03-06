package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func HashObject(sourcePath string, printHash bool) [20]byte {
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal("HashObject could not read file")
	}

	formattedContent := fmt.Sprintf("blob %d\x00%s", len(content), content)
	hash := utils.WriteObject([]byte(formattedContent))
	hexHash := fmt.Sprintf("%x", hash)

	if printHash {
		fmt.Println(hexHash)
	}

	return hash
}
