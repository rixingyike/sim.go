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
	"math"
)

func init() {
	type (
		Joke struct {
			sim.Model `xorm:"extends"`
			UserId int64 `json:"user_id"`
			Title        string `xorm:"char(50)" json:"title"`
			Content        string `xorm:"text" json:"content"`
			Image        string `xorm:"tinytext" json:"image"`

			User sim.WeappUser `xorm:"-" json:"user"`
		}
		JokePage struct {
			List       []Joke `json:"list"`
			Size       int `json:"size"`
			Page       int `json:"page"`   //自1开始计数
			Count      int `json:"count"`
			TotalPage  int `json:"total_page"`
		}
	)

	const URL_BASE = ""
	const MODEL_NAME = "joke"

	if web.DB == nil {
		return
	}
	if web.Config.Debug {
		if exist,_ := web.DB.IsTableExist(MODEL_NAME); !exist {
			web.DB.Sync2(new(Joke))
		}
	}

	/*
	使用curl测试接口
	POST   /blog     创建
	curl -X POST -d '{"title":"x"}' http://localhost:4000/blog

	GET    /blog/1 查看
	curl http://localhost:4000/blog/1

	GET    /blogs
	curl http://localhost:4000/blogs

	DELETE /blog/1 删除
	curl -X DELETE http://localhost:4000/blog/1

	PUT    /blog/1 更新
	curl -X PUT -d '{"title":"xxx"}' http://localhost:4000/blog/1
	*/

	var sub = web.Router
	if URL_BASE != "" {
		sub = web.Party(URL_BASE)
	}
	// 拉取单条记录
	sub.Get(fmt.Sprintf("/%s/:id", MODEL_NAME), func(c *iris.Context) {
		var r sim.Result

		var id, err = c.ParamInt64("id")
		if err != nil {
			r.Code = -1
			r.Message = "id unvalid"
			c.JSON(200, r)
			return
		}

		var data Joke
		if _, err := web.DB.ID(id).Get(&data); err != nil {
			r.Code = -2
			r.Message = "db id.get error"
			c.JSON(200, r)
			return
		}

		r.Code = 1
		r.Message = "success"
		r.Data = data
		c.JSON(200, r)
	})

	// 分页拉取记录,如果没有指定size,拉取所有限1000条以内
	sub.Get(fmt.Sprintf("/%ss", MODEL_NAME), func(c *iris.Context) {
		var r sim.Result
		var page,_ = c.URLParamInt("page")
		var size,_ = c.URLParamInt("size")

		if page == 0 {
			page = 1
		}
		if size == 0 {
			size = 1000
		}
		var data = JokePage{Page:page,Size:size}
		var offset = (data.Page - 1) * data.Size

		if err := web.DB.Desc("id").Limit(data.Size, offset).Find(&data.List); err == nil {
			if total, err := web.DB.Count(new(Joke)); err == nil {
				data.Count = int(total)
				data.TotalPage = int(math.Ceil(float64(total) / float64(data.Size)))
			}else{
				sim.Debug("select count err", err.Error())
			}
			for k,v := range data.List{
				web.DB.ID(v.UserId).Get(&v.User)
				data.List[k].User = v.User
			}
			r.Code = 1
			r.Data = data
		}else{
			sim.Debug("select page err", err.Error())
		}

		c.JSON(200, r)
	})

	// 新增单条记录
	sub.Post(fmt.Sprintf("/%s", MODEL_NAME), func(c *iris.Context) {
		var r sim.Result

		var data Joke
		if err := c.ReadJSON(&data); err != nil {
			sim.Debug("read data err", err.Error())
			r.Code = -1
			r.Message = "read data err"
			c.JSON(200, r)
			return
		}

		var my = web.GetWeappUser(c)
		data.UserId = my.ID

		if affected, err := web.DB.Insert(&data); err == nil {
			if affected > 0 {
				r.Code = 1
				r.Data = data.ID
			}
		}else{
			sim.Debug("post new bean err",err.Error())
		}

		c.JSON(200, r)
	})

	// 删除单条记录
	sub.Delete(fmt.Sprintf("/%s/:id", MODEL_NAME), func(c *iris.Context) {
		var r sim.Result

		var id, err = c.ParamInt64("id")
		if err != nil {
			r.Code = -1
			r.Message = "id unvalid"
			c.JSON(200, r)
			return
		}

		var data Joke
		if _, err := web.DB.ID(id).Get(&data); err != nil {
			r.Code = -2
			r.Message = "not found"
			c.JSON(200, r)
			return
		}

		if affected, err := web.DB.Id(id).Delete(&data); err == nil {
			if affected > 0 {
				r.Code = 1
			}
		}

		c.JSON(200, r)
	})

	// 更新单条记录
	sub.Put(fmt.Sprintf("/%s/:id", MODEL_NAME), func(c *iris.Context) {
		var r sim.Result

		var id, err = c.ParamInt64("id")
		if err != nil {
			r.Code = -1
			r.Message = "id unvalid"
			c.JSON(200, r)
			return
		}

		var data Joke
		if _, err := web.DB.ID(id).Get(&data); err != nil {
			r.Code = -2
			r.Message = "not found"
			c.JSON(200, r)
			return
		}

		if err := c.ReadJSON(&data); err != nil {
			r.Code = -3
			r.Message = "read data err"
			c.JSON(200, r)
			return
		}

		if _, err := web.DB.Id(id).Update(&data); err == nil {
			//如果当前内容与原内容一样,affected会返回0
			r.Code = 1
			r.Data = data.ID
		}

		c.JSON(200, r)
	})
}
