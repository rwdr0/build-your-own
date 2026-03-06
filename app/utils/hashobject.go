package utils

import (
	"fmt"
	"log"
	"os"
)

func HashObject(sourcePath string, printHash bool) [20]byte {
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal("HashObject could not read file")
	}

	formattedContent := fmt.Sprintf("blob %d\x00%s", len(content), content)
	hash, hexHash := WriteObject([]byte(formattedContent))

	if printHash {
		fmt.Println(hexHash)
	}

	return hash
}
