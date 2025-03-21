// Examples for join
//
//	db.Table("users").Join("card", "users.id", "card.uid")
//	db.Table("users").Join("card", "users.id", "=", "card.uid")
//	db.Table("users").Join("card", "users.id", "=", "card.uid", "left")
//	db.Table("users").Join("card", "users.age", "=", db().Table("user_info").Max("age"))
//	db.Table("users").Join(db().Table("user_info").Where("uid", 2))
//	db.xxx.LeftJoin(xxx, <same as before>)
//	db.Table(gorose.As("users", "a")).Join(gorose.As("card", "b"), "a.id", "b.uid")

package builder

import (
	"errors"
)

type IJoinOn interface {
	On(column string, args ...string) IJoinOn
	OrOn(column string, args ...string) IJoinOn
}

// JoinClause 描述JOIN操作
type JoinClause struct {
	JoinItems []any
	Err       error
}

// TypeJoinSub 描述JOIN操作
type TypeJoinSub struct {
	IBuilder
}
type TypeJoinStandard struct {
	TableClause
	Type     string // JOIN类型（INNER, LEFT, RIGHT等）
	Column1  string
	Operator string
	Column2  string
}
type TypeJoinOn struct {
	TableClause
	OnClause func(IJoinOn)
	Type     string // JOIN类型（INNER, LEFT, RIGHT等）
}
type TypeJoinOnCondition struct {
	Conditions []TypeJoinOnConditionItem
}
type TypeJoinOnConditionItem struct {
	Relation string // and/or
	Column1  string
	Operator string
	Column2  string
}

func (jc *TypeJoinOnCondition) onClause(relation, column1 string, column2OrOperator ...string) IJoinOn {
	if len(column2OrOperator) == 0 {
		return jc
	}
	var operator = "="
	var column2 = column2OrOperator[0]
	if len(column2OrOperator) == 2 {
		operator = column2OrOperator[0]
		column2 = column2OrOperator[1]
	}
	jc.Conditions = append(jc.Conditions, TypeJoinOnConditionItem{
		Relation: relation,
		Column1:  column1,
		Operator: operator,
		Column2:  column2,
	})
	return jc
}
func (jc *TypeJoinOnCondition) On(column1 string, operatorOrColumn2 ...string) IJoinOn {
	return jc.onClause("AND", column1, operatorOrColumn2...)
}
func (jc *TypeJoinOnCondition) OrOn(column1 string, operatorOrColumn2 ...string) IJoinOn {
	return jc.onClause("OR", column1, operatorOrColumn2...)
}

func (db *JoinClause) join(joinType string, table any, argOrFn ...any) *JoinClause {
	var tab TableClause
	switch table.(type) {
	case string:
		tab.Tables = table
	case TableClause:
		tab = table.(TableClause)
	case IBuilder:
		db.JoinItems = append(db.JoinItems, TypeJoinSub{table.(IBuilder)})
		return db
	}

	switch len(argOrFn) {
	case 1:
		if v, ok := argOrFn[0].(func(on IJoinOn)); ok {
			db.JoinItems = append(db.JoinItems, TypeJoinOn{
				TableClause: tab,
				OnClause:    v,
				Type:        joinType,
			})
		}
	case 2:
		db.JoinItems = append(db.JoinItems, TypeJoinStandard{
			TableClause: tab,
			Column1:     argOrFn[0].(string),
			Operator:    "=",
			Column2:     argOrFn[1].(string),
			Type:        joinType,
		})
	case 3:
		db.JoinItems = append(db.JoinItems, TypeJoinStandard{
			TableClause: tab,
			Column1:     argOrFn[0].(string),
			Operator:    argOrFn[1].(string),
			Column2:     argOrFn[2].(string),
			Type:        joinType,
		})
	default:
		db.Err = errors.New("join args error")
	}
	return db
}

func (db *JoinClause) Join(table any, argOrFn ...any) *JoinClause {
	return db.join("INNER JOIN", table, argOrFn...)
}
func (db *JoinClause) JoinOn(table any, fn func(on IJoinOn)) *JoinClause {
	return db.join("INNER JOIN", table, fn)
}

// LeftJoin 描述LEFT JOIN操作
func (db *JoinClause) LeftJoin(table any, argOrFn ...any) *JoinClause {
	return db.join("LEFT JOIN", table, argOrFn...)
}

// RightJoin 描述RIGHT JOIN操作
func (db *JoinClause) RightJoin(table any, argOrFn ...any) *JoinClause {
	return db.join("RIGHT JOIN", table, argOrFn...)
}

// CrossJoin 描述CROSS JOIN操作
func (db *JoinClause) CrossJoin(table any, argOrFn ...any) *JoinClause {
	return db.join("CROSS JOIN", table, argOrFn...)
}
