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
	"io/ioutil"
	"net/url"
	"net/http"
	"bufio"
	"io"
	"net/http/cookiejar"
	"github.com/Jeffail/gabs"
)



//发起一个post请求,获取字符串内容
//@values params参数,可以为nil
//@headers 可为nil
func HttpPost(link string, values url.Values, headers map[string]string) (result string) {
	var requestBody io.Reader
	if values != nil {
		requestBody = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequest("POST", link, requestBody)

	if headers == nil {
		headers = map[string]string{}
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	jar,_ := cookiejar.New(nil)
	client := &http.Client{
		Jar:jar,
	}

	if headers["Cookie"] != "" {
		req.Header.Set("Set-Cookie",headers["Cookie"])

		var url2,_ = url.Parse(link)
		if cookies,err := ParseCookie(headers["Cookie"]); err == nil{
			client.Jar.SetCookies(url2, cookies)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		//Debug("http post do err",err.Error())
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		//Debug("resp.StatusCode",resp.StatusCode)
		return result
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Debug("http post readall err",err.Error())
	}else{
		result = string(body)
	}

	return
}

//发起一个get请求,获取字符串内容
func HttpGet(link string, values url.Values, headers map[string]string) (result string) {
	//Debug("link",link)
	if values != nil && len(values) > 0 {
		if strings.Contains(link,"?") {
			link += "?" + values.Encode()
		}else{
			link += "&" + values.Encode()
		}
	}
	req, err := http.NewRequest("GET", link, nil)

	if headers == nil {
		headers = map[string]string{}
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	jar,_ := cookiejar.New(nil)
	client := &http.Client{
		Jar:jar,
	}

	if headers["Cookie"] != "" {
		req.Header.Set("Set-Cookie",headers["Cookie"])

		var url2,_ = url.Parse(link)
		if cookies,err := ParseCookie(headers["Cookie"]); err == nil{
			//Debug("client.Jar",client.Jar)
			client.Jar.SetCookies(url2, cookies)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result
	}

	body, _ := ioutil.ReadAll(resp.Body)
	result = string(body)

	return
}

// 将cookie字符串解析为*http.Cookie对象
func ParseCookie(cookie string) ([]*http.Cookie, error) {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", cookie))))
	if err != nil {
		return nil, err
	}
	cookies := req.Cookies()
	if len(cookies) == 0 {
		return nil, fmt.Errorf("no cookies")
	}
	return cookies, nil
}

// 将对象转为map对象,字段以小写字段显示,终端打印显示字段名
func ToJsonObject(v interface{}) (r interface{}) {
	if jsonStr, err := ToJson(v); err == nil {
		ParseJson(&r, jsonStr)
	}
	return
}

// json string生成
func ToJson(v interface{}) (result string,err error) {
	if bytes, err := json.Marshal(v); err == nil {
		result = string(bytes)
	}
	return
}

// 从json中反序列化对象(静态的数据结构)
func ParseJson(v interface{}, s string) (err error) {
	//error:invalid character '\n' in string literal
	s = strings.TrimSpace(s)
	err = json.Unmarshal([]byte(s), v)
	return
}


// 将动态的json解析为文档对象
// 参见:https://github.com/Jeffail/gabs
//使用示例:
//value, ok = doc.Path("outter.inner.value1").Data().(float64)
// value == 10.0, ok == true
//value, ok = doc.Search("outter", "inner", "value1").Data().(float64)
// value == 10.0, ok == true
//value, ok = doc.Path("does.not.exist").Data().(float64)
// value == 0.0, ok == false
//exists := doc.Exists("outter", "inner", "value1")
// exists == true
//exists := doc.Exists("does", "not", "exist")
// exists == false
func ParseJsonToDocument(s string) (doc *gabs.Container, err error) {

	doc, err = gabs.ParseJSON([]byte(s))
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
