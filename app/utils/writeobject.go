package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

// WriteObject compresses content with zlib, writes it to .git/objects under
// its SHA-1 hash, and returns the raw hash and its hex-encoded string.
func WriteObject(content []byte) ([20]byte, string) {
	hash := sha1.Sum(content)
	hexHash := fmt.Sprintf("%x", hash)
	destinationPath := fmt.Sprintf(".git/objects/%s/%s", hexHash[:2], hexHash[2:])

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write([]byte(content))
	w.Close()

	os.MkdirAll(filepath.Dir(destinationPath), 0o755)
	os.WriteFile(destinationPath, buf.Bytes(), 0o644)

	return hash, hexHash
}
