package sim

import (
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"golang.org/x/net/context"
	"github.com/kataras/iris/v12"
	"io"
	"io/ioutil"
	"net/url"
	"log"
	"mime/multipart"
	"os"
	// "qiniupkg.com/x/url.v7"
	// "github.com/qiniu/api.v7"
)

//QINIU_IMAGE_SCALEMODE = "imageView2/2/w/360"
type QiniuController struct {
	Web                  *WebEngine
	Scope                string //bucket
	AccessKey, SecretKey string

	Watermark  string //图片的水印尾缀
	ServerBase string //七牛空间的基地址
	policy     *storage.PutPolicy
	mac        *qbox.Mac
	cfg        storage.Config
}

func (this *QiniuController) Init() {
	this.mac = qbox.NewMac(this.AccessKey, this.SecretKey)
	this.cfg = storage.Config{}
	this.cfg.Zone = &storage.ZoneHuadong
	this.cfg.UseHTTPS = false
	this.cfg.UseCdnDomains = false

	if this.ServerBase == "" {
		this.ServerBase = "http://7xndm1.com1.z0.glb.clouddn.com"
	}

	this.Web.Framework.Get("/qiniu/uptoken", func(c iris.Context) {
		c.JSON(H{
			"uptoken": this.newUptoken(),
		})
	})

	// 从微信接口上传取到的mediaid,上传至七牛空间,支持自定义key
	this.Web.Framework.Post("/qiniu/weixin/{mediaid}", func(c iris.Context) {
		var r Result
		var mediaId = c.Params().Get("mediaid")// c.Param("mediaid")
		var key = c.URLParam("key")

		if this.Web.Weapp != nil {
			var imageUrl = this.Web.Weapp.ImageUrlForWeixinMediaId(mediaId)
			if qiniuImgUrl := this.UploadImageFromUrl(imageUrl, key); qiniuImgUrl != "" {
				r.Code = 1
				r.Data = qiniuImgUrl
			}
		}

		c.JSON(r)
	})

	// 开启给simditor编辑器上传图片的接口
	this.Web.Framework.Post("/qiniu/simditor", func(c iris.Context) {
		_, info, _ := c.FormFile("imgfile") //iris v6
		file, _ := info.Open()
		defer file.Close()

		var url = this.UploadImageFromForm(file)
		var r = H{
			"success":   true,
			"msg":       "",
			"file_path": url,
		}
		c.JSON(r)
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

	var localFile = f.Name()

	putPolicy := storage.PutPolicy{
		Scope: this.Scope,
	}
	upToken := putPolicy.UploadToken(this.mac)

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&this.cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	err = formUploader.PutFileWithoutKey(context.Background(), &ret, upToken, localFile, &putExtra)
	if err != nil {
		path = ""
	} else {
		path = this.qiniuImageUrlFromKey(ret.Key)
	}
	return
}

//从本地上传文件对象至七牛
func (this *QiniuController) UploadImageFromLocale(f *os.File) (path string) {
	var localFile = f.Name()

	putPolicy := storage.PutPolicy{
		Scope: this.Scope,
	}
	upToken := putPolicy.UploadToken(this.mac)

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&this.cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	err := formUploader.PutFileWithoutKey(context.Background(), &ret, upToken, localFile, &putExtra)
	if err != nil {
		path = ""
	} else {
		path = this.qiniuImageUrlFromKey(ret.Key)
	}
	return
}

//从远程文件上传七牛云存储,key可以传空字符串
func (this *QiniuController) UploadImageFromUrl(url string, key string) string {
	var path string
	if key == "" {
		key = base64.StdEncoding.EncodeToString([]byte(url))
	}

	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(this.mac, &cfg)
	_, err := bucketManager.Fetch(url, this.Scope, key)
	if err != nil {
		fmt.Println("fetch error,", err)
	} else {
		path = this.qiniuImageUrlFromKey(key)
	}

	return path
}

func (this *QiniuController) qiniuImageUrlFromKey(key string) string {
	return fmt.Sprintf("%s/%s", this.ServerBase, url.QueryEscape(key))
}

func (this *QiniuController) newUptoken() string {
	putPolicy := storage.PutPolicy{
		Scope: this.Scope,
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	//生成一个上传token
	return putPolicy.UploadToken(this.mac)

}
