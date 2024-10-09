package gorose

import (
	"github.com/gohouse/gorose/v3/driver"
	"testing"
)

type User struct {
	Id   int64  `db:"id,pk"`
	Name string `db:"name"`
}

// var dbg = Open("mysql") // just test toSql
// var dbg = Open("postgresql") // just test toSql
// var dbg = Open("mssql") // just test toSql
// var dbg = Open("oracle") // just test toSql
var dbg = Open("sqlite3") // just test toSql

func db() *Database {
	return dbg.NewDatabase()
}

func TestDatabase_ToSqlTo(t *testing.T) {
	var user = User{Id: 1}
	prepare, values, err := db().ToSqlTo(&user)
	driver.AssertsError(t, err)
	//db().Driver.Dialect.Placeholder()
	var expect = map[string]string{
		"mysql":      "SELECT `id`, `name` FROM `User` WHERE `id` = ? LIMIT 1",
		"postgresql": `SELECT "id", "name" FROM "User" WHERE "id" = $1 LIMIT 1`,
		"mssql":      `SELECT [id], [name] FROM [User] WHERE [id] = @p1 OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY`,
		"oracle":     `SELECT "id", "name" FROM "User" WHERE "id" = @p1 OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY`,
		"sqlite3":    `SELECT "id", "name" FROM "User" WHERE "id" = ? LIMIT 1`,
	}
	driver.AssertsEqual(t, expect[dbg.driver], prepare)
	var expectValues = []interface{}{1}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlToSlice(t *testing.T) {
	var user []User
	prepare, values, err := db().Where("id", ">", 1).OrderBy("id").Limit(10).Page(2).ToSqlTo(&user)
	driver.AssertsError(t, err)
	var expect = "SELECT `id`, `name` FROM `User` WHERE `id` > ? ORDER BY `id` LIMIT 10 OFFSET 10"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []interface{}{1}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSql(t *testing.T) {
	prepare, values, err := db().Table("users").Select("b").Where("c", 1).GroupBy("a").Having("a", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	driver.AssertsError(t, err)
	var expect = "SELECT `b` FROM `users` WHERE `c` = ? GROUP BY `a` HAVING `a` = ? ORDER BY `id` LIMIT 10 OFFSET 10"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []int{1, 1}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlInsert(t *testing.T) {
	var user = User{Name: "john"}
	prepare, values, err := db().ToSqlInsert(&user)
	driver.AssertsError(t, err)
	var expect = "INSERT INTO `User` (`name`) VALUES (?)"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []string{"john"}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlInserts(t *testing.T) {
	var user = []User{{Name: "John"}, {Name: "Alice"}}
	prepare, values, err := db().ToSqlInsert(&user)
	driver.AssertsError(t, err)
	var expect = "INSERT INTO `User` (`name`) VALUES (?),(?)"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []string{"John", "Alice"}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlUpdate(t *testing.T) {
	var user = User{Id: 1, Name: "john"}
	prepare, values, err := db().ToSqlUpdate(&user)
	driver.AssertsError(t, err)
	var expect = "UPDATE `User` SET `name` = ? WHERE `id` = ?"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []any{"john", 1}
	driver.AssertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlDelete(t *testing.T) {
	var user = User{Id: 1}
	prepare, values, err := db().ToSqlDelete(&user, "name")
	driver.AssertsError(t, err)
	var expect = "DELETE FROM `User` WHERE `id` = ? AND `name` = ?"
	driver.AssertsEqual(t, expect, prepare)
	var expectValues = []any{1, ""}
	driver.AssertsEqual(t, expectValues, values)
}
