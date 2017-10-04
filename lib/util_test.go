package sim

import (
	// "fmt"
	"testing"
)

// curl -X POST -d '{}' https://api.douban.com/v2/user/ahbei
func TestNewRequest(t *testing.T) {
	url := "https://www.sogou.com"

	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Host":         "api.douban.com",
		"Content-Type": "application/json",
	}

	cookies := map[string]string{
		"userId":    "12",
		"loginTime": "15045682199",
	}

	queries := map[string]string{
		"page": "2",
		"act":  "update",
	}

	postData := map[string]interface{}{
		"name":      "mike",
		"age":       24,
		"interests": []string{"basketball", "reading", "coding"},
		"isAdmin":   true,
	}

	// 链式操作
	req := NewRequest()
	resp, err := req.
		SetUrl(url).
		SetHeaders(headers).
		SetCookies(cookies).
		SetQueries(queries).
		SetPostData(postData).Get()

	if err == nil {
		// fmt.Println(resp.Body)
		if resp.IsOk() {
			t.Log("ok")
		} else {
			t.Errorf("not is ok raw:%s", resp.Raw)
		}
	} else {
		t.Errorf("err %s", err)
	}
}

func TestDivision(t *testing.T) {
	if i := Division(1, 2); i != 12 {
		t.Error("除法函数测试没通过")
	} else {
		t.Log("第一个测试通过了")
	}
}
