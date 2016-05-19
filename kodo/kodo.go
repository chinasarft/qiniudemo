// kodoapi project main.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	context "golang.org/x/net/context"
	kodo "qiniupkg.com/api.v7/kodo"
	kodocli "qiniupkg.com/api.v7/kodocli"
)

func init() {
	fmt.Println("initlizing")
	//初始化ak pk
	kodo.SetMac("xx", "xx")
}

func myupload() {
	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象

	bucket := "lmkbucket"
	key := "test/test1.jpg"
	policy := &kodo.PutPolicy{
		Scope:   bucket + ":" + key, // 上传文件的限制条件，这里限制只能上传一个名为 "foo/bar.jpg" 的文件
		Expires: 3600,               // 这是限制上传凭证(uptoken)的过期时长，3600 是一小时
	}
	uptoken := c.MakeUptoken(policy) // 生成上传凭证
	fmt.Println(uptoken)

	zone = 0
	uploader := kodocli.NewUploader(zone, nil)
	ctx := context.Background()

	localFile := "d:/file/t01d947f4a4e932678a.jpg"
	err := uploader.PutFile(ctx, nil, uptoken, key, localFile, nil)
	if err != nil {
		fmt.Println("upload fail", err) // 上传文件失败处理
		return
	} else {
		fmt.Println("upload successful")
	}
}
func blockupload() {
	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象

	bucket := "lmkbucket"
	key := "test/testblk1.blk"
	policy := &kodo.PutPolicy{
		Scope:   bucket + ":" + key, // 上传文件的限制条件，这里限制只能上传一个名为 "foo/bar.jpg" 的文件
		Expires: 3600,               // 这是限制上传凭证(uptoken)的过期时长，3600 是一小时
	}
	uptoken := c.MakeUptoken(policy) // 生成上传凭证
	fmt.Println(uptoken)

	zone = 0
	uploader := kodocli.NewUploader(zone, nil)
	ctx := context.Background()

	localFile := "d:/testuploadtoqiniu.rar"
	err := uploader.RputFile(ctx, nil, uptoken, key, localFile, nil)
	if err != nil {
		fmt.Println("blockupload fail", err) // 上传文件失败处理
		return
	} else {
		fmt.Println("blockupload successful")
	}
}
func pulibcDownload() {
	domain := "7xpfbz.com1.z0.glb.clouddn.com"            // 您的空间绑定的域名，这个可以在七牛的Portal中查到
	baseUrl := kodo.MakeBaseUrl(domain, "test/test1.jpg") // 得到下载 url
	resp, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("download fail", err) // 上传文件失败处理
		return
	} else {
		fmt.Println("download successful")
		fmt.Printf("%v", resp)
	}
	file, _ := os.Create("d:/lmk/download/ttest1.jpg")
	io.Copy(file, resp.Body)
	fmt.Println("下载完成！")
}
func privateDownload() {
	domain := "7xpfbz.com1.z0.glb.clouddn.com"            // 您的空间绑定的域名，这个可以在七牛的Portal中查到
	baseUrl := kodo.MakeBaseUrl(domain, "test/test1.jpg") // 得到下载 url
	zone := 0
	c := kodo.New(zone, nil)                     // 创建一个 Client 对象
	privateUrl := c.MakePrivateUrl(baseUrl, nil) // 用默认的下载策略去生成私有下载的 url
	resp, err := http.Get(privateUrl)
	if err != nil {
		fmt.Println("download fail", err) // 上传文件失败处理
		return
	} else {
		fmt.Println("download successful")
		fmt.Printf("%v", resp)
	}
	file, _ := os.Create("d:/lmk/download/ttest1.jpg")
	io.Copy(file, resp.Body)
	fmt.Println("下载完成！")
}
func main() {
	//privateDownload()
	blockupload()
}
