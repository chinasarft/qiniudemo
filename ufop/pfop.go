package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
type CmdArgs struct {
	Bucket    string
	Key       string
	Force     string
	Fops      string
	NotifyURL string
	Pipeline  string

	PrintBody bool
	Urlencode bool

	Interval int
	Times    int

	TokenOnly bool
	Body      string
	Path      string
}
type PfopItem struct {
	Cmd       string `json:"id"`
	Code      int    `json:"code"`
	Desc      string `json:"desc"`
	Error     string `json:"error"`
	Hash      string `json:"hash"`
	Key       string `json:"key"`
	ReturnOld int    `json:"returnOld"`
}
type PfopStatus struct {
	Id          string     `json:"id"`
	Code        int        `json:"code"`
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
		log.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(bytes, &keyPair)
	log.Println(keyPair)
}

func (a *CmdArgs) getPostBody() (body string) {
	if a.Urlencode {
		/*
			body = "bucket=" + base64.URLEncoding.EncodeToString([]byte(a.Bucket)) +
				"&key=" + base64.URLEncoding.EncodeToString([]byte(a.Key)) +
				"&force=" + a.Force + "&fops=" + base64.URLEncoding.EncodeToString([]byte(a.Fops))
		*/
		body = "bucket=" + url.QueryEscape(a.Bucket) +
			"&key=" + url.QueryEscape(a.Key) +
			"&force=" + a.Force + "&fops=" + url.QueryEscape(a.Fops)
		if a.Pipeline != "" {
			body = body + "&pipeline=" + a.Pipeline
		}
		if a.NotifyURL != "" {
			//body = body + "&notifyURL=" + base64.URLEncoding.EncodeToString([]byte(a.NotifyURL))
			body = body + "&notifyURL=" + url.QueryEscape(a.NotifyURL)
		}

	} else {
		body = "bucket=" + a.Bucket + "&key=" + a.Key + "&force=" + a.Force + "&fops=" + a.Fops
		if a.Pipeline != "" {
			body = body + "&pipeline=" + a.Pipeline
		}
		if a.NotifyURL != "" {
			body = body + "&notifyURL=" + a.NotifyURL
		}
	}
	return
}

//token的生成
//http://developer.qiniu.com/article/developer/security/access-token.html
func getPfopToken(path string, body string) (token string) {
	log.Println("path:", path)
	key := []byte(keyPair.Sk)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(path + "\n" + body))
	log.Printf("%x\n", mac.Sum(nil))
	//b64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	b64 := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	log.Println(b64)
	token = keyPair.Ak + ":" + b64
	log.Println(" body:", body)
	log.Println("token:", token)
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
		log.Println(url, err)
		if resp == nil {
			log.Println("resp is nil")
			return nil
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read resp body")
		return nil
	}
	//log.Println(string(body))
	return body
}
func preparePost(a *CmdArgs) (res *PfopId) {
	var b []byte
	bstr := a.getPostBody()
	//b = post("http://api.qiniu.com/pfop/",  url.QueryEscape(bstr))
	b = post("http://api.qiniu.com/pfop/", bstr)

	var fop PfopId
	log.Println(string(b))
	err := json.Unmarshal(b, &fop)
	if err != nil {
		log.Println("Unmarshal fail:", err)
		return nil
	}
	res = &fop
	return
}

func get(id string) (st *PfopStatus) {
	resp, err := http.Get("http://api.qiniu.com/status/get/prefop?id=" + id)
	if err != nil {
		log.Println("http.Get fail:", err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	//	log.Println(string(body))
	var sts PfopStatus
	err = json.Unmarshal(body, &sts)
	if err != nil {
		log.Println(err)
		return nil
	}
	st = &sts
	return
}
func initArg(a *CmdArgs) {
	flag.StringVar(&a.Bucket, "bucket", "", "specify a bucket name")
	flag.StringVar(&a.Key, "key", "", "specify a key name")
	flag.StringVar(&a.Force, "force", "1", "ufop force")
	flag.StringVar(&a.Fops, "fops", "", "ufop fops")
	flag.StringVar(&a.NotifyURL, "notifyURL", "", "callback url")
	flag.StringVar(&a.Pipeline, "pipeline", "", "specify dedicated queue")

	flag.BoolVar(&a.PrintBody, "printBody", false, "defualt is false")
	flag.BoolVar(&a.Urlencode, "urlencode", true, "defualt is true")

	flag.IntVar(&a.Interval, "interval", 10, "querty PersistentId interval. 0 mean not querty")
	flag.IntVar(&a.Times, "times", 5, "querty PersistentId times")

	flag.BoolVar(&a.TokenOnly, "tokenonly", false, "just calculate token")
	flag.StringVar(&a.Body, "body", "", "use this to calculate token")
	flag.StringVar(&a.Path, "path", "", `use this to calculate token.
	sample:
	  path value:/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ=
	  ak: MY_ACCESS_KEY sk:MY_SECRET_KEY
	  result:MY_ACCESS_KEY:FXsYh0wKHYPEsIAgdPD9OfjkeEM=
	`)
	flag.Parse()
}
func checkArg(a *CmdArgs) {
	if a.TokenOnly {
		if a.Path == "" {
			log.Println("must give path value")
			os.Exit(1)
		}
		return
	}
	if a.Bucket == "" {
		log.Println("must specify a bucket name use:-bucket")
		os.Exit(1)
	}
	if a.Key == "" {
		log.Println("must specify a key(file) name use:-key")
		os.Exit(1)
	}
	if a.Fops == "" {
		log.Println("must specify fops(cmd) use:-fops")
		os.Exit(1)
	}
}

func main() {
	//post("http://api.qiniu.com/pfop/", "pipeline=jjj&bucket=hikvision&force=1&fops=hkconv%2Fbucket%2FaGlrdmlzaW9uYQ%3D%3D%2Fkey%2FaGlraW5nMTA0ODU3Ni5tcDQ%3D%2Fstart%2F1048576%2Fend%2F2379776&key=2016-04-18.avi")

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var arg CmdArgs
	initArg(&arg)
	checkArg(&arg)

	if arg.TokenOnly {
		getPfopToken(arg.Path, arg.Body)
		return
	}
	if arg.PrintBody {
		log.Println(arg.getPostBody())
		return
	}

	id := preparePost(&arg)
	if id == nil {
		log.Println("preparePost fail")
		return
	}
	if id.PersistentId == "" {
		log.Println(id.Error)
		return
	}

	log.Println(id.PersistentId)
	if arg.Interval > 0 && arg.Times > 0 {
		for i := 0; i < arg.Times; i++ {
			st := get(id.PersistentId)
			log.Println(st.Desc)
			time.Sleep(time.Second * (time.Duration(arg.Interval)))
		}
	}
}
