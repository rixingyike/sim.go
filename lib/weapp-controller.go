/*
 * 小程序用户自动登陆
 * author: liyi
 * email: 9830131#qq.com
 * date: 2017/4/14
 */
package sim

import (
	"gopkg.in/kataras/iris.v6"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type (
	WeappController struct {
		Web              *WebEngine
		AppId, AppSecret string
		token            WeixinAccessToken
	}
	WeappUser struct {
		Model  `xorm:"extends"`
		OpenID    string `xorm:"char(32)" json:"openId"`
		UnionID   string `xorm:"char(32)" json:"unionId"`
		Nickname  string `xorm:"char(50)" json:"nickName"`
		Gender    int    `xorm:"tinyint" json:"gender"`
		City      string `xorm:"char(20)" json:"city"`
		Province  string `xorm:"char(20)" json:"province"`
		Country   string `xorm:"char(20)" json:"country"`
		AvatarURL string `xorm:"tinytext" json:"avatarUrl"`
		Language  string `xorm:"char(10)" json:"language"`
		Watermark struct {
					  Timestamp int64  `json:"timestamp"`
					  AppID     string `json:"appid"`
				  } `xorm:"json" json:"watermark"`
	}
	// 微信json2session返回的数据结构
	WeappSession struct {
		SessionKey string `json:"session_key"`
		ExpiresIn  int `json:"expires_in"`
		Openid     string `json:"openid"`
	}

	WeixinAccessToken struct {
		AccessToken  string            `json:"access_token"`        // 网页授权接口调用凭证
		ExpiresIn    int64            `json:"expires_in"`           // access_token 接口调用凭证超时时间, 单位: 秒
		RefreshToken string            `json:"refresh_token"`       //刷新 access_token 的凭证
		OpenId       string            `json:"openid,omitempty"`
		Scope        string            `json:"scope,omitempty"`     //用户授权的作用域, 使用逗号(,)分隔

		ExpiresTime  time.Time        `json:"-"`

		CreatedAt    int64            `json:"created_at,omitempty"` //access_token 创建时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP
		UnionId      string            `json:"unionid,omitempty"`
	}
)

func (this *WeappController) Init() {
	this.Web.Get("/weapp/login", func(c *iris.Context) {
		var r Result
		var code, encryptedData, iv = this.getLoginArgsFromHeader(c.Request.Header)
		var session = this.retrieveSession(code)
		//Debug("weapp login args",code, encryptedData, iv, session)

		if user, err := this.decryptEncryptedData(encryptedData, iv, session.SessionKey); err == nil {
			if exist, err := this.Web.DB.Where("open_id = ?", user.OpenID).NoAutoCondition().Get(&user); err == nil {
				if !exist {
					this.Web.DB.InsertOne(&user)
				}
				r.Data = user
				r.Code = 1
			}else{
				r.Code = -2
			}
		}else{
			r.Code = -1
		}

		c.JSON(200, r)
	})
}


//以mediaId得到获取媒体文件的url
func (this *WeappController) ImageUrlForWeixinMediaId(mediaId string) string {
	var accessToken = this.getAccessToken()
	var wxImgUrl = fmt.Sprintf("http://file.api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",accessToken,mediaId)
	return wxImgUrl
}

//返回access token
func (this *WeappController) getAccessToken() (result string) {
	now := time.Now()

	if this.token.AccessToken != "" && this.token.ExpiresTime.After(now) {
		//如果未过期,直接返回
		result = this.token.AccessToken
	}else{
		//如果过期,先取再返
		var url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", this.AppId, this.AppSecret)
		// 曾经body变量名为result,结果出现了不易排查的错误,变量名称尽量有意义,与具体操作有关
		if body := HttpGet(url,nil,nil); body != "" {
			//DebugPrintln("get access token", body)
			ParseJson(&this.token, body)
			this.token.ExpiresTime = now.Add( time.Duration(this.token.ExpiresIn*1000) )
			result = this.token.AccessToken
		}
	}

	return
}

// 从请求头中获取小程序登陆的参数
func (this *WeappController) getLoginArgsFromHeader(header http.Header) (code, encryptedData, iv string) {
	code = header.Get("X-WX-Code")
	encryptedData = header.Get("X-WX-Encrypted-Data")
	iv = header.Get("X-WX-IV")
	return
}

// 拉取小程序登陆必须的session,依code拉取
func (this *WeappController) retrieveSession(code string) (weappSession WeappSession) {
	var url = fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", this.AppId, this.AppSecret, code);
	var res = HttpGet(url, nil, nil)
	//Debug("retrieveSession res",res)
	//map[session_key:vroPndZ1jeWdxkVAo5V05A== expires_in:7200 openid:o-hrq0EVYOTJHX9MWqk-LF-_KL0o]
	ParseJson(&weappSession, res)
	return
}

// 解密小程序加密信息
func (this *WeappController) decryptEncryptedData(encryptedData, iv, sessionKey string) (user WeappUser, err error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = this.pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return
	}
	err = json.Unmarshal(cipherText, &user)
	if err != nil {
		return
	}
	if user.Watermark.AppID != this.AppId {
		err = errors.New("app id not match")
		return
	}
	return
}

// pkcs7Unpad returns slice of the original data without padding
func (this *WeappController) pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errors.New("invalid block size")
	}
	if len(data) % blockSize != 0 || len(data) == 0 {
		return nil, errors.New("invalid PKCS7 data")
	}
	c := data[len(data) - 1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, errors.New("invalid padding on input")
	}
	for i := 0; i < n; i++ {
		if data[len(data) - n + i] != c {
			return nil, errors.New("invalid padding on input")
		}
	}
	return data[:len(data) - n], nil
}