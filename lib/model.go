/*
 * 通用数据模型
 * author: sban
 * email: 9830131#qq.com
 * date: 2017/3/16
 */
package sim

type(
	Model struct {
		ID      int64        `xorm:"autoincr" json:"id"`
		Updated int64 `xorm:"updated" json:"updated"` //`xorm:"updated_at"`
		Created int64 `xorm:"created" json:"created"` //`xorm:"created_at" json:"created_at"`
		Deleted int64 `xorm:"deleted" json:"-"`
	}
)