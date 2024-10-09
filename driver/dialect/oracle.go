package dialect

import "fmt"

type OracleDialect struct {
	placeHolderIndex int
}

func init() {
	Register("oracle", &OracleDialect{})
}

func (d *OracleDialect) New() IDialect {
	return &OracleDialect{}
}

func (d *OracleDialect) Placeholder() string {
	//return fmt.Sprintf("$%d", idx)
	d.placeHolderIndex += 1
	return fmt.Sprintf("@p%d", d.placeHolderIndex)
}

func (d *OracleDialect) AutoIncrement() string {
	return "" // 需要使用序列
}

func (d *OracleDialect) LimitOffset(limit, offset int) string {
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
}

func (d *OracleDialect) QuoteIdentifier(identifier string) string {
	if identifier == "" || identifier == "*" {
		return identifier
	}
	return fmt.Sprintf("\"%s\"", identifier)
}

func (d *OracleDialect) Upsert() string {
	return "MERGE INTO"
}

func (d *OracleDialect) LockInShareMode() string { return "" }
func (d *OracleDialect) LockForUpdate() string   { return "FOR UPDATE" }
