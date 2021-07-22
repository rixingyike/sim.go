/*
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/3/21
 */
 package sim

 import (
	 "github.com/kataras/iris/v12"
	 "github.com/kataras/iris/v12/middleware/logger"
	 "github.com/kataras/iris/v12/middleware/recover"
	 "github.com/kataras/iris/v12/sessions"
	 "github.com/iris-contrib/middleware/cors"
	 "fmt"
	 "os"
	 "strconv"
	 "encoding/json"
 )
 
 type WebEngine struct {
	Framework *iris.Application
	 Config Config
	 DB     *DB
	 Weapp *WeappController
	 Qiniu *QiniuController
	 Sessions *sessions.Sessions
 }
 
 var Web = &WebEngine{}
 
 func init() {
	 fmt.Println("web engine init")
	 Web.Init()
 }
 
 // 获取当前登陆的微信小程序用户
 func (this *WebEngine) GetWeappUser(c iris.Context) (u WeappUser) {
	 const USER_KEY = "WEAPPUSER"
	 var id = c.GetHeader("X-WX-Id")
	 //  Debug("id",id)
 
	 if c.GetCookie(USER_KEY) != "" {
		 userJson := c.GetCookie(USER_KEY)
		 json.Unmarshal([]byte(userJson), &u)
	 }else if id, err := strconv.ParseInt(id,10,64); err == nil {
		 this.DB.ID(id).Get(&u)
		 userJson,_ := json.Marshal(u)
		 c.SetCookieKV(USER_KEY, string(userJson))
	 }
 
	 return u
 }
 
 func (this *WebEngine) Init() {
	 var confileFile = "./config.ini"
	 if _, err := os.Stat(confileFile); os.IsNotExist(err) {
		 fmt.Println("config.ini文件未找到")
		 panic(err)
	 }
 
	 // 无论是以gin热编译,还是直接启动,都能从程序执行的当前目录找到config.ini
	 if err := ReadTomlConfig(&this.Config, confileFile); err == nil {
		 this.Log("已读取配置文件")
		 //this.Log("config", ToMapObject(this.Config))
	 }else{
		 panic(err)
	 }
 
	 this.Framework = iris.New()
	 //  this.Framework.Adapt(httprouter.New())
	 if this.Config.Debug {
		this.Framework.Logger().SetLevel("debug")
		// 设置recover从panics恢复，设置log记录
		this.Framework.Use(recover.New())
		this.Framework.Use(logger.New())
	 }
 
	 // 开始CORS跨域支持
	 if this.Config.Cors.Enable {
			this.Framework.UseRouter(cors.New(cors.Options{
					AllowedOrigins:   []string{"*"},
					AllowCredentials: true,
			}))
	 }
 
	 // 初始化数据库
	 if this.Config.Mysql.Enable {
		 this.DB = &DB{ShowSql:this.Config.Debug,DataSource:this.Config.Mysql.DataSource, DriverName:DB_DRIVER_MYSQL}
		 this.DB.Init()
	 }else if this.Config.Sqlite3.Enable {
		 Debug("init sqlite3 db")
		 this.DB = &DB{ShowSql:this.Config.Debug,DataSource:this.Config.Sqlite3.Filepath, DriverName:DB_DRIVER_SQLITE3}
		 this.DB.Init()
	 }
 
	 // 启用小程序用户自动登陆
	 if this.DB != nil && this.Config.Weapp.Enable {
		 if exist,_ := this.DB.IsTableExist(new(WeappUser)); !exist {
			 this.DB.Sync2(new(WeappUser))
		 }
		 this.Weapp = &WeappController{Web:this,AppId:this.Config.Weapp.AppId,AppSecret:this.Config.Weapp.AppSecret}
		 this.Weapp.Init()
	 }
 
	 // 初始化网页模板
	 if this.Config.Html.Enable {
			djangoAdapt := iris.Django(this.Config.Html.TemplateDir, ".html")
			djangoAdapt.Reload(this.Config.Debug)
			this.Framework.Layout("layout.html")
	 }
 
	 // 初始化静态目录
	 if this.Config.Static.Enable {
		 this.Framework.HandleDir("/static", this.Config.Static.StaticDir)
	 }
 
	 // 初始化session
	 if this.Config.Session.Enable {
		this.Sessions = sessions.New(sessions.Config{Cookie: this.Config.Session.Key})
	 }
 
	 // 启用七牛图片上传
	 if this.Config.Qiniu.Enable {
		 var qs = this.Config.Qiniu
		 this.Qiniu = &QiniuController{Web:this,Scope:qs.Scope,AccessKey:qs.AccessKey,SecretKey:qs.SecretKey,ServerBase:qs.ServerBase,Watermark:qs.Watermark}
		 this.Qiniu.Init()
	 }
 
	 this.Framework.Any("/hi", func(c iris.Context) {
		 if user := this.GetWeappUser(c); user.ID > 0 {
			 c.WriteString(fmt.Sprintf("hi,%s",user.Nickname))
			 return
		 }
		 c.WriteString("hi,sim.go")
	 })
 }
 
 func (this *WebEngine) Start() {
	 this.Framework.Listen(this.Config.Addr)
 }
 
 func (this *WebEngine) Log(v ...interface{}) {
	 if this.Config.Debug {
		 Debug(v...)
	 }
 }