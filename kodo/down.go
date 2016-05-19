// kodoapi project main.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	kodo "qiniupkg.com/api.v7/kodo"
)

func init() {
	fmt.Println("initlizing")
	//初始化ak pk
        kodo.SetMac("xx", "xx")
}

func pulibcDownload() {
	domain := "7xpfbz.com1.z0.glb.clouddn.com"            // 您的空间绑定的域名，这个可以在七牛的Portal中查到
	baseUrl := kodo.MakeBaseUrl(domain, "test/test1.jpg") // 得到下载 url
	resp, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("download fail", err) //下载文件处理失败
		return
	} else {
		if resp.Status == "200 OK" {
			fmt.Println("download successful")
		} else {
			fmt.Println("http status:", resp.Status)
			fmt.Printf("%v", resp)
			fmt.Println("download fail", err) //下载文件处理失败
			return
		}
	}
	file, _ := os.Create("d:/lmk/download/ttest1.jpg")
	io.Copy(file, resp.Body)
	fmt.Println("下载完成！")
}
func privateDownload() {
	domain := "7xpfbz.com1.z0.glb.clouddn.com" // 您的空间绑定的域名，这个可以在七牛的Portal中查到
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
	privateDownload()
	//pulibcDownload()
}
