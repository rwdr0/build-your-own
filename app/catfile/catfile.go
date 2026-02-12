// Package catfile => stage #4
package catfile

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func CatFile() {
	objectHash := utils.GetArgumentsForStage(3)[0]
	data := readObject(objectHash)
	_, content, _ := bytes.Cut(data, []byte{0})
	fmt.Print(string(content))
}

func readObject(hash string) []byte {
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
