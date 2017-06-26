/*
 * 文件暂无名
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/6/26
 */
package controller

import (
	"../lib"
	"gopkg.in/kataras/iris.v6"
	"fmt"
)

type(
	Blog struct {
		sim.Model `xorm:"extends"`
		Title        string `xorm:"char(50)" json:"title"`
	}
)

/*
POST   /blog     创建
curl -X POST -d '{}' http://localhost:4000/blog
GET    /blog/1 查看
curl http://localhost:4000/blog/1
GET    /blogs
curl http://localhost:4000/blogs
DELETE /blog/1 删除
curl -X DELETE http://localhost:4000/blog/1
PUT    /blog/1 更新
curl -X PUT -d '{}' http://localhost:4000/blog/1
*/

func init() {
	const model = "blog"
	if web.Config.Debug {
		if exist,_ := web.Db.IsTableExist(model); !exist {
			web.Db.Sync2(new(Blog))
		}
	}

	web.Get(fmt.Sprintf("/%s/:id",model), func(c *iris.Context) {
		var r sim.Result
		var id, err = c.ParamInt64("id")

		if err != nil {
			r.Code = -1
			r.Message = "id unvalid"
			c.JSON(200, r)
			return
		}

		var bean Blog
		if _, err := web.Db.ID(id).Get(&bean); err != nil {
			r.Code = -2
			r.Message = "db id.get error"
			c.JSON(200, r)
			return
		}

		r.Code = 1
		r.Message = "success"
		r.Data = bean
		c.JSON(200, r)
	})
	web.Get(fmt.Sprintf("/%ss",model), func(c *iris.Context) {
		c.WriteString("gets ok")
	})
	web.Post(fmt.Sprintf("/%s",model), func(c *iris.Context) {
		c.WriteString("post ok")
	})
	web.Delete(fmt.Sprintf("/%s/:id",model), func(c *iris.Context) {
		c.WriteString("delete ok")
	})
	web.Put(fmt.Sprintf("/%s/:id",model), func(c *iris.Context) {
		c.WriteString("put ok")
	})
}
