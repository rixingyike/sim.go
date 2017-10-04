/*
 * 通用数据模型
 * author: sban
 * email: 9830131#qq.com
 * date: 2017/3/16
 */
package sim

type (
	Model struct {
		// 如果ID没有pk tag，使用DB.ID(id)条件时，将报错：‘ID condition is error, expect 0 primarykeys, there are 1’
		ID      int64 `xorm:"autoincr pk" json:"id"`
		Updated int64 `xorm:"updated" json:"updated"` //`xorm:"updated_at"`
		Created int64 `xorm:"created" json:"created"` //`xorm:"created_at" json:"created_at"`
		Deleted int64 `xorm:"deleted" json:"-"`
	}
)
