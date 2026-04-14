package utils

import "fmt"

// HashObject formats content as a git blob object, writes it to the object
// store, and returns its SHA-1 hash. If printHash is true the hex hash is
// printed to stdout.
func HashObject(content []byte, printHash bool) [20]byte {
	formattedContent := fmt.Sprintf("blob %d\x00%s", len(content), content)
	hash, hexHash := WriteObject([]byte(formattedContent))

	if printHash {
		fmt.Println(hexHash)
	}

	return hash
}
