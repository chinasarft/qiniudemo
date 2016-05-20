package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type KeyPair struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

var keyPair KeyPair

type PfopId struct {
	PersistentId string `json:"persistentId"`
	Error        string `json:"error"`
}
type CmdArgs struct {
	Mp4file   *string
	Queue     *string
	Start     *int
	End       *int
	Interval  *int
	Times     *int
	Urlencode *int
	TokenOnly *int
	Body      *string
	Path      *string
}
type PfopItem struct {
	cmd       string `json:"id"`
	code      int    `json:"code"`
	Desc      string `json:"desc"`
	Error     string `json:"error"`
	Hash      string `json:"hash"`
	Key       string `json:"key"`
	ReturnOld int    `json:"returnOld"`
}
type PfopStatus struct {
	Id          string     `json:"id"`
	code        int        `json:"code"`
	Desc        string     `json:"desc"`
	InputKey    string     `json:"inputKey"`
	InputBucket string     `json:"inputBucket"`
	Pipeline    string     `json:"pipeline"`
	Reqid       string     `json:"reqid"`
	Items       []PfopItem `json:"items"`
}

func init() {
	bytes, err := ioutil.ReadFile("key.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(bytes, &keyPair)
	fmt.Println(keyPair)
}

//token的生成
//http://developer.qiniu.com/article/developer/security/access-token.html
func getPfopToken(path string, body string) (token string) {
	fmt.Println("path:", path)
	key := []byte(keyPair.Sk)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(path + "\n" + body))
	fmt.Printf("%x\n", mac.Sum(nil))
	b64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(b64)
	token = keyPair.Ak + ":" + b64
	fmt.Println(" body:", body)
	fmt.Println("token:", token)
	return
}
func post(url string, postStr string) []byte {

	client := &http.Client{
		CheckRedirect: nil,
	}

	postBytesReader := bytes.NewReader([]byte(postStr))
	req, _ := http.NewRequest("POST", url, postBytesReader)

	req.Header.Add("Authorization", "QBox "+getPfopToken("/pfop/", postStr))
	req.Header.Set("Host", "api.qiniu.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("User-Agent", "QiniuJava/7.0.7 (Windows 7 amd64 6.1) Java 1.8.0_91")

	//req.Header.Add("Connection", "Keep-Alive")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(url, err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("fail to read resp body")
		return nil
	}
	//fmt.Println(string(body))
	return body
}
func preparePost(a *CmdArgs, useUrlencode int) (res *PfopId) {
	var b []byte
	if useUrlencode == 0 {
		//pfop只转义了/,所以不调用url.QueryEscape，调用会发生奇怪的错误，比如token
		//有时对有时错，并且有些时候参数有问题
		urlStr := "bucket=hikvision&force=1&fops=hkconv%2Fbucket%2FaGlrdmlzaW9u%2Fkey%2F" + base64.StdEncoding.EncodeToString([]byte(*a.Mp4file)) + "%2Fstart%2F" + strconv.Itoa(*a.Start) + "%2Fend%2F" + strconv.Itoa(*a.End) + "&key=2016-04-18.avi"
		if *a.Queue != "" {
			urlStr = "pipeline=" + (*a.Queue) + "&" + urlStr
		}
		fmt.Println(urlStr)
		b = post("http://api.qiniu.com/pfop/", urlStr)
	} else {
		urlStr := "bucket=hikvision&force=1&fops=hkconv/bucket/aGlrdmlzaW9u/key/" + base64.StdEncoding.EncodeToString([]byte(*a.Mp4file)) + "/start/" + strconv.Itoa(*a.Start) + "/end/" + strconv.Itoa(*a.End) + "&key=2016-04-18.avi"
		if *a.Queue != "" {
			urlStr = "pipeline=" + (*a.Queue) + "&" + urlStr
		}
		fmt.Println(urlStr)
		urlenc := url.QueryEscape(urlStr)
		b = post("http://api.qiniu.com/pfop/", urlenc)
	}

	var fop PfopId
	fmt.Println(string(b))
	err := json.Unmarshal(b, &fop)
	if err != nil {
		fmt.Println("Unmarshal fail:", err)
		return nil
	}
	res = &fop
	return
}

func get(id string) (st *PfopStatus) {
	resp, err := http.Get("http://api.qiniu.com/status/get/prefop?id=" + id)
	if err != nil {
		fmt.Println("http.Get fail:", err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//	fmt.Println(string(body))
	var sts PfopStatus
	err = json.Unmarshal(body, &sts)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	st = &sts
	return
}
func initArg(a *CmdArgs) {
	a.End = flag.Int("end", 0, "must specified")
	a.Start = flag.Int("start", 0, "must specified")
	a.Queue = flag.String("queue", "", "specify dedicated queue")
	a.Mp4file = flag.String("fname", "", "output mp4 file name")
	a.Interval = flag.Int("interval", 0, "querty PersistentId interval. 0 mean not querty")
	a.Times = flag.Int("times", 0, "querty PersistentId times")
	a.Urlencode = flag.Int("urlencode", 0, "defualt, just encode /")
	a.TokenOnly = flag.Int("tokenonly", 0, "just calculate token")
	a.Body = flag.String("body", "", "use this to calculate token")
	a.Path = flag.String("path", "", `use this to calculate token.
	sample:
	  path value:/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ=
	  ak: MY_ACCESS_KEY sk:MY_SECRET_KEY
	  result:MY_ACCESS_KEY:FXsYh0wKHYPEsIAgdPD9OfjkeEM=
	`)
	flag.Parse()
	fmt.Println("end     :", *a.End)
	fmt.Println("start   :", *a.Start)
	fmt.Println("interval:", *a.Interval)
	fmt.Println("times   :", *a.Times)
	fmt.Println("queue   :", *a.Queue)
	fmt.Println("mp4file :", *a.Mp4file)
	fmt.Println("urlenc  :", *a.Urlencode)
	fmt.Println("path    :", *a.Path)
}
func checkArg(a *CmdArgs) {
	if *a.TokenOnly != 0 {
		if *a.Path == "" {
			fmt.Println("must give path value")
			os.Exit(1)
		}
		return
	}
	if *a.End < *a.Start {
		fmt.Println("wrong start end")
		os.Exit(2)
	}
	if *a.Mp4file == "" {
		fmt.Println("fname must sepcified")
		os.Exit(2)
	}

}
func main() {
	//post("http://api.qiniu.com/pfop/", "pipeline=jjj&bucket=hikvision&force=1&fops=hkconv%2Fbucket%2FaGlrdmlzaW9uYQ%3D%3D%2Fkey%2FaGlraW5nMTA0ODU3Ni5tcDQ%3D%2Fstart%2F1048576%2Fend%2F2379776&key=2016-04-18.avi")
	var arg CmdArgs
	initArg(&arg)
	checkArg(&arg)
	if *arg.TokenOnly != 0 {
		getPfopToken(*arg.Path, *arg.Body)
		return
	}

	id := preparePost(&arg, *arg.Urlencode)
	if id == nil {
		fmt.Println("preparePost fail")
		return
	}
	if id.PersistentId == "" {
		fmt.Println(id.Error)
		return
	}

	fmt.Println(id.PersistentId)
	if *arg.Interval > 0 && *arg.Times > 0 {
		for i := 0; i < *arg.Times; i++ {
			st := get(id.PersistentId)
			fmt.Println(st.Desc)
			time.Sleep(time.Second * (time.Duration(*arg.Interval)))
		}
	}
}
