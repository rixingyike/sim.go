package sim

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type MysqlDb struct {
	*xorm.Engine
	DataSource string
	ShowSql bool
}

func (this *MysqlDb) Init() {
	if engine, err := xorm.NewEngine("mysql", this.DataSource); err == nil {
		engine.SetMapper(core.GonicMapper{})
		this.Engine = engine
		this.Engine.ShowSQL(this.ShowSql)
		this.Engine.ShowExecTime(this.ShowSql)
	}
}