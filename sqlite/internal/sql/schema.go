package sql

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	SCHEMA_REGEX = `(?s).+?\(([^\)]+)`
	COLUMN_REGEX = `(?i)(?:\s*(?:(\w+))|(\".+?\"))\s*(NULL|INTEGER|REAL|TEXT|BLOB)`
)

var (
	schemaRegex = regexp.MustCompile(SCHEMA_REGEX)
	columnRegex = regexp.MustCompile(COLUMN_REGEX)
)

const SQLITE_MASTER_SCHEMA = `
    CREATE TABLE sqlite_schema(
    type text,
    name text,
    tbl_name text,
    rootpage integer,
    sql text
  );`

type parsedColumn struct {
	columnName  string
	columnType  string
	columnIndex int
}

func parseSchema(schemaSql string) []parsedColumn {
	matches := schemaRegex.FindStringSubmatch(schemaSql)
	columns := commaSeparatorRegex.Split(matches[1], -1)
	parsedSchema := make([]parsedColumn, len(columns))

	for i, column := range columns {
		matches := columnRegex.FindStringSubmatch(column)
		var columnName string

		if matches[1] == "" {
			columnName = matches[2]
		} else {
			columnName = matches[1]
		}

		parsedSchema[i] = parsedColumn{
			columnName:  columnName,
			columnType:  strings.ToLower(matches[3]),
			columnIndex: i,
		}
	}

	return parsedSchema
}

func getTableSchema(tableName string) string {
	if tableName == "sqlite_schema" || tableName == "sqlite_master" {
		return SQLITE_MASTER_SCHEMA
	} else {
		schemaQuery := fmt.Sprintf("SELECT sql FROM sqlite_schema WHERE name = '%s'", tableName)
		queryResult := ExecuteSelect(schemaQuery)
		return queryResult[:len(queryResult)-1] // trimming the last \n
	}
}
