package sim

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const (
	DB_DRIVER_MYSQL = "mysql"
	DB_DRIVER_SQLITE3 = "sqlite3"
)

type DB struct {
	*xorm.Engine
	DataSource,DriverName string
	ShowSql bool
}

func (this *DB) Init() {
	if engine, err := xorm.NewEngine(this.DriverName, this.DataSource); err == nil {
		engine.SetMapper(core.GonicMapper{})
		this.Engine = engine
		this.Engine.ShowSQL(this.ShowSql)
		//this.Engine.ShowExecTime(true)
	}else{
		Debug("db init err", err.Error())
	}
}