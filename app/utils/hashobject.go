package utils

import "fmt"

func HashObject(content []byte, printHash bool) [20]byte {
	formattedContent := fmt.Sprintf("blob %d\x00%s", len(content), content)
	hash, hexHash := WriteObject([]byte(formattedContent))

	if printHash {
		fmt.Println(hexHash)
	}

	return hash
}
