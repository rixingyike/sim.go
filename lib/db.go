package sim

import (
	_ "github.com/go-sql-driver/mysql"
	// sqlite3需要下载gcc支持，在win上
	// _ "github.com/mattn/go-sqlite3" 
	"xorm.io/core"
	"xorm.io/xorm"
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