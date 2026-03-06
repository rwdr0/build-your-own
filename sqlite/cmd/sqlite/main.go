package main

import (
	"fmt"
	"os"

	"github.com/rudrowo/sqlite/internal/api"
	"github.com/rudrowo/sqlite/internal/btree"
	"github.com/rudrowo/sqlite/internal/sql"
)

func main() {
	fileName := os.Args[1]
	userCommand := os.Args[2]

	dbFile := api.Init(fileName)
	defer dbFile.Close()

	switch userCommand {
	case ".dbinfo":
		fmt.Printf("database page size: %v\n", btree.PAGE_SIZE)
		fmt.Printf("number of tables: %v", sql.ExecuteSelect("SELECT COUNT(*) FROM sqlite_schema"))
	case ".tables":
		fmt.Print(sql.ExecuteSelect("SELECT tbl_name FROM sqlite_schema WHERE tbl_name != 'sqlite_sequence'"))
	default:
		fmt.Print(sql.ExecuteSelect(userCommand))
	}
}
