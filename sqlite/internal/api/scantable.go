package api

import (
	"github.com/rudrowo/sqlite/internal/btree"
	"github.com/rudrowo/sqlite/internal/dataformat"
)

func ScanTable(columnIndicesToSerialize []int, rowLength int, rootPageOffset int64, filter func(row []any) bool, rowsChannel chan<- []any) {
	leafPagesChannel := make(chan btree.LeafTablePage, BTREE_BUFFER_SIZE)
	go btree.LoadAllLeafTablePages(rootPageOffset, dbFile, leafPagesChannel, true)

	for page := range leafPagesChannel {
		for _, cell := range page.Cells { // each cell corresponds to a row
			row := make([]any, rowLength)
			recordBody := cell.Payload.RecordBody
			j, k := 0, 0

			for i, columnType := range cell.Payload.ColumnTypes {
				contentSize := int(dataformat.GetContentSize(columnType))

				//  Lazy serializer
				if j < len(columnIndicesToSerialize) && i == columnIndicesToSerialize[j] {
					var content any

					switch {
					case columnType == 0 && i == 0: // RowId
						content = int64(cell.RowId)
					case columnType == 0: // NULL
						content = nil
					case columnType >= 1 && columnType <= 6: // int
						content = dataformat.DeserializeInteger(recordBody[k : k+contentSize])
					case columnType == 7: // float
						content = dataformat.DeserializeFloat(recordBody[k : k+contentSize])
					default: // string
						content = string(recordBody[k : k+contentSize])
					}

					row[i] = content
					j += 1
				} else {
					row[i] = nil
				}

				k += contentSize
			}
			if filter(row) {
				rowsChannel <- row
			}
		}
	}
	close(rowsChannel)
}
