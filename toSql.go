package gorose

import (
	"errors"
	"fmt"
	"github.com/gohouse/gorose/v3/builder"
	"github.com/gohouse/gorose/v3/parser"
	"reflect"
)

func (db *Database) ToSqlSelect() (sql4prepare string, binds []any) {
	return db.Driver.ToSqlSelect(db.Context)
}

func (db *Database) ToSqlTable() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSqlTable(db.Context)
}
func (db *Database) ToSqlJoin() (sql4prepare string, binds []any, err error) {
	return db.Driver.ToSqlJoin(db.Context)
}

func (db *Database) ToSqlWhere() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSqlWhere(db.Context)
}

func (db *Database) ToSqlOrderBy() (sql4prepare string) {
	return db.Driver.ToSqlOrderBy(db.Context)
}

func (db *Database) ToSqlLimitOffset() (sqlSegment string, binds []any) {
	return db.Driver.ToSqlLimitOffset(db.Context)
}

func (db *Database) ToSql() (sql4prepare string, values []any, err error) {
	return db.Driver.ToSql(db.Context)
}

func (db *Database) ToSqlExists(bind ...any) (sql4prepare string, values []any, err error) {
	if len(bind) > 0 {
		sql4prepare, values, err = db.ToSqlTo(bind[0])
	} else {
		sql4prepare, values, err = db.Driver.ToSql(db.Context)
	}
	if err != nil {
		return
	}
	sql4prepare = fmt.Sprintf("SELECT EXISTS (%s) AS exist", sql4prepare)
	return
}

func (db *Database) ToSqlAggregate(function, column string) (sql4prepare string, values []any, err error) {
	var ctx = *db.Context
	ctx.SelectClause.Columns = append(ctx.SelectClause.Columns, builder.Column{
		Name:  fmt.Sprintf("%s(%s)", function, column),
		Alias: function,
		IsRaw: true,
		Binds: []any{},
	})
	return db.Driver.ToSql(&ctx)
}

func (db *Database) ToSqlTo(obj any, mustColumn ...string) (sql4prepare string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	columns, fieldStruct, _ := parser.StructsParse(obj)
	switch rfv.Kind() {
	case reflect.Struct:
		var data = make(map[string]any)
		data, err = parser.StructDataToMap(rfv, columns, fieldStruct, mustColumn...)
		if err != nil {
			return
		}
		sql4prepare, binds, err = db.Table(obj).Select(columns...).Where(data).Limit(1).ToSql()
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			sql4prepare, binds, err = db.Table(obj).Select(columns...).ToSql()
		}
	default:
		err = errors.New("obj must be struct(slice) or map(slice)")
	}
	return
}

func (db *Database) ToSqlInsert(obj any, args ...builder.TypeToSqlInsertCase) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlInsert(db.Context, obj, args...)
}
func (db *Database) ToSqlDelete(obj any, mustColumn ...string) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlDelete(db.Context, obj, mustColumn...)
}

func (db *Database) ToSqlUpdate(obj any, mustColumn ...string) (sqlSegment string, binds []any, err error) {
	return db.Driver.ToSqlUpdate(db.Context, builder.TypeToSqlUpdateCase{BindOrData: obj, MustColumn: mustColumn})
}

// ToSqlIncDec
//
//	symbol: +/-
//	data: {count: 2}	=> count = count + 2
func (db *Database) ToSqlIncDec(symbol string, data map[string]any) (sql4prepare string, values []any, err error) {
	//return db.Driver.ToSqlIncDec(db.Context, symbol, data)
	return db.Driver.ToSqlUpdate(db.Context, builder.TypeToSqlIncDecCase{Symbol: symbol, Data: data})
}
