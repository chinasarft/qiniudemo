package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	context "golang.org/x/net/context"
	kodo "qiniupkg.com/api.v7/kodo"
)

type KeyPair struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

var keyPair KeyPair

//用hkconv这个ufop把key=2016-04-18.avi文件的start到end转为mp4放到bucket/aGlrdmlzaW9u/key/MjAxNjA5MDgxMDQwLm1wNA==
//bucket=hikvision&key=2016-04-18.avi&force=1&fops=hkconv/bucket/aGlrdmlzaW9u/key/MjAxNjA5MDgxMDQwLm1wNA==/start/10000/end/80000
//如果pfop以如上方式发起
//force这个参数是不会传递到ufop的
type ReqArgs struct {
	Cmd  string `json:"cmd"`  //就是fops=的内容
	Mode uint32 `json:"mode"` //标志是否为异步,ufop里几乎也不用理这个字段
	Src  struct {
		Url      string `json:"url"`      //自动把bucket=hikvision&key=2016-04-18.avi这个转为一个http的url地址
		Mimetype string `json:"mimetype"` //key=2016-04-18.avi文件的类型
		Fsize    int32  `json:"fsize"`    //key=2016-04-18.avi这个文件的大小
		Bucket   string `json:"bucket"`   //bucket=hikvision
		Key      string `json:"key"`      //key=2016-04-18.avi
	} `json: "src"`
}

func init() {
	bytes, err := ioutil.ReadFile("key.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(bytes, &keyPair)
	fmt.Println(keyPair)
	kodo.SetMac(keyPair.Ak, keyPair.Sk)
}

func serverUpload() {
	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象

	bucket := c.Bucket("lmkbucket")
	ctx := context.Background()
	localFile := "tmpfile111"
	err := bucket.PutFile(ctx, nil, "ufop/test.txt", localFile, nil)
	if err != nil {
		// 上传文件失败处理
		return
	}

}

func demoHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("start")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("fail to read req body")
		return
	}

	var args ReqArgs
	err = json.Unmarshal(body, &args)
	if err != nil {
		w.WriteHeader(500)
		log.Println("fail to unmarshal:", err)
		return
	}

	resp, err := http.Get(args.Src.Url)
	if err != nil {
		w.WriteHeader(500)
		log.Println("fail to fetch file")
		return
	}
	defer resp.Body.Close()
	fout, err := os.Create("tmpfile111")
	if err != nil {
		w.Write([]byte("open file tmpfile111 fail"))
		return

	}
	defer fout.Close()

	buf := make([]byte, 1024)
	for {
		size, _ := resp.Body.Read(buf)
		if size == 0 {
			break
		} else {
			w.Write(buf[:size])
			fout.Write(buf[:size])
		}
	}
	fout.Write(body)
	w.Write(body)

	zone := 0
	c := kodo.New(zone, nil) // 创建一个 Client 对象

	bucket := c.Bucket("lmkbucket")
	ctx := context.Background()
	localFile := "tmpfile111"
	err = bucket.PutFile(ctx, nil, "ufop/test.txt", localFile, nil)
	if err != nil {
		// 上传文件失败处理
		fmt.Println(err)
		w.Write([]byte("upload fail"))
	}

	w.Write([]byte("\nnew version"))
}

func main() {
	http.HandleFunc("/uop", demoHandler)
	err := http.ListenAndServe(":9100", nil)
	if err != nil {
		log.Fatal("Demo server failed to start:", err)
	}
}
