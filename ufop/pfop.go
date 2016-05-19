package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"os"
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

func getPfopToken(path string, body string) (token string) {
	key := []byte(keyPair.Sk)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(path + "\n" + body))
	//fmt.Printf("%x\n", mac.Sum(nil))
	b64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	//fmt.Println(b64)
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

	req.Header.Set("Host", "api.qiniu.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//http://developer.qiniu.com/article/developer/security/access-token.html
	req.Header.Add("Authorization", " QBox "+getPfopToken("/pfop/", postStr))

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
func preparePost() (res *PfopId) {
	//pipeline=jjj&bucket=hikvision&force=1&fops=hkconv/bucket/aGlrdmlzaW9uYQ==/key/aGlraW5nMTA0ODU3Ni5tcDQ=/start/1048576/end/2379776&key=2016-04-18.avi
	//urlStr := "bucket=hikvision&force=1&fops=hkconv/bucket/aGlrdmlzaW9u/key/" + base64.StdEncoding.EncodeToString([]byte(os.Args[1])) + "/start/" + os.Args[2] + "/end/" + os.Args[3] + "&key=2016-04-18.avi"
	urlStr := "bucket=hikvision&force=1&fops=hkconv%2Fbucket%2FaGlrdmlzaW9u%2Fkey%2F" + base64.StdEncoding.EncodeToString([]byte(os.Args[1])) + "%2Fstart%2F" + os.Args[2] + "%2Fend%2F" + os.Args[3] + "&key=2016-04-18.avi"
	if len(os.Args) == 5 {
		urlStr = "pipeline=" + os.Args[4] + "&" + urlStr
	}
	fmt.Println(urlStr)
	//pfop只转义了/,所以不调用url.QueryEscape，调用会发生奇怪的错误，比如token
	//有时对有时错，并且有些时候参数有问题
	//urlenc := url.QueryEscape(urlStr)
	//b := post("http://api.qiniu.com/pfop/", urlenc)
	b := post("http://api.qiniu.com/pfop/", urlStr)

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
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
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
func main() {
	//post("http://api.qiniu.com/pfop/", "pipeline=jjj&bucket=hikvision&force=1&fops=hkconv%2Fbucket%2FaGlrdmlzaW9uYQ%3D%3D%2Fkey%2FaGlraW5nMTA0ODU3Ni5tcDQ%3D%2Fstart%2F1048576%2Fend%2F2379776&key=2016-04-18.avi")
	if len(os.Args) < 4 {
		fmt.Println("usage as:", os.Args[0], " mp4name startoffset endoffset [queue] ")
		return
	}
	id := preparePost()
	if id == nil {
		fmt.Println("preparePost fail")
		return
	}
	if id.PersistentId == "" {
		fmt.Println(id.Error)
		return
	}

	fmt.Println(id.PersistentId)
	for i := 0; i < 9; i++ {
		st := get(id.PersistentId)
		fmt.Println(st.Desc)
		time.Sleep(time.Second * 3)
	}
}
