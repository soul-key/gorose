package dialect

import "fmt"

type SQLite3Dialect struct{}

func init() {
	Register("sqlite3", &SQLite3Dialect{})
}

func (d *SQLite3Dialect) New() IDialect {
	return &SQLite3Dialect{}
}

func (d *SQLite3Dialect) Placeholder() string {
	return "?"
}

func (d *SQLite3Dialect) AutoIncrement() string {
	return "AUTOINCREMENT"
}

func (d *SQLite3Dialect) LimitOffset(limit, offset int) string {
	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (d *SQLite3Dialect) QuoteIdentifier(identifier string) string {
	if identifier == "" || identifier == "*" {
		return identifier
	}
	return fmt.Sprintf("\"%s\"", identifier)
}

func (d *SQLite3Dialect) Upsert() string {
	return "INSERT OR REPLACE"
}

func (d *SQLite3Dialect) LockInShareMode() string { return "" }
func (d *SQLite3Dialect) LockForUpdate() string   { return "" }
