package dialect

import "fmt"

type MsSQLDialect struct {
	placeHolderIndex int
}

func init() {
	Register("mssql", &MsSQLDialect{})
}

func (d *MsSQLDialect) New() IDialect {
	return &MsSQLDialect{}
}

func (d *MsSQLDialect) Placeholder() string {
	//return fmt.Sprintf("$%d", idx)
	d.placeHolderIndex += 1
	return fmt.Sprintf("@p%d", d.placeHolderIndex)
}

func (d *MsSQLDialect) AutoIncrement() string {
	return "IDENTITY(1,1)"
}

func (d *MsSQLDialect) LimitOffset(limit, offset int) string {
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)
}

func (d *MsSQLDialect) QuoteIdentifier(identifier string) string {
	if identifier == "" || identifier == "*" {
		return identifier
	}
	return fmt.Sprintf("[%s]", identifier)
}

func (d *MsSQLDialect) Upsert() string {
	return "MERGE INTO"
}

func (d *MsSQLDialect) LockInShareMode() string { return "" }

//func (d *MsSQLDialect) LockInShareMode() string { return "WITH (HOLDLOCK)" }

func (d *MsSQLDialect) LockForUpdate() string { return "WITH (ROWLOCK)" }
