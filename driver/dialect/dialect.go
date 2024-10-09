package dialect

import "sync"

type IDialect interface {
	New() IDialect
	Placeholder() string                  // 占位符处理，例如 MySQL 使用 `?`, PostgreSQL 使用 `$1`
	AutoIncrement() string                // 自增字段的声明方式
	LimitOffset(limit, offset int) string // 分页查询的 SQL 片段
	//InsertQuery(table string, columns []string, values [][]interface{}) (string, []interface{}) // 批量插入 SQL 生成
	QuoteIdentifier(identifier string) string // 对标识符（字段名、表名）加引号
	Upsert() string

	LockInShareMode() string
	LockForUpdate() string
}

var dialectMap = map[string]IDialect{}
var dialectLock sync.RWMutex

func Register(driver string, parser IDialect) {
	dialectLock.Lock()
	defer dialectLock.Unlock()
	dialectMap[driver] = parser
}

func GetDialect(driver string) IDialect {
	dialectLock.RLock()
	defer dialectLock.RUnlock()
	return dialectMap[driver]
}

func DialectList() (dr []string) {
	dialectLock.RLock()
	defer dialectLock.RUnlock()
	for d := range dialectMap {
		dr = append(dr, d)
	}
	return
}
