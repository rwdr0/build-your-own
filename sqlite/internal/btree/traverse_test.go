package btree

import (
	"os"
	"testing"
)

func TestLoadAllLeafTablePages(t *testing.T) {
	dbFile, err := os.Open("../../sample.db")
	if err != nil {
		t.Errorf(`Error Opening db file`)
	}
	defer dbFile.Close()

	testChannel := make(chan LeafTablePage, 1)
	go LoadAllLeafTablePages(0, dbFile, testChannel, true)

	count := uint16(0)
	for c := range testChannel {
		count += c.Header.CellCount
	}

	t.Logf("\n%+v\n", count)
}
