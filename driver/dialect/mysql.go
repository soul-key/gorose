package dialect

import "fmt"

type MySQLDialect struct{}

func init() {
	Register("mysql", &MySQLDialect{})
}

func (d *MySQLDialect) New() IDialect {
	return &MySQLDialect{}
}

func (d *MySQLDialect) Placeholder() string {
	return "?"
}

func (d *MySQLDialect) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (d *MySQLDialect) LimitOffset(limit, offset int) string {
	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (d *MySQLDialect) QuoteIdentifier(identifier string) string {
	if identifier == "" || identifier == "*" {
		return identifier
	}
	return fmt.Sprintf("`%s`", identifier)
}

func (d *MySQLDialect) Upsert() string { return "ON DUPLICATE KEY UPDATE" }

func (d *MySQLDialect) LockInShareMode() string { return "LOCK IN SHARE MODE" }
func (d *MySQLDialect) LockForUpdate() string   { return "FOR UPDATE" }
