package main

import (
  "net/http"
  log "logging"
  "io/ioutil"
  "io"
  "compress/gzip"
  "bytes"
  "strings"
context "golang.org/x/net/context"
	kodo "qiniupkg.com/api.v7/kodo"
	kodocli "qiniupkg.com/api.v7/kodocli"
)

func init() {
	fmt.Println("initlizing")
	//初始化ak pk
	kodo.SetMac("l7zE60fO13jF2dW6csKAWIq-8XibFHBuzGvhLcQq", "7UISPqobJECpdWt3L0d6Yaq2t3XJhx-1SOPNuKTk")
}

func post(url string,postStr string) []byte{
 
  log.Debug("let's post :"+url)
 
  client := &http.Client{
    CheckRedirect: nil,
  }
 
  postBytesReader := bytes.NewReader([]byte(postStr))
  reqest, _ := http.NewRequest("POST", url, postBytesReader)
 
  reqest.Header.Set("Host","api.qiniu.com")
  reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//http://developer.qiniu.com/article/developer/security/access-token.html
  reqest.Header.Add("Authorization", "QBox AccessToken")
 
  resp, err := client.Do(reqest)
 
  if err != nil {
    log.Error(url,err)
    return nil
  }
 
  defer resp.Body.Close()
 
  var reader io.ReadCloser
  switch resp.Header.Get("Content-Encoding") {
  case "gzip":
    reader, err = gzip.NewReader(resp.Body)
    if err != nil {
      log.Error(url,err)
      return nil
    }
    defer reader.Close()
  default:
    reader = resp.Body
  }
 
 
  if(reader!=nil){
    body, err := ioutil.ReadAll(reader)
    if err != nil {
      log.Error(url,err)
      return nil
    }
    return body
  }
  return nil
}
