package utils

import (
	"compress/zlib"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadObject(hash string) []byte {
	path := filepath.Join(".git", "objects", hash[:2], hash[2:])

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open object file: %v", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		log.Fatalf("zlib reader: %v", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("read object: %v", err)
	}

	return data
}
