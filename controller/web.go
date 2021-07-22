/*
 * 文件暂无名
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/6/26
 */
 package controller

 import (
	 "sim.go/lib"
 )
 
 type Web struct {
	 *sim.WebEngine
 }
 
 var web = &Web{WebEngine:sim.Web}
 
 func init() {
 }
 
 func Start()  {
	 web.Start()
 }