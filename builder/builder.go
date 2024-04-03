package builder

type IBuilder interface {
	ToSql() (sql4prepare string, binds []any, err error)
	ToSqlSelect() (sql4prepare string, binds []any)
	ToSqlTable() (sql4prepare string, values []any, err error)
	ToSqlJoin() (sql4prepare string, binds []any, err error)
	ToSqlWhere() (sql4prepare string, values []any, err error)
	ToSqlOrderBy() (sql4prepare string)
	ToSqlLimitOffset() (sqlSegment string, binds []any)
	ToSqlInsert(obj any, args ...TypeToSqlInsertCase) (sqlSegment string, binds []any, err error)
	ToSqlDelete(obj any, mustColumn ...string) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(obj any, mustColumn ...string) (sqlSegment string, binds []any, err error)
	ToSqlIncDec(symbol string, data map[string]any) (sql4prepare string, values []any, err error)
}

// LimitOffsetClause 存储LIMIT和OFFSET信息。
type LimitOffsetClause struct {
	Limit  int
	Offset int
	Page   int
}

type TypeToSqlUpdateCase struct {
	BindOrData any
	MustColumn []string
}

type TypeToSqlIncDecCase struct {
	Symbol string
	Data   map[string]any
}

type TypeToSqlInsertCase struct {
	IsReplace       bool
	IsIgnoreCase    bool
	OnDuplicateKeys []string
	UpdateFields    []string
	MustColumn      []string
}
