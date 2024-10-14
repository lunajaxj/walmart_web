package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var ProxyUrl = "l752.kdltps.com:15818"
var ProxyUser = "t19932187800946"
var ProxyPass = "wsad123456"

func main() {
	var a = 1
	var b = 2
	fmt.Println(float64(a) / float64(b))

	//proxy_str := fmt.Sprintf("http://%s:%s@%s", ProxyUser, ProxyPass, ProxyUrl)
	//proxy, _ := url.Parse(proxy_str)
	//		var client *http.Client
	if config.IsC {
		client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
	} else {
		client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	}
	//req, _ := http.NewRequest("GET", "https://walmart.com/ip/537921685", nil)
	//req.Header.Add("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	//req.Header.Add("Accept", "*/*")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	//
	//response, err := client.Do(req)
	//if err != nil {
	//	if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
	//		log.Println("代理IP无效，开始重试")
	//	} else if strings.Contains(err.Error(), "441") {
	//		log.Println("代理超频！暂停10秒后继续...")
	//		time.Sleep(time.Second * 10)
	//	} else if strings.Contains(err.Error(), "440") {
	//		log.Println("代理宽带超频！暂停5秒后继续...")
	//		time.Sleep(time.Second * 5)
	//	} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
	//		log.Println("代理宽带超频！暂停5秒后继续...")
	//		time.Sleep(time.Second * 5)
	//	} else {
	//		log.Println("错误信息：" + err.Error())
	//
	//	}
	//}
	//result := ""
	//if response.Header.Get("Content-Encoding") == "gzip" {
	//	reader, err := gzip.NewReader(response.Body) // gzip解压缩
	//	if err != nil {
	//		log.Println("出现错误，重新开始：", err.Error())
	//
	//	}
	//	defer reader.Close()
	//	con, err := io.ReadAll(reader)
	//	if err != nil {
	//		log.Println("gzip解压错误，重新开始：")
	//
	//	}
	//	result = string(con)
	//} else {
	//	dataBytes, err := io.ReadAll(response.Body)
	//	if err != nil {
	//		log.Println("请求超时，重新开始：")
	//
	//	}
	//	defer response.Body.Close()
	//	result = string(dataBytes)
	//}
	//
	//fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
	//if len(fk) > 0 {
	//	log.Println("被风控，更换IP重新开始")
	//
	//}
	//log.Println(result)
}
