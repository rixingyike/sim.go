package sim

import "github.com/kataras/iris/v12"

type Result struct {
	Code	int `json:"code"`
	Message	string `json:"message"`
	Data interface{} `json:"data"`
}

func (r Result) ToMap() iris.Map {
	var m = iris.Map{
		"code": r.Code,
		"message": r.Message,
		"data": r.Data}

	return m
}