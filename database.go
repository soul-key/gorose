package gorose

import (
	"database/sql"
	"fmt"
	"github.com/gohouse/gorose/v3/builder"
	"github.com/gohouse/gorose/v3/driver"
	"reflect"
	"strings"
)

type Database struct {
	*Engin
	Driver  *driver.Driver
	Context *builder.Context
}

func NewDatabase(g *GoRose) *Database {
	return &Database{
		Driver:  driver.NewDriver(g.driver),
		Engin:   NewEngin(g),
		Context: builder.NewContext(g.prefix),
	}
}
func (db *Database) Table(table any, alias ...string) *Database {
	db.Context.TableClause.Table(table, alias...)
	return db
}

// Select specifies the columns to retrieve.
// Select("a","b")
// Select("a.id as aid","b.id bid")
// Select("id,nickname name")
func (db *Database) Select(columns ...string) *Database {
	db.Context.SelectClause.Select(columns...)
	return db
}

// AddSelect 添加选择列
func (db *Database) AddSelect(columns ...string) *Database {
	db.Context.SelectClause.AddSelect(columns...)
	return db
}

// SelectRaw 允许直接在查询中插入原始SQL片段作为选择列。
func (db *Database) SelectRaw(raw string, binds ...any) *Database {
	db.Context.SelectClause.SelectRaw(raw, binds...)
	return db
}

// Join clause
func (db *Database) Join(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.Join(table, argOrFn...)
	return db
}
func (db *Database) JoinOn(table any, fn func(on builder.IJoinOn)) *Database {
	db.Context.JoinClause.JoinOn(table, fn)
	return db
}

// LeftJoin clause
func (db *Database) LeftJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.LeftJoin(table, argOrFn...)
	return db
}

// RightJoin clause
func (db *Database) RightJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.RightJoin(table, argOrFn...)
	return db
}

// CrossJoin clause
func (db *Database) CrossJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.CrossJoin(table, argOrFn...)
	return db
}
func (db *Database) Where(column any, argsOrclosure ...any) *Database {
	db.Context.WhereClause.Where(column, argsOrclosure...)
	return db
}
func (db *Database) OrWhere(column any, argsOrclosure ...any) *Database {
	db.Context.WhereClause.OrWhere(column, argsOrclosure...)
	return db
}

// WhereRaw 在查询中添加一个原生SQL“where”条件。
//
// sql: 原生SQL条件字符串。
// bindings: SQL绑定参数数组。
func (db *Database) WhereRaw(raw string, bindings ...any) *Database {
	db.Context.WhereClause.WhereRaw(raw, bindings...)
	return db
}
func (db *Database) OrWhereRaw(raw string, bindings ...any) *Database {
	db.Context.WhereClause.OrWhereRaw(raw, bindings...)
	return db
}

// GroupBy 添加 GROUP BY 子句
func (db *Database) GroupBy(columns ...string) *Database {
	db.Context.GroupClause.GroupBy(columns...)
	return db
}
func (db *Database) GroupByRaw(columns ...string) *Database {
	db.Context.GroupClause.GroupByRaw(columns...)
	return db
}

// Having 添加 HAVING 子句, 同where
func (db *Database) Having(column any, argsOrClosure ...any) *Database {
	db.Context.HavingClause.Where(column, argsOrClosure...)
	return db
}
func (db *Database) OrHaving(column any, argsOrClosure ...any) *Database {
	db.Context.HavingClause.OrWhere(column, argsOrClosure...)
	return db
}

// HavingRaw 添加 HAVING 子句, 同where
func (db *Database) HavingRaw(raw string, argsOrClosure ...any) *Database {
	db.Context.HavingClause.WhereRaw(raw, argsOrClosure...)
	return db
}
func (db *Database) OrHavingRaw(raw string, argsOrClosure ...any) *Database {
	db.Context.HavingClause.OrWhereRaw(raw, argsOrClosure...)
	return db
}

// OrderBy adds an ORDER BY clause to the query.
func (db *Database) OrderBy(column string, directions ...string) *Database {
	db.Context.OrderByClause.OrderBy(column, directions...)
	return db
}
func (db *Database) OrderByRaw(column string) *Database {
	db.Context.OrderByClause.OrderByRaw(column)
	return db
}

// Limit 设置查询结果的限制数量。
func (db *Database) Limit(limit int) *Database {
	db.Context.Limit(limit)
	return db
}

// Offset 设置查询结果的偏移量。
func (db *Database) Offset(offset int) *Database {
	db.Context.Offset(offset)
	return db
}

// Page 页数,根据limit确定
func (db *Database) Page(num int) *Database {
	db.Context.Page(num)
	return db
}

// SharedLock 4 select ... locking in share mode
func (db *Database) SharedLock() *Database {
	db.Context.SharedLock()
	return db
}

// LockForUpdate 4 select ... for update
func (db *Database) LockForUpdate() *Database {
	db.Context.LockForUpdate()
	return db
}

func (db *Database) toBind(bind any) (err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.ToSql()
	if err != nil {
		return
	}

	err = db.queryToBindResult(bind, prepare, binds...)
	return
}

// Get 获取查询结果集。
//
// columns: 要获取的列名数组，如果不提供，则获取所有列。
func (db *Database) Get(columns ...string) (res []map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Select(columns...).ToSql()
	if err != nil {
		return
	}

	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db *Database) First(columns ...string) (res map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Select(columns...).Limit(1).ToSql()
	if err != nil {
		return
	}

	res = make(map[string]any)
	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db *Database) Find(id int) (res map[string]any, err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.Where("id", id).Limit(1).ToSql()
	if err != nil {
		return
	}

	res = make(map[string]any)
	err = db.queryToBindResult(&res, prepare, binds...)
	return
}
func (db *Database) queryToBindResult(bind any, query string, args ...any) (err error) {
	return db.Engin.QueryTo(bind, query, args...)
}

func (db *Database) insert(obj any, arg builder.TypeToSqlInsertCase) (res sql.Result, err error) {
	//segment, binds, err := db.ToSqlInsert(obj, ignoreCase, onDuplicateKeys, mustColumn...)
	segment, binds, err := db.ToSqlInsert(obj, arg)
	if err != nil {
		return res, err
	}
	return db.Engin.Exec(segment, binds...)
}
func (db *Database) Insert(obj any, mustColumn ...string) (affectedRows int64, err error) {
	result, err := db.insert(obj, builder.TypeToSqlInsertCase{MustColumn: mustColumn})
	if err != nil {
		return affectedRows, err
	}
	return result.RowsAffected()
}

// InsertGetId 插入数据,获取并自增id
//
// 参考 https://laravel.com/docs/10.x/queries#auto-incrementing-ids
func (db *Database) InsertGetId(obj any, mustColumn ...string) (lastInsertId int64, err error) {
	result, err := db.insert(obj, builder.TypeToSqlInsertCase{MustColumn: mustColumn})
	if err != nil {
		return lastInsertId, err
	}
	return result.LastInsertId()
}

// InsertOrIgnore 插入数据，忽略错误。
//
// 参考 https://laravel.com/docs/10.x/queries#insert-statements
func (db *Database) InsertOrIgnore(obj any, mustColumn ...string) (affectedRows int64, err error) {
	result, err := db.insert(obj, builder.TypeToSqlInsertCase{IsIgnoreCase: true, MustColumn: mustColumn})
	if err != nil {
		return affectedRows, err
	}
	return result.RowsAffected()
}

// Upsert 插入数据，如果存在则更新。
//
// 参考 https://laravel.com/docs/10.x/queries#upserts
// 如果是mysql,则不需要填写第二个参数,MySQL会自动处理唯一索引和主键冲突问题
//
//	eg: Upsert(obj, []string{"id"}, []string{"age"}, "id", "name")
func (db *Database) Upsert(obj any, onDuplicateKeys, updateFields []string, mustColumn ...string) (affectedRows int64, err error) {
	result, err := db.insert(obj, builder.TypeToSqlInsertCase{OnDuplicateKeys: onDuplicateKeys, UpdateFields: updateFields, MustColumn: mustColumn})
	if err != nil {
		return affectedRows, err
	}
	return result.RowsAffected()
}

// Replace 插入数据，如果存在则替换。
//
// 参考 mysql replace into 用法
func (db *Database) Replace(obj any, mustColumn ...string) (affectedRows int64, err error) {
	result, err := db.insert(obj, builder.TypeToSqlInsertCase{IsReplace: true, MustColumn: mustColumn})
	if err != nil {
		return affectedRows, err
	}
	return result.RowsAffected()
}

// UpdateOrInsert 更新数据，如果存在则更新，否则插入。
//
// 参考 https://laravel.com/docs/10.x/queries#update-or-insert
func (db *Database) UpdateOrInsert(conditions, data map[string]any) (affectedRows int64, err error) {
	dbTmp := db.Where(conditions)
	var exists bool
	if exists, err = dbTmp.Exists(); err != nil {
		return
	}
	if exists {
		return dbTmp.Update(data)
	}
	return dbTmp.Insert(data)
}

func (db *Database) Update(obj any, mustColumn ...string) (affectedRows int64, err error) {
	segment, binds, err := db.ToSqlUpdate(obj, mustColumn...)
	if err != nil {
		return affectedRows, err
	}
	return db.Engin.execute(segment, binds...)
}

func (db *Database) Delete(obj any, mustColumn ...string) (affectedRows int64, err error) {
	segment, binds, err := db.ToSqlDelete(obj, mustColumn...)
	if err != nil {
		return affectedRows, err
	}
	return db.Engin.execute(segment, binds...)
}

func (db *Database) incDecEach(symbol string, data map[string]any) (affectedRows int64, err error) {
	prepare, values, err := db.ToSqlIncDec(symbol, data)
	if err != nil {
		return affectedRows, err
	}
	return db.Engin.execute(prepare, values...)
}
func (db *Database) incDec(symbol string, column string, steps ...any) (affectedRows int64, err error) {
	var step any = 1
	if len(steps) > 0 {
		step = steps[0]
	}
	return db.incDecEach(symbol, map[string]any{column: step})
}
func (db *Database) Increment(column string, steps ...any) (affectedRows int64, err error) {
	return db.incDec("+", column, steps...)
}
func (db *Database) Decrement(column string, steps ...any) (affectedRows int64, err error) {
	return db.incDec("-", column, steps...)
}
func (db *Database) IncrementEach(data map[string]any) (affectedRows int64, err error) {
	return db.incDecEach("+", data)
}
func (db *Database) DecrementEach(data map[string]any) (affectedRows int64, err error) {
	return db.incDecEach("-", data)
}

// func (db *Database) Aggregate(functions, columns string) (float64, error) {}
func (db *Database) aggregateSingle(bind any, function, column string) error {
	prepare, values, err := db.ToSqlAggregate(function, column)
	if err != nil {
		return err
	}
	return db.Engin.QueryRow(prepare, values...).Scan(bind)
}
func (db *Database) Max(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "max", column)
	return
}
func (db *Database) Min(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "min", column)
	return
}
func (db *Database) Sum(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "sum", column)
	return
}
func (db *Database) Avg(column string) (res float64, err error) {
	err = db.aggregateSingle(&res, "avg", column)
	return
}
func (db *Database) Count() (res int64, err error) {
	err = db.aggregateSingle(&res, "count", "*")
	return
}

// List 获取指定列的值列表。
func (db *Database) List(column string) (res []any, err error) {
	ress, err := db.Get(column)
	if err != nil {
		return res, err
	}
	for _, v := range ress {
		res = append(res, v[column])
	}
	return
}

// Pluck 从查询结果集中获取键值对列表。
func (db *Database) Pluck(column string, keyColumn string) (res map[any]any, err error) {
	ress, err := db.Get(column, keyColumn)
	if err != nil {
		return res, err
	}
	res = make(map[any]any)
	for _, v := range ress {
		res[v[keyColumn]] = v[column]
	}
	return
}
func (db *Database) Value(column string) (res any, err error) {
	first, err := db.First(column)
	if err != nil {
		return res, err
	}
	return first[column], err
}
func (db *Database) Exists(bind ...any) (b bool, err error) {
	prepare, values, err := db.ToSqlExists(bind...)
	if err != nil {
		return b, err
	}
	err = db.Engin.QueryRow(prepare, values...).Scan(&b)
	return
}
func (db *Database) DoesntExist(bind ...any) (b bool, err error) {
	b, err = db.Exists(bind...)
	return !b, err
}

func (db *Database) Union(b ...builder.IBuilder) *Database {
	db.Context.UnionClause.Union(b...)
	return db
}

func (db *Database) UnionAll(b ...builder.IBuilder) *Database {
	db.Context.UnionClause.UnionAll(b...)
	return db
}

func (db *Database) Truncate(obj ...any) (affectedRows int64, err error) {
	var table string
	var dbTmp = db
	if len(obj) > 0 {
		dbTmp = db.Table(obj[0])
	}
	table, _, err = dbTmp.ToSqlTable()
	if err != nil {
		return
	}
	return db.Engin.execute(fmt.Sprintf("TRUNCATE TABLE %s", table))
}

type TxHandler func() *Database

func (db *Database) Begin() (tx TxHandler, err error) {
	return func() *Database {
		db.Context = builder.NewContext(db.prefix)
		return db
	}, db.Engin.Begin()
}

func (db *Database) Transaction(closure ...func(TxHandler) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, v := range closure {
		err = v(tx)
		if err != nil {
			err2 := db.Rollback()
			if err2 != nil {
				return err2
			}
			return err
		}
	}
	return db.Commit()
}

////////////////////// To 相关的操作 //////////////////////
// 绑定到具体类型上
// Get(),First(),Find() => To()/Bind()
// Value() => ValueTo()
// List()  => ListTo()
// Pluck() => PluckTo()
// Value() => ValueTo()
// Max()   => MaxTo()
// Min()   => MinTo()
// Sum()   => SumTo()

// To 通用查询,go 绑定 struct/map
func (db *Database) To(obj any, mustColumn ...string) (err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.ToSqlTo(obj, mustColumn...)
	if err != nil {
		return
	}

	err = db.queryToBindResult(obj, prepare, binds...)
	return
}

// Bind 查询结果,绑定到结构体
// 与 To 的区别是,绑定字段不作为查询依据
// 经常用在join语句中,手动指定查询字段,然后直接绑定到一个结构体
func (db *Database) Bind(obj any) (err error) {
	var prepare string
	var binds []any
	prepare, binds, err = db.ToSql()
	if err != nil {
		return
	}

	err = db.queryToBindResult(obj, prepare, binds...)
	return
}

// ListTo 获取指定列的值列表。
func (db *Database) ListTo(column string, obj any) (err error) {
	return db.Select(column).toBind(obj)
}

// PluckTo 从查询结果集中获取键值对列表。
func (db *Database) PluckTo(column string, keyColumn string, obj any) (err error) {
	ress, err := db.Get(column, keyColumn)
	if err != nil {
		return err
	}
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	for _, v := range ress {
		rfv2 := reflect.ValueOf(v)
		keys := rfv2.MapKeys()
		key0 := keys[0].String()
		key1 := keys[1].String()
		if strings.HasSuffix(keyColumn, key0) {
			rfv.SetMapIndex(reflect.ValueOf(v[key0]), reflect.ValueOf(v[key1]))
		} else {
			rfv.SetMapIndex(reflect.ValueOf(v[key1]), reflect.ValueOf(v[key0]))
		}
	}
	return
}

// ValueTo 获取指定字段的值,并绑定到给定的变量中
func (db *Database) ValueTo(column string, obj any) (err error) {
	prepare, values, err := db.Select(column).ToSql()
	if err != nil {
		return err
	}
	return db.QueryRow(prepare, values...).Scan(obj)
}

// MaxTo 同 Max
//
//	obj为具体类型的变量,如: var a int, obj 为 &a, 可以得到具体类型
func (db *Database) MaxTo(column string, obj any) (err error) {
	err = db.aggregateSingle(obj, "max", column)
	return
}

// MinTo 同 Min, 参考 MaxTo
func (db *Database) MinTo(column string, obj any) (err error) {
	err = db.aggregateSingle(obj, "min", column)
	return
}

// SumTo 同 Sum, 参考 MaxTo
func (db *Database) SumTo(column string, obj any) (err error) {
	err = db.aggregateSingle(obj, "sum", column)
	return
}
