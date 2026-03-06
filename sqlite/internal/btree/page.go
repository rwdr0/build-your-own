/*
Based on
https://www.sqlite.org/fileformat.html#record_format
https://saveriomiroddi.github.io/SQLIte-database-file-format-diagrams/
*/
package btree

import (
	"encoding/binary"

	"github.com/rudrowo/sqlite/internal/dataformat"
)

var PAGE_SIZE = int64(4096)

const (
	INTERIOR_INDEX_PAGE_TYPE = 0x02
	LEAF_INDEX_PAGE_TYPE     = 0x0a

	INTERIOR_TABLE_PAGE_TYPE = 0x05
	LEAF_TABLE_PAGE_TYPE     = 0x0d

	SQLITE3_HEADER_SIZE = 100
)

// Headers
type (
	leafHeader struct {
		PageType  uint8
		CellCount uint16
	}
	interiorHeader struct {
		pageType         uint8
		cellCount        uint16
		rightmostPointer uint32
	}
	recordHeader struct {
		ColumnTypes []uint64 // []varint
		HeaderSize  uint64   // varint
	}
)

// Cells
type (
	interiorTableCell struct {
		leftChildPointer uint32
		rowId            uint64 // varint
	}
	leafTableCell struct {
		Payload struct {
			RecordBody []byte
			recordHeader
		}
		RowId uint64 // varint
	}
)

// Pages
type (
	interiorTablePage struct {
		cellPointers []uint16
		cells        []interiorTableCell
		header       interiorHeader
	}
	LeafTablePage struct {
		CellPointers []uint16
		Cells        []leafTableCell
		Header       leafHeader
	}
)

func (page *interiorTablePage) loadFromBuffer(fileBuffer []byte, isSchemaPage bool) {
	var bi int // bi is the buffer index
	if isSchemaPage {
		bi = SQLITE3_HEADER_SIZE
	} else {
		bi = 0
	}

	page.header.pageType = fileBuffer[bi]
	bi += 3
	page.header.cellCount = binary.BigEndian.Uint16(fileBuffer[bi : bi+2])
	bi += 5
	page.header.rightmostPointer = binary.BigEndian.Uint32(fileBuffer[bi : bi+4])
	bi += 4

	page.cellPointers = make([]uint16, page.header.cellCount)
	for i := range page.cellPointers {
		page.cellPointers[i] = binary.BigEndian.Uint16(fileBuffer[bi : bi+2])
		bi += 2
	}

	page.cells = make([]interiorTableCell, page.header.cellCount)
	for i := range page.cells {
		bi := page.cellPointers[i]
		cell := &(page.cells[i])

		cell.leftChildPointer = binary.BigEndian.Uint32(fileBuffer[bi : bi+4])
		bi += 4
		cell.rowId, _ = dataformat.DeserializeVarint(fileBuffer[bi:])
	}
}

func (page *LeafTablePage) loadFromBuffer(fileBuffer []byte, isSchemaPage bool) {
	var bi int // bi is the buffer index
	if isSchemaPage {
		bi = SQLITE3_HEADER_SIZE
	} else {
		bi = 0
	}

	page.Header.PageType = fileBuffer[bi]
	bi += 3
	page.Header.CellCount = binary.BigEndian.Uint16(fileBuffer[bi : bi+2])
	bi += 5

	page.CellPointers = make([]uint16, page.Header.CellCount)
	for i := range page.CellPointers {
		page.CellPointers[i] = binary.BigEndian.Uint16(fileBuffer[bi : bi+2])
		bi += 2
	}

	page.Cells = make([]leafTableCell, page.Header.CellCount)
	for i := range page.Cells {
		bi := page.CellPointers[i]
		cell := &(page.Cells[i])

		payloadSize, bytesRead := dataformat.DeserializeVarint(fileBuffer[bi:])
		bi += bytesRead
		cell.RowId, bytesRead = dataformat.DeserializeVarint(fileBuffer[bi:])
		bi += bytesRead

		payload := fileBuffer[bi : bi+uint16(payloadSize)]
		headerSize, bytesRead := dataformat.DeserializeVarint(payload)
		bi = bytesRead // bi is now index of the payload buffer

		columnTypes := make([]uint64, headerSize)
		columnTypesRead := 0
		for bi < uint16(headerSize) {
			columnTypes[columnTypesRead], bytesRead = dataformat.DeserializeVarint(payload[bi:])
			bi += bytesRead
			columnTypesRead += 1
		}
		cell.Payload.ColumnTypes = make([]uint64, columnTypesRead)
		copy(cell.Payload.ColumnTypes, columnTypes[0:columnTypesRead]) // Discard excess allocated memory

		cell.Payload.RecordBody = payload[bi:]
		cell.Payload.HeaderSize = headerSize
	}
}
