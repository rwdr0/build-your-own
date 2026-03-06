package sql

import (
	"fmt"
	"strconv"

	"github.com/rudrowo/sqlite/internal/btree"
)

const SQLITE_SCHEMA_ROOT_OFFSET = 0

func GetRootPageOFFSET(tableName string) int64 {
	if tableName == "sqlite_schema" || tableName == "sqlite_master" {
		return SQLITE_SCHEMA_ROOT_OFFSET
	} else {
		query := fmt.Sprintf(`SELECT rootpage FROM sqlite_schema WHERE name = '%s'`, tableName)
		rootPageStr := ExecuteSelect(query)
		rootPage, err := strconv.ParseInt(rootPageStr[:len(rootPageStr)-1], 10, 64) // trimming the last \n
		if err != nil {
			panic("Could not parse rootpage in GetRootPageOFFSET")
		}
		return (rootPage - 1) * btree.PAGE_SIZE
	}
}
