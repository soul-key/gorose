package builder

// TableClause table clause
type TableClause struct {
	Tables any // table name or struct(slice) or subQuery
	Alias  string
}

// Table sets the table name for the query.
func (db *TableClause) Table(table any, alias ...string) {
	var as string
	if len(alias) > 0 {
		as = alias[0]
	}
	db.Tables = table
	db.Alias = as
}
