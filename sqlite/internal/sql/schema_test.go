package sql

import "testing"

func TestParseSchema(t *testing.T) {
	tableStatement := `
  CREATE TABLE users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL,
      email TEXT UNIQUE NOT NULL,
      age INTEGER
  );
  `
	p := parseSchema(tableStatement)
	t.Log(p)
}
