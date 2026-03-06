package api

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/rudrowo/sqlite/internal/btree"
)

var dbFile *os.File

const BTREE_BUFFER_SIZE = 100

func Init(fileName string) *os.File {
	var err error
	dbFile, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	btree.PAGE_SIZE = int64(readPageSize())
	return dbFile
}

func readPageSize() uint16 {
	dbHeader := make([]byte, btree.SQLITE3_HEADER_SIZE)
	_, err := dbFile.Read(dbHeader)
	if err != nil {
		log.Fatal(err)
	}

	pageSize := binary.BigEndian.Uint16(dbHeader[16:18])
	return pageSize
}
