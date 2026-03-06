package sql

import (
	"testing"

	"github.com/rudrowo/sqlite/internal/api"
)

func TestGetRootPageOffset(t *testing.T) {
	dbFile := api.Init("../../sample.db")
	defer dbFile.Close()

	// t.Log(getRootPageOffset("sqlite_schema"))
	t.Log(GetRootPageOFFSET("oranges"))
}
