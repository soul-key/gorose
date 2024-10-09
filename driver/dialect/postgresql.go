package dialect

import (
	"fmt"
)

type PostgresqlDialect struct {
	placeHolderIndex int
}

func init() {
	Register("postgresql", &PostgresqlDialect{})
}

func (d *PostgresqlDialect) New() IDialect {
	return &PostgresqlDialect{}
}

func (d *PostgresqlDialect) Placeholder() string {
	//return fmt.Sprintf("$%d", idx)
	d.placeHolderIndex += 1
	return fmt.Sprintf("$%d", d.placeHolderIndex)
}

//func (d *PostgresqlDialect) IsExpression(obj any) (b bool) {
//	rfv := reflect.Indirect(reflect.ValueOf(obj))
//	if rfv.Kind() == reflect.String && strings.Contains(rfv.String(), "$") {
//		b = true
//	}
//	return
//}
//func (d *PostgresqlDialect) PlaceholderMulti(total int) string {
//	var arr = make([]string, 0, total)
//	for i := 0; i < total; i++ {
//		arr = append(arr, fmt.Sprintf("$%d", i+1))
//	}
//	return strings.Join(arr, ",")
//}
//
//func (d *PostgresqlDialect) PlaceholderMultiInsert(cols, rows int) string {
//	placeholders := make([]string, 0, rows)
//
//	for i := 0; i < rows; i++ {
//		ph := make([]string, cols)
//		for j := 0; j < cols; j++ {
//			ph[j] = fmt.Sprintf("$%d", i*cols+j+1)
//		}
//		placeholders = append(placeholders, "("+strings.Join(ph, ",")+")")
//	}
//	return strings.Join(placeholders, ",")
//}

func (d *PostgresqlDialect) AutoIncrement() string {
	return "SERIAL"
}

func (d *PostgresqlDialect) LimitOffset(limit, offset int) string {
	if offset == 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

func (d *PostgresqlDialect) QuoteIdentifier(identifier string) string {
	if identifier == "" || identifier == "*" {
		return identifier
	}
	return fmt.Sprintf("\"%s\"", identifier)
}

func (d *PostgresqlDialect) Upsert() string {
	return "ON CONFLICT DO UPDATE"
}

func (d *PostgresqlDialect) LockInShareMode() string { return "FOR SHARE" }
func (d *PostgresqlDialect) LockForUpdate() string   { return "FOR UPDATE" }
