package sim

import (
	"io"
	"os"
	"mime/multipart"
	"io/ioutil"
	"log"
	"github.com/qiniu/api.v7/kodo"
	"qiniupkg.com/x/url.v7"
	"gopkg.in/kataras/iris.v6"
	"fmt"
	"encoding/base64"
)

//QINIU_IMAGE_SCALEMODE = "imageView2/2/w/360"
type QiniuController struct {
	Web              *WebEngine
	Scope     string
	AccessKey, SecretKey string

	Watermark string //图片的水印尾缀
	ServerBase    string //七牛空间的基地址
	client    *kodo.Client
	policy *kodo.PutPolicy
}

func (this *QiniuController) Init() {
	kodo.SetMac(this.AccessKey, this.SecretKey)
	//创建一个Client
	var zone = 0 //空间Bucket所在的区域
	this.client = kodo.New(zone, nil)

	if this.ServerBase == "" {
		this.ServerBase = "http://7xndm1.com1.z0.glb.clouddn.com"
	}

	this.Web.Get("/qiniu/uptoken", func(c *iris.Context) {
		c.JSON(200, H{
			"uptoken":this.newUptoken(),
		})
	})

	// 从微信接口上传取到的mediaid,上传至七牛空间,支持自定义key
	this.Web.Post("/qiniu/weixin/:mediaid", func(c *iris.Context) {
		var r Result
		var mediaId = c.Param("mediaid")
		var key = c.URLParam("key")

		if this.Web.Weapp != nil {
			var imageUrl = this.Web.Weapp.ImageUrlForWeixinMediaId(mediaId)
			if qiniuImgUrl := this.UploadImageFromUrl(imageUrl,key); qiniuImgUrl != "" {
				r.Code = 1
				r.Data = qiniuImgUrl
			}
		}

		c.JSON(200, r)
	})

	// 开启给simditor编辑器上传图片的接口
	this.Web.Post("/qiniu/simditor", func(c *iris.Context) {
		_,info, _ := c.FormFile("imgfile") //iris v6
		file, _ := info.Open()
		defer file.Close()

		var url = this.UploadImageFromForm(file)
		var r = H{
			"success":true,
			"msg":"",
			"file_path":url,
		}
		c.JSON(200, r)
	})
}

//上传文件对象至七牛
func (this *QiniuController) UploadImageFromForm(file multipart.File) (path string) {
	f, err := ioutil.TempFile(os.TempDir(), "qiniu")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	defer os.Remove(f.Name())

	var localfile = f.Name()
	var uploader = this.client.Bucket(this.Scope)
	var ret kodo.PutRet
	var extra = &kodo.PutExtra{
		CheckCrc: 0,
	}
	res := uploader.PutFileWithoutKey(nil, &ret, localfile, extra)
	if res != nil {
		log.Println("io.Put failed:", res)
		path = ""
	} else {
		path = this.qiniuImageUrlFromKey(ret.Key)
	}
	return
}

//从本地上传文件对象至七牛
func (this *QiniuController) UploadImageFromLocale(f *os.File) (path string) {
	var localfile = f.Name()
	var uploader = this.client.Bucket(this.Scope)
	var ret kodo.PutRet
	var extra = &kodo.PutExtra{
		CheckCrc: 0,
	}
	res := uploader.PutFileWithoutKey(nil, &ret, localfile, extra)
	if res != nil {
		log.Println("io.Put failed:", res)
		path = ""
	} else {
		path = this.qiniuImageUrlFromKey(ret.Key)
	}
	return
}

//从远程文件上传七牛云存储,key可以传空字符串
func (this *QiniuController) UploadImageFromUrl(url string,key string) string {
	var path string
	var uploader = this.client.Bucket(this.Scope)

	if key == "" {
		key = base64.StdEncoding.EncodeToString([]byte(url))
	}

	var err = uploader.Fetch(nil, key, url)
	if err != nil {
		Debug("io.Put failed:", err.Error())
	} else {
		path = this.qiniuImageUrlFromKey(key)
	}
	return path
}

func (this *QiniuController) qiniuImageUrlFromKey(key string) string {
	return fmt.Sprintf("%s/%s",this.ServerBase,url.Escape(key))
}

func (this *QiniuController) newUptoken() string {
	//设置上传的策略
	if this.policy == nil {
		this.policy = &kodo.PutPolicy{
			Scope:   this.Scope,
			//设置Token过期时间
			Expires: 3600,
		}
	}
	//生成一个上传token
	return this.client.MakeUptoken(this.policy);
}