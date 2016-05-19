// kodoapi project main.go
package main

import (
	"fmt"

	context "golang.org/x/net/context"
	kodo "qiniupkg.com/api.v7/kodo"
	kodocli "qiniupkg.com/api.v7/kodocli"
)

func init() {
	fmt.Println("initlizing")
	//初始化ak pk
        kodo.SetMac("xx", "xx")
}

func serverUpload() {
	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象

	bucket := c.Bucket("lmkbucket")
	ctx := context.Background()
	localFile := "D:\\tmp\\m64\\server20160317.exe"
	err := bucket.PutFile(ctx, nil, "ipcam/server20160317", localFile, nil)
	if err != nil {
		fmt.Println("putfile error")
		// 上传文件失败处理
		return
	}
	fmt.Println("putfile success")

}
func clientUpload() {
	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象
	bucket := "lmkbucket"
	key := "ipcam/server20160317"
	policy := &kodo.PutPolicy{
		Scope:   bucket + ":" + key, // 上传文件的限制条件，这里限制只能上传一个名为 "foo/bar.jpg" 的文件
		Expires: 3600,               // 这是限制上传凭证(uptoken)的过期时长，3600 是一小时
	}
	uptoken := c.MakeUptoken(policy) // 生成上传凭证
	fmt.Println(uptoken)

	zone = 0
	uploader := kodocli.NewUploader(zone, nil)
	ctx := context.Background()

	localFile := "D:/tmp/m64/server20160317.exe"
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

func main() {
	//serverUpload()
	clientUpload()
}
