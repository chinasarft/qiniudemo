package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

//官方例子
//http://developer.qiniu.com/article/developer/security/access-token.html
/*
# 假设有如下的管理请求：
AccessKey = "MY_ACCESS_KEY"
SecretKey = "MY_SECRET_KEY"
url = "http://rs.qiniu.com/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ="

#则待签名的原始字符串是：
signingStr = "/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ=\n"

#签名字符串是：
sign = "157b18874c0a1d83c4b0802074f0fd39f8e47843"
注意：签名结果是二进制数据，此处输出的是每个字节的十六进制表示，以便核对检查。

#编码后的签名字符串是：
encodedSign = "FXsYh0wKHYPEsIAgdPD9OfjkeEM="

#最终的管理凭证是：
accessToken = "MY_ACCESS_KEY:FXsYh0wKHYPEsIAgdPD9OfjkeEM="
*/
var (
	ak = "MY_ACCESS_KEY"
	sk = "MY_SECRET_KEY"
)

func test() {
	fmt.Println("pfop和存储都是用的这个算法")
	fmt.Println("requst url:http://rs.qiniu.com/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ=")

	signingStr := "/move/bmV3ZG9jczpmaW5kX21hbi50eHQ=/bmV3ZG9jczpmaW5kLm1hbi50eHQ=\n"
	key := []byte(sk)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(signingStr))
	fmt.Printf("%x\n", mac.Sum(nil)) //157b18874c0a1d83c4b0802074f0fd39f8e47843

	b64 := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(b64) //FXsYh0wKHYPEsIAgdPD9OfjkeEM=
	fmt.Println("token:", ak+":"+b64)
}
func main() {
	test()
}
