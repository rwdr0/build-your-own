package api

import (
	btree "github.com/rudrowo/sqlite/internal/btree"
)

func CountRows(rootPageOffset int64) uint16 {
	leafPagesChannel := make(chan btree.LeafTablePage, BTREE_BUFFER_SIZE)
	go btree.LoadAllLeafTablePages(rootPageOffset, dbFile, leafPagesChannel, true)

	cellsCount := uint16(0)
	for page := range leafPagesChannel {
		cellsCount += page.Header.CellCount
	}

	return cellsCount
}
