package utils

import (
	"crypto/sha1"
	"fmt"
)

// HashOptions controls the side effects of HashObject.
type HashOptions struct {
	PrintHash bool
	Write     bool
}

// HashObject formats content as a git object of the given type and returns its
// SHA-1 hash. When opts.Write is true the object is written to the object
// store; when opts.PrintHash is true the hex hash is printed to stdout.
func HashObject(content []byte, objectType string, opts HashOptions) [20]byte {
	formattedContent := fmt.Sprintf("%s %d\x00%s", objectType, len(content), content)

	var hash [20]byte
	var hexHash string
	if opts.Write {
		hash, hexHash = WriteObject([]byte(formattedContent))
	} else {
		hash = sha1.Sum([]byte(formattedContent))
		hexHash = fmt.Sprintf("%x", hash)
	}

	if opts.PrintHash {
		fmt.Println(hexHash)
	}

	return hash
}
