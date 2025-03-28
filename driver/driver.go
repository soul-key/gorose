package driver

import (
	"errors"
	"fmt"
	"github.com/gohouse/gorose/v3/builder"
	"github.com/gohouse/gorose/v3/driver/dialect"
	"github.com/gohouse/gorose/v3/parser"
	"regexp"
	"sort"

	"reflect"
	"strings"
)

type IDriver interface {
	ToSqlSelect(c *builder.Context) (sql4prepare string, binds []any)
	ToSqlTable(c *builder.Context) (sql4prepare string, values []any, err error)
	ToSqlJoin(c *builder.Context) (sql4prepare string, binds []any, err error)
	ToSqlWhere(c *builder.Context) (sql4prepare string, values []any, err error)
	ToSqlOrderBy(c *builder.Context) (sql4prepare string)
	ToSqlLimitOffset(c *builder.Context) (sqlSegment string, binds []any)

	ToSql(c *builder.Context) (sql4prepare string, binds []any, err error)
	ToSqlInsert(c *builder.Context, obj any, args ...builder.TypeToSqlInsertCase) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(c *builder.Context, arg any) (sqlSegment string, binds []any, err error)
	ToSqlDelete(c *builder.Context, obj any, mustColumn ...string) (sqlSegment string, binds []any, err error)
}

type Driver struct {
	//driver string
	Dialect dialect.IDialect
}

func NewDriver(driver string) *Driver {
	return &Driver{Dialect: dialect.GetDialect(driver).New()}
}

func (d Driver) ToSql(c *builder.Context) (sql4prepare string, binds []any, err error) {
	sql4prepare, binds, err = d.toSql(c)
	if len(c.UnionClause.Unions) > 0 {
		for _, u := range c.UnionClause.Unions {
			sql4prepare2, binds2, err2 := u.ToSql()
			if err2 != nil {
				return
			}
			if sql4prepare2 == "" {
				continue
			}
			var unions = "UNION"
			if u.IsUnionAll {
				unions = "UNION ALL"
			}
			sql4prepare = fmt.Sprintf("%s %s %s", sql4prepare, unions, sql4prepare2)
			binds = append(binds, binds2...)
		}
	}
	return
}
func (d Driver) toSql(c *builder.Context) (sql4prepare string, binds []any, err error) {
	selects, anies := d.ToSqlSelect(c)
	table, binds2, err := d.ToSqlTable(c)
	if err != nil {
		return sql4prepare, binds2, err
	}
	joins, binds3, err := d.ToSqlJoin(c)
	if err != nil {
		return sql4prepare, binds3, err
	}
	wheres, binds4, err := d.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, binds4, err
	}
	orderBy := d.ToSqlOrderBy(c)
	limit, binds5 := d.ToSqlLimitOffset(c)
	groupBys := d.ToSqlGroupBy(c)
	havings, binds6, err := d.ToSqlHaving(c)

	binds = append(binds, anies...)
	binds = append(binds, binds2...)
	binds = append(binds, binds3...)
	binds = append(binds, binds4...)
	binds = append(binds, binds5...)
	binds = append(binds, binds6...)

	var locking string
	if c.PessimisticLocking == builder.TypeLockInShareMode {
		locking = d.Dialect.LockInShareMode()
	} else if c.PessimisticLocking == builder.TypeLockForUpdate {
		locking = d.Dialect.LockForUpdate()
	}

	//sql4prepare = NamedSprintf("SELECT :selects FROM :table :join :wheres :groupBys :havings :orderBy :pagination :PessimisticLocking", selects, table, joins, wheres, groupBys, havings, orderBy, limit, c.PessimisticLocking)
	sql4prepare = fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s %s %s", selects, table, joins, wheres, groupBys, havings, orderBy, limit, locking)
	sql4prepare = regexp.MustCompile(`\s{2,}`).ReplaceAllString(strings.TrimSpace(sql4prepare), " ")
	return
}

func (d Driver) ToSqlSelect(c *builder.Context) (sql4prepare string, binds []any) {
	var cols []string
	for _, col := range c.SelectClause.Columns {
		if col.IsRaw {
			cols = append(cols, col.Name)
			binds = append(binds, col.Binds...)
		} else {
			if col.Alias == "" {
				cols = append(cols, d.Dialect.QuoteIdentifier(col.Name))
			} else {
				cols = append(cols, fmt.Sprintf("%s AS %s", d.Dialect.QuoteIdentifier(col.Name), col.Alias))
			}
		}
	}
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	var distinct string
	if c.SelectClause.Distinct {
		distinct = "DISTINCT "
	}
	sql4prepare = fmt.Sprintf("%s%s", distinct, strings.Join(cols, ", "))
	return
}

func (d Driver) ToSqlTable(c *builder.Context) (sql4prepare string, binds []any, err error) {
	return d.buildSqlTable(c.TableClause, c.Prefix)
}

func (d Driver) buildSqlTable(tab builder.TableClause, prefix string) (sql4prepare string, binds []any, err error) {
	if v, ok := tab.Tables.(builder.IBuilder); ok {
		sql4prepare, binds, err = v.ToSql()
		if tab.Alias != "" {
			sql4prepare = fmt.Sprintf("(%s) %s", sql4prepare, d.Dialect.QuoteIdentifier(tab.Alias))
		}
		return
	}
	rfv := reflect.Indirect(reflect.ValueOf(tab.Tables))
	switch rfv.Kind() {
	case reflect.String:
		sql4prepare = d.Dialect.QuoteIdentifier(fmt.Sprintf("%s%s", prefix, tab.Tables))
	case reflect.Struct:
		sql4prepare = d.buildTableName(rfv.Type(), prefix)
	case reflect.Slice:
		if rfv.Type().Elem().Kind() == reflect.Struct {
			sql4prepare = d.buildTableName(rfv.Type().Elem(), prefix)
		} else {
			err = errors.New("table param must be string or struct(slice) bind with 1 or 2 params")
			return
		}
	default:
		err = errors.New("table must be string | struct | slice")
		return
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", sql4prepare, d.Dialect.QuoteIdentifier(tab.Alias))), binds, err
}

func (d Driver) toSqlWhere(c *builder.Context, wc builder.WhereClause) (sql4prepare string, binds []any, err error) {
	if len(wc.Conditions) == 0 {
		return
	}
	var sql4prepareArr []string
	for _, v := range wc.Conditions {
		switch item := v.(type) {
		case builder.TypeWhereRaw:
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s", item.LogicalOp, item.Column))
			binds = append(binds, item.Bindings...)
		case builder.TypeWhereStandard:
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s %s", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), item.Operator, d.Dialect.Placeholder()))
			binds = append(binds, item.Value)
		case builder.TypeWhereIn:
			values := ToSlice(item.Value)
			var phs []string
			for range values {
				phs = append(phs, d.Dialect.Placeholder())
			}
			//sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), item.Operator, strings.Repeat("?,", len(values)-1)+"?"))
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), item.Operator, strings.Join(phs, ",")))
			binds = append(binds, values...)
		case builder.TypeWhereBetween:
			values := ToSlice(item.Value)
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s %s AND %s", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), d.Dialect.Placeholder(), item.Operator, d.Dialect.Placeholder()))
			binds = append(binds, values...)
		case builder.TypeWhereNested:
			var tmp = builder.Context{}
			item.WhereNested(&tmp.WhereClause)
			prepare, anies, err := d.ToSqlWhere(&tmp)
			if err != nil {
				return sql4prepare, binds, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s (%s)", item.LogicalOp, strings.TrimPrefix(prepare, "WHERE ")))
			binds = append(binds, anies...)
		case builder.TypeWhereSubQuery:
			query, anies, err := item.SubQuery.ToSql()
			if err != nil {
				return sql4prepare, binds, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), item.Operator, query))
			binds = append(binds, anies...)
		case builder.TypeWhereSubHandler:
			var ctx = builder.NewContext(c.Prefix)
			item.Sub(ctx)
			query, anies, err := d.ToSql(ctx)
			if err != nil {
				return sql4prepare, binds, err
			}
			sql4prepareArr = append(sql4prepareArr, fmt.Sprintf("%s %s %s (%s)", item.LogicalOp, d.Dialect.QuoteIdentifier(item.Column), item.Operator, query))
			binds = append(binds, anies...)
		}
	}
	if len(sql4prepareArr) > 0 {
		sql4prepare = strings.TrimSpace(strings.Trim(strings.Trim(strings.TrimSpace(strings.Join(sql4prepareArr, " ")), "AND"), "OR"))
	}
	return
}
func (d Driver) ToSqlWhere(c *builder.Context) (sql4prepare string, binds []any, err error) {
	sql4prepare, binds, err = d.toSqlWhere(c, c.WhereClause)
	if sql4prepare != "" {
		if c.WhereClause.Not {
			sql4prepare = fmt.Sprintf("WHERE NOT %s", sql4prepare)
		}
		sql4prepare = fmt.Sprintf("WHERE %s", sql4prepare)
	}
	return
}

func (d Driver) ToSqlJoin(c *builder.Context) (sql4prepare string, binds []any, err error) {
	if c.JoinClause.Err != nil {
		return sql4prepare, binds, c.JoinClause.Err
	}
	if len(c.JoinClause.JoinItems) == 0 {
		return
	}
	for _, v := range c.JoinClause.JoinItems {
		var prepare string
		var sql4 string
		var bind []any
		switch item := v.(type) {
		case builder.TypeJoinStandard:
			prepare, bind, err = d.buildSqlTable(item.TableClause, c.Prefix)
			if err != nil {
				return
			}
			sql4 = fmt.Sprintf("%s %s ON %s %s %s", item.Type, prepare, d.Dialect.QuoteIdentifier(item.Column1), item.Operator, d.Dialect.QuoteIdentifier(item.Column2))
		case builder.TypeJoinSub:
			sql4, bind, err = item.ToSql()
			if err != nil {
				return
			}
		case builder.TypeJoinOn:
			var tjo builder.TypeJoinOnCondition
			item.OnClause(&tjo)
			if len(tjo.Conditions) == 0 {
				return
			}
			var sqlArr []string
			for _, cond := range tjo.Conditions {
				sqlArr = append(sqlArr, fmt.Sprintf("%s %s %s %s", cond.Relation, d.Dialect.QuoteIdentifier(cond.Column1), cond.Operator, d.Dialect.QuoteIdentifier(cond.Column2)))
			}

			sql4 = TrimPrefixAndOr(strings.Join(sqlArr, " "))
		}
		sql4prepare = fmt.Sprintf("%s %s", sql4prepare, sql4)
		binds = append(binds, bind...)
	}
	return
}

func (d Driver) ToSqlGroupBy(c *builder.Context) (sql4prepare string) {
	if len(c.GroupClause.Groups) > 0 {
		var tmp []string
		for _, col := range c.GroupClause.Groups {
			if col.IsRaw {
				tmp = append(tmp, col.Column)
			} else {
				tmp = append(tmp, d.Dialect.QuoteIdentifier(col.Column))
			}
		}
		sql4prepare = fmt.Sprintf("GROUP BY %s", strings.Join(tmp, ","))
	}
	return
}
func (d Driver) ToSqlHaving(c *builder.Context) (sql4prepare string, binds []any, err error) {
	sql4prepare, binds, err = d.toSqlWhere(c, c.HavingClause.WhereClause)
	if sql4prepare != "" {
		sql4prepare = fmt.Sprintf("HAVING %s", sql4prepare)
	}
	return
}
func (d Driver) ToSqlOrderBy(c *builder.Context) (sql4prepare string) {
	if len(c.OrderByClause.Columns) == 0 {
		return
	}
	var orderBys []string
	for _, v := range c.OrderByClause.Columns {
		if v.IsRaw {
			orderBys = append(orderBys, v.Column)
		} else {
			if v.Direction == "" {
				orderBys = append(orderBys, d.Dialect.QuoteIdentifier(v.Column))
			} else {
				orderBys = append(orderBys, fmt.Sprintf("%s %s", d.Dialect.QuoteIdentifier(v.Column), v.Direction))
			}
		}
	}
	sql4prepare = fmt.Sprintf("ORDER BY %s", strings.Join(orderBys, ", "))
	return
}

func (d Driver) ToSqlLimitOffset(c *builder.Context) (sqlSegment string, binds []any) {
	var offset int
	if c.LimitOffsetClause.Offset > 0 {
		offset = c.LimitOffsetClause.Offset
	} else if c.LimitOffsetClause.Page > 0 {
		offset = c.LimitOffsetClause.Limit * (c.LimitOffsetClause.Page - 1)
	}
	if c.LimitOffsetClause.Limit > 0 {
		sqlSegment = d.Dialect.LimitOffset(c.LimitOffsetClause.Limit, offset)
		//if offset > 0 {
		//	sqlSegment = "LIMIT ? OFFSET ?"
		//	binds = append(binds, c.LimitOffsetClause.Limit, offset)
		//} else {
		//	sqlSegment = "LIMIT ?"
		//	binds = append(binds, c.LimitOffsetClause.Limit)
		//}
	}
	return
}

// ToSqlInsert insert
func (d Driver) ToSqlInsert(c *builder.Context, obj any, args ...builder.TypeToSqlInsertCase) (sqlSegment string, binds []any, err error) {
	var arg builder.TypeToSqlInsertCase
	if len(args) > 0 {
		arg = args[0]
	}
	var ctx = *c
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		var datas []map[string]any
		datas, err = parser.StructsToInsert(obj, arg.MustColumn...)
		if err != nil {
			return
		}
		ctx.TableClause.Table(obj)
		return d.toSqlInsert(&ctx, datas, arg)
	case reflect.Slice:
		switch rfv.Type().Elem().Kind() {
		case reflect.Struct:
			c.TableClause.Table(obj)
			var datas []map[string]any
			datas, err = parser.StructsToInsert(obj, arg.MustColumn...)
			if err != nil {
				return
			}
			return d.toSqlInsert(c, datas, arg)
		default:
			return d.toSqlInsert(c, obj, arg)
		}
	default:
		return d.toSqlInsert(c, obj, arg)
	}
}

func (d Driver) ToSqlDelete(c *builder.Context, obj any, mustColumn ...string) (sqlSegment string, binds []any, err error) {
	var ctx = *c
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		data, err := parser.StructToDelete(obj, mustColumn...)
		if err != nil {
			return sqlSegment, binds, err
		}
		ctx.TableClause.Table(obj)
		ctx.WhereClause.Where(data)
		return d.toSqlDelete(&ctx)
	case reflect.Int64, reflect.Int32, reflect.String:
		ctx.WhereClause.Where("id", obj)
		return d.toSqlDelete(&ctx)
	default:
		err = errors.New("obj must be struct or id value")
	}
	return
}

func (d Driver) ToSqlUpdate(c *builder.Context, arg any) (sqlSegment string, binds []any, err error) {
	switch v := arg.(type) {
	case builder.TypeToSqlUpdateCase:
		return d.toSqlUpdate(c, v.BindOrData, v.MustColumn...)
	case builder.TypeToSqlIncDecCase:
		return d.toSqlIncDec(c, v.Symbol, v.Data)
	default:
		return
	}
}

func (d Driver) toSqlUpdate(c *builder.Context, obj any, mustColumn ...string) (sqlSegment string, binds []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	switch rfv.Kind() {
	case reflect.Struct:
		dataMap, pk, pkValue, err := parser.StructToUpdate(obj, mustColumn...)
		if err != nil {
			return sqlSegment, binds, err
		}
		var ctx = *c
		ctx.TableClause.Table(obj)
		if pk != "" {
			ctx.WhereClause.Where(pk, pkValue)
		}
		return d.toSqlUpdateReal(&ctx, dataMap)
	case reflect.Map:
		return d.toSqlUpdateReal(c, obj)
	default:
		err = errors.New("no support update obj")
		return
	}
}

func (d Driver) toSqlIncDec(c *builder.Context, symbol string, data map[string]any) (sql4prepare string, values []any, err error) {
	prepare, anies, err := d.ToSqlTable(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, anies...)

	var tmp []string
	for k, v := range data {
		tmp = append(tmp, fmt.Sprintf("%s=%s%s%s", d.Dialect.QuoteIdentifier(k), d.Dialect.QuoteIdentifier(k), symbol, d.Dialect.Placeholder()))
		values = append(values, v)
	}

	where, val, err := d.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, val...)

	sql4prepare = fmt.Sprintf("UPDATE %s SET %s %s", prepare, strings.Join(tmp, ","), where)
	return
}

func (d Driver) buildTableName(rft reflect.Type, prefix string) (tab string) {
	return d.Dialect.QuoteIdentifier(fmt.Sprintf("%s%s", prefix, parser.StructsToTableName(rft)))
}

// func (b Driver) toSqlInsert(c *gorose.Context, data any, ignoreCase string, onDuplicateKeys []string) (sql4prepare string, values []any, err error) {
func (d Driver) toSqlInsert(c *builder.Context, data any, insertCase builder.TypeToSqlInsertCase) (sql4prepare string, values []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	var fields []string
	var valuesPlaceholderArr []string
	switch rfv.Kind() {
	case reflect.Map:
		keys := rfv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		var valuesPlaceholderTmp []string
		for _, key := range keys {
			fields = append(fields, d.Dialect.QuoteIdentifier(key.String()))
			valuesPlaceholderTmp = append(valuesPlaceholderTmp, d.Dialect.Placeholder())
			values = append(values, rfv.MapIndex(key).Interface())
		}
		valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
	case reflect.Slice:
		if rfv.Len() == 0 {
			return
		}
		if rfv.Type().Elem().Kind() == reflect.Map {
			// 先获取到插入字段
			keys := rfv.Index(0).MapKeys()
			sort.Slice(keys, func(i, j int) bool {
				return keys[i].String() < keys[j].String()
			})
			for _, key := range keys {
				fields = append(fields, d.Dialect.QuoteIdentifier(key.String()))
			}
			// 组合插入数据
			for i := 0; i < rfv.Len(); i++ {
				var valuesPlaceholderTmp []string
				for _, key := range keys {
					valuesPlaceholderTmp = append(valuesPlaceholderTmp, d.Dialect.Placeholder())
					values = append(values, rfv.Index(i).MapIndex(key).Interface())
				}
				valuesPlaceholderArr = append(valuesPlaceholderArr, fmt.Sprintf("(%s)", strings.Join(valuesPlaceholderTmp, ",")))
			}
		} else {
			err = errors.New("only map(slice) data supported")
			return
		}
	default:
		err = errors.New("only map(slice) data supported")
		return
	}
	if err != nil {
		return
	}

	var onDuplicateKey string
	if len(insertCase.UpdateFields) > 0 {
		var tmp []string
		for _, v := range insertCase.UpdateFields {
			tmp = append(tmp, fmt.Sprintf("%s=VALUES(%s)", d.Dialect.QuoteIdentifier(v), d.Dialect.QuoteIdentifier(v)))
		}
		onDuplicateKey = fmt.Sprintf("%s %s", d.Dialect.Upsert(), strings.Join(tmp, ", "))
	}

	var insert = "INSERT"
	if insertCase.IsReplace {
		insert = "REPLACE"
	} else if insertCase.IsIgnoreCase {
		insert = "INSERT IGNORE"
	}

	var tables string
	tables, _, err = d.ToSqlTable(c)
	if err != nil {
		return
	}
	//sql4prepare = NamedSprintf(":insert INTO :tables (:fields) VALUES :placeholder :onDuplicateKey", insert, tables, strings.Join(fields, ","), strings.Join(valuesPlaceholderArr, ","), onDuplicateKey)
	sql4prepare = fmt.Sprintf("%s INTO %s (%s) VALUES %s %s", insert, tables, strings.Join(fields, ","), strings.Join(valuesPlaceholderArr, ","), onDuplicateKey)
	return
}

func (d Driver) toSqlUpdateReal(c *builder.Context, data any) (sql4prepare string, values []any, err error) {
	rfv := reflect.Indirect(reflect.ValueOf(data))
	var updates []string
	switch rfv.Kind() {
	case reflect.Map:
		keys := rfv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, key := range keys {
			updates = append(updates, fmt.Sprintf("%s = %s", d.Dialect.QuoteIdentifier(key.String()), d.Dialect.Placeholder()))
			values = append(values, rfv.MapIndex(key).Interface())
		}
	default:
		err = errors.New("only map data supported")
		return
	}
	var tables string
	tables, _, err = d.ToSqlTable(c)
	if err != nil {
		return
	}
	wheres, binds, err := d.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, binds...)

	//sql4prepare = NamedSprintf("UPDATE :tables SET :updates :wheres", tables, strings.Join(updates, ", "), wheres)
	sql4prepare = fmt.Sprintf("UPDATE %s SET %s %s", tables, strings.Join(updates, ", "), wheres)

	return
}

func (d Driver) toSqlDelete(c *builder.Context) (sql4prepare string, values []any, err error) {
	var tables string
	tables, _, err = d.ToSqlTable(c)
	if err != nil {
		return
	}
	wheres, binds, err := d.ToSqlWhere(c)
	if err != nil {
		return sql4prepare, values, err
	}
	values = append(values, binds...)
	//sql4prepare = NamedSprintf("DELETE FROM :tables :wheres", tables, wheres)
	sql4prepare = fmt.Sprintf("DELETE FROM %s %s", tables, wheres)
	return
}
