/*
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/3/21
 */
package sim

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"fmt"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"os"
)

type IrisWebEngine struct {
	*iris.Framework
	Config Config
	Db *MysqlDb
}

var Web = &IrisWebEngine{}

func init() {
	fmt.Println("iris init")
	Web.Init()
}

func (this *IrisWebEngine) Init() {
	var confileFile = "./config.ini"
	if _, err := os.Stat(confileFile); os.IsNotExist(err) {
		fmt.Println("config.ini文件未找到")
		panic(err)
	}

	// 无论是以gin热编译,还是直接启动,都能从程序执行的当前目录找到config.ini
	if err := ReadTomlConfig(&this.Config, confileFile); err == nil {
		this.Log("已读取配置文件")
		this.Log("config", ToMapObject(this.Config))
	}else{
		panic(err)
	}

	this.Framework = iris.New(iris.Configuration{Gzip: true, Charset: "UTF-8"})
	this.Framework.Adapt(httprouter.New())
	if this.Config.Debug {
		this.Framework.Adapt(iris.DevLogger())
	}

	// 开始CORS跨域支持
	if this.Config.Cors.Enable {
		this.Framework.Adapt(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowCredentials:true,
		}))
	}

	// 初始化数据库
	if this.Config.Mysql.Enable {
		this.Db = &MysqlDb{ShowSql:this.Config.Debug,DataSource:this.Config.Mysql.DataSource}
		this.Db.Init()
	}

	// 初始化模板
	if this.Config.Html.Enable {
		this.Framework.StaticWeb("/static", this.Config.Html.StaticDir)
		var djangoAdapt = view.Django(this.Config.Html.TemplateDir, ".html")
		//原来的Config.IsDevelopment转移到了这里,设置为true,表示模板文件热加载
		djangoAdapt.Reload(this.Config.Debug)
		this.Framework.Adapt(djangoAdapt)
	}

	// 初始化session
	if this.Config.Session.Enable {
		var sessionAdapt = sessions.New(sessions.Config{Cookie: this.Config.Session.Key})
		this.Framework.Adapt(sessionAdapt)
	}

	this.Any("/hi", func(c *iris.Context) {
		c.WriteString("hi,sim.go")
	})
}

func (this *IrisWebEngine) Start() {
	this.Framework.Listen(this.Config.Addr)
}

func (this *IrisWebEngine) Log(v ...interface{}) {
	if this.Config.Debug {
		Debug(v...)
	}
}