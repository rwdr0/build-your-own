package sql

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/rudrowo/sqlite/internal/api"
	"github.com/rudrowo/sqlite/internal/dataformat"
)

const (
	SELECT_STATEMENT_REGEX = `(?i)^SELECT\s+(.*?)\s+FROM\s+(\w+)\s*(?:\s+WHERE\s+(.*))?$`
	WHERE_CLAUSE_REGEX     = `([a-zA-Z_][a-zA-Z0-9_]*)\s*(=|!=|<=|>=|<|>)\s*'?([^']+)'?\s*`
	COUNT_REGEX            = `(?i)COUNT\(\*\)`
)

var (
	whereClauseRegex      = regexp.MustCompile(WHERE_CLAUSE_REGEX)
	selectStatementRegex  = regexp.MustCompile(SELECT_STATEMENT_REGEX)
	commaSeparatorRegex   = regexp.MustCompile(`\s*,\s*`)
	countExpresseionRegex = regexp.MustCompile(COUNT_REGEX)
)

func ExecuteSelect(query string) string {
	matches := selectStatementRegex.FindStringSubmatch(query)

	columnTokens := commaSeparatorRegex.Split(matches[1], -1)
	tableName := matches[2]
	whereClause := matches[3]
	rootPageOffset := GetRootPageOFFSET(tableName)

	if countExpresseionRegex.MatchString(matches[1]) {
		rowCount := api.CountRows(rootPageOffset)
		return strconv.FormatInt(int64(rowCount), 10) + "\n"
	}

	schemaSql := getTableSchema(tableName)
	parsedSchema := parseSchema(schemaSql)
	selectedColumnIndices := make([]int, len(columnTokens))

	for i, columnName := range columnTokens {
		for _, parsedColumn := range parsedSchema {
			if columnName == parsedColumn.columnName {
				selectedColumnIndices[i] = parsedColumn.columnIndex
			}
		}
	}

	filter, lhsIndex := parseWhereClause(whereClause, parsedSchema)

	columnIndicesToSerialize := make([]int, len(selectedColumnIndices))
	copy(columnIndicesToSerialize, selectedColumnIndices)
	if lhsIndex >= 0 && !slices.Contains(columnIndicesToSerialize, lhsIndex) {
		columnIndicesToSerialize = append(columnIndicesToSerialize, lhsIndex)
	}
	slices.Sort(columnIndicesToSerialize)

	rowsChannel := make(chan []any)
	go api.ScanTable(columnIndicesToSerialize, len(parsedSchema), rootPageOffset, filter, rowsChannel)

	var result strings.Builder

	for row := range rowsChannel {
		firstPrintInRow := true

		for _, si := range selectedColumnIndices {
			if !firstPrintInRow {
				result.WriteByte('|')
			}

			content := row[si]
			if content != nil {
				switch c := row[si].(type) {
				case int64:
					result.WriteString(strconv.FormatInt(c, 10))
				case float64:
					result.WriteString(strconv.FormatFloat(c, 'f', 2, 64))
				case string: // string
					result.WriteString(c)
				}
			}

			firstPrintInRow = false
		}
		result.WriteByte('\n')
	}
	return result.String()
}

func parseWhereClause(whereClause string, parsedSchema []parsedColumn) (func(row []any) bool, int) {
	if whereClause == "" {
		return func(row []any) bool {
			return true
		}, -1
	}

	tokens := whereClauseRegex.FindStringSubmatch(whereClause)
	lhsToken, operatorToken, rhsToken := tokens[1], tokens[2], tokens[3]

	var intPrimitive func(int64, int64) bool
	var floatPrimitive func(float64, float64) bool
	var stringPrimitive func(string, string) bool

	switch operatorToken {
	case "=":
		intPrimitive = equalToPrimitive[int64]
		floatPrimitive = equalToPrimitive[float64]
		stringPrimitive = equalToPrimitive[string]
	case "!=":
		intPrimitive = notEqualToPrimitive[int64]
		floatPrimitive = notEqualToPrimitive[float64]
		stringPrimitive = notEqualToPrimitive[string]
	case ">":
		intPrimitive = strictlyGreaterThanPrimitive[int64]
		floatPrimitive = strictlyGreaterThanPrimitive[float64]
		stringPrimitive = strictlyGreaterThanPrimitive[string]
	case ">=":
		intPrimitive = greaterThanOrEqualToPrimitive[int64]
		floatPrimitive = greaterThanOrEqualToPrimitive[float64]
		stringPrimitive = greaterThanOrEqualToPrimitive[string]
	case "<":
		intPrimitive = strictlyLessThanPrimitive[int64]
		floatPrimitive = strictlyLessThanPrimitive[float64]
		stringPrimitive = strictlyLessThanPrimitive[string]
	case "<=":
		intPrimitive = lessThanOrEqualToPrimitive[int64]
		floatPrimitive = lessThanOrEqualToPrimitive[float64]
		stringPrimitive = lessThanOrEqualToPrimitive[string]
	}

	var targetColumn parsedColumn

	for _, parsedColumn := range parsedSchema {
		if lhsToken == parsedColumn.columnName {
			targetColumn = parsedColumn
			break
		}
	}

	lhsIndex := targetColumn.columnIndex
	switch targetColumn.columnType {
	case "integer":
		rhsArg, err := strconv.ParseInt(rhsToken, 10, 64)
		if err != nil {
			panic("Type conversion failed for deserialized int64 value")
		}

		return func(row []any) bool {
			if row[lhsIndex] == nil {
				return false
			}
			return intPrimitive(row[lhsIndex].(int64), rhsArg)
		}, lhsIndex
	case "real":
		rhsArg, err := strconv.ParseFloat(rhsToken, 64)
		if err != nil {
			panic("Type conversion failed for deserialized float64 value")
		}

		return func(row []any) bool {
			if row[lhsIndex] == nil {
				return false
			}
			return floatPrimitive(row[lhsIndex].(float64), rhsArg)
		}, lhsIndex
	case "text":
		rhsArg := rhsToken

		return func(row []any) bool {
			if row[lhsIndex] == nil {
				return false
			}
			return stringPrimitive(row[lhsIndex].(string), rhsArg)
		}, lhsIndex
	default:
		panic("Malformed schema passed to where clause")
	}
}

// Comparison Primitives
type Constraint dataformat.DeserializedTypes

func equalToPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs == rhs
}

func notEqualToPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs != rhs
}

func strictlyGreaterThanPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs > rhs
}

func greaterThanOrEqualToPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs >= rhs
}

func strictlyLessThanPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs < rhs
}

func lessThanOrEqualToPrimitive[T Constraint](lhs T, rhs T) bool {
	return lhs <= rhs
}
