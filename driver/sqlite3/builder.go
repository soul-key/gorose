package sqlite3

import (
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v3/driver"
	"github.com/gohouse/gorose/v3/driver/mysql"
)

const DriverName = "sqlite3"

type Builder struct {
	//prefix string
	mysql.Builder
}

func init() {
	driver.Register(DriverName, &Builder{})
}
