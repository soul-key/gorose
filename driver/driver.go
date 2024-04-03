package driver

import (
	"github.com/gohouse/gorose/v3/builder"
	"sync"
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

var driverMap = map[string]IDriver{}
var driverLock sync.RWMutex

func Register(driver string, parser IDriver) {
	driverLock.Lock()
	defer driverLock.Unlock()
	driverMap[driver] = parser
}

func GetDriver(driver string) IDriver {
	driverLock.RLock()
	defer driverLock.RUnlock()
	return driverMap[driver]
}

func DriverList() (dr []string) {
	driverLock.RLock()
	defer driverLock.RUnlock()
	for d := range driverMap {
		dr = append(dr, d)
	}
	return
}
