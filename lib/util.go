/*
 * 文件暂无名
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/6/26
 */
package sim

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
	"encoding/json"
)

// 将对象转为map对象,字段以小写字段显示,终端打印显示字段名
func ToMapObject(v interface{}) (r interface{}) {
	if jsonStr, err := ToJsonString(v); err == nil {
		ParseJsonString(&r, jsonStr)
	}
	return
}

// json string生成
func ToJsonString(v interface{}) (result string,err error) {
	if bytes, err := json.Marshal(v); err == nil {
		result = string(bytes)
	}
	return
}

// 从json中反序列化对象
func ParseJsonString(v interface{}, s string) (err error) {
	//error:invalid character '\n' in string literal
	s = strings.TrimSpace(s)
	err = json.Unmarshal([]byte(s), v)
	return
}

// 读取toml配置文件内容
func ReadTomlConfig(ref interface{}, filePath string) (err error) {
	//如果toml文件中没有结构体中定义的字段,例如WebDescription,程序不会报错
	_, err = toml.DecodeFile(filePath, ref)
	return
}

// 控制台打印消息
func Debug(v ...interface{}) {
	fmt.Println(v...)
}
