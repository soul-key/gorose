package gorose

import (
	"github.com/gohouse/gorose/v3/builder"
	"math"
)

// Pagination 是用于分页查询结果的结构体，包含当前页数据及分页信息。
type Pagination struct {
	Limit       int              `json:"limit"`
	Pages       int              `json:"pages"`
	CurrentPage int              `json:"currentPage"`
	PrevPage    int              `json:"prevPage"`
	NextPage    int              `json:"nextPage"`
	Total       int64            `json:"total"`
	Data        []map[string]any `json:"data"`
}

func (db *Database) Paginate(obj ...any) (result Pagination, err error) {
	if len(obj) > 0 {
		db.Table(obj[0])
	}
	var count int64
	count, err = db.Count()
	if err != nil || count == 0 {
		return
	}
	if db.Context.LimitOffsetClause.Limit == 0 {
		db.Limit(15)
	}
	if db.Context.LimitOffsetClause.Page == 0 {
		db.Page(1)
	}

	res, err := db.Get()
	if err != nil {
		return result, err
	}

	result.Total = count
	result.Data = res
	result.Limit = db.Context.LimitOffsetClause.Limit
	result.Pages = int(math.Ceil(float64(count) / float64(db.Context.LimitOffsetClause.Limit)))
	result.CurrentPage = db.Context.LimitOffsetClause.Page
	result.PrevPage = db.Context.LimitOffsetClause.Page - 1
	result.NextPage = db.Context.LimitOffsetClause.Page + 1
	if db.Context.LimitOffsetClause.Page == 1 {
		result.PrevPage = 1
	}
	if db.Context.LimitOffsetClause.Page == result.Pages {
		result.NextPage = result.Pages
	}
	return
}


func (db *Database) WhereSub(column string, operation string, sub builder.WhereSubHandler) *Database {
	db.Context.WhereClause.WhereSub(column, operation, sub)
	return db
}
func (db *Database) OrWhereSub(column string, operation string, sub builder.WhereSubHandler) *Database {
	db.Context.WhereClause.OrWhereSub(column, operation, sub)
	return db
}
func (db *Database) WhereBuilder(column string, operation string, sub builder.IBuilder) *Database {
	db.Context.WhereClause.WhereBuilder(column, operation, sub)
	return db
}
func (db *Database) OrWhereBuilder(column string, operation string, sub builder.IBuilder) *Database {
	db.Context.WhereClause.OrWhereBuilder(column, operation, sub)
	return db
}
func (db *Database) WhereNested(handler builder.WhereNestedHandler) *Database {
	db.Context.WhereClause.WhereNested(handler)
	return db
}
func (db *Database) OrWhereNested(handler builder.WhereNestedHandler) *Database {
	db.Context.WhereClause.OrWhereNested(handler)
	return db
}
func (db *Database) WhereIn(column string, value any) *Database {
	db.Context.WhereClause.WhereIn(column, value)
	return db
}
func (db *Database) OrWhereIn(column string, value any) *Database {
	db.Context.WhereClause.WhereIn(column, value)
	return db
}
func (db *Database) WhereNull(column string) *Database {
	db.Context.WhereClause.WhereNull(column)
	return db
}
func (db *Database) OrWhereNull(column string) *Database {
	db.Context.WhereClause.WhereNull(column)
	return db
}
func (db *Database) WhereBetween(column string, value any) *Database {
	db.Context.WhereClause.WhereBetween(column, value)
	return db
}
func (db *Database) OrWhereBetween(column string, value any) *Database {
	db.Context.WhereClause.WhereBetween(column, value)
	return db
}
func (db *Database) WhereExists(clause builder.IBuilder) {
	db.Context.WhereClause.WhereExists(clause)
}
func (db *Database) WhereLike(column, value string) *Database {
	db.Context.WhereClause.WhereLike(column, value)
	return db
}
func (db *Database) OrWhereLike(column, value string) *Database {
	db.Context.WhereClause.OrWhereLike(column, value)
	return db
}
//func (db *Database) WhereNotIn(column string, value any) *Database {
//	db.Context.WhereClause.WhereNotIn(column, value)
//	return db
//}
//func (db *Database) OrWhereNotIn(column string, value any) *Database {
//	db.Context.WhereClause.WhereNotIn(column, value)
//	return db
//}
//func (db *Database) WhereNotNull(column string) *Database {
//	db.Context.WhereClause.whereNull("AND", column, true)
//	return db
//}
//func (db *Database) OrWhereNotNull(column string) *Database {
//	db.Context.WhereClause.whereNull("OR", column, true)
//	return db
//}
//func (db *Database) WhereNotBetween(column string, value any) *Database {
//	db.Context.WhereClause.whereBetween("AND", column, value, true)
//	return db
//}
//func (db *Database) OrWhereNotBetween(column string, value any) *Database {
//	db.Context.WhereClause.whereBetween("OR", column, value, true)
//	return db
//}
//func (db *Database) WhereNotExists(clause IBuilder) {
//	db.Context.WhereClause.WhereNotExists(clause)
//}
//func (db *Database) WhereNotLike(column, value string) *Database {
//	db.Context.WhereClause.whereLike("AND", column, value, true)
//	return db
//}
//func (db *Database) OrWhereNotLike(column, value string) *Database {
//	db.Context.WhereClause.whereLike("OR", column, value, true)
//	return db
//}
func (db *Database) WhereNot(column any, args ...any) *Database {
	db.Context.WhereClause.WhereNot(column, args...)
	return db
}

func (db *Database) OrderByAsc(column string) *Database {
	return db.OrderBy(column, "ASC")
}

func (db *Database) OrderByDesc(column string) *Database {
	return db.OrderBy(column, "DESC")
}
