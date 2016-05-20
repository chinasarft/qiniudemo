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

type ReqArgs struct {
	Cmd  string `json:"cmd"`
	Mode uint32 `json:"mode"`
	Src  struct {
		Url      string `json:"url"`
		Mimetype string `json:"mimetype"`
		Fsize    int32  `json:"fsize"`
		Bucket   string `json:"bucket"`
		Key      string `json:"key"`
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
