package product

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"walmart_web/app/activity"
	"walmart_web/app/config"
	"walmart_web/app/keyword"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
	"walmart_web/app/walLog"
)

var TaskKeysNum int
var TaskUrlNum int
var TaskActsNum int
var TaskIdsNum int
var ChGtin = make(chan int, 30)

func init() {
	TaskIds = config.IdsFile
	TaskKeys = config.KeyssFile
}

func CrawlerActivity(urlll string) {
	TaskActsNum++
	defer func() {
		TaskActsNum--
		<-config.Ch
		config.ActWg.Done()
	}()
	var page = 1
	var xc = 1
	lo := mode.Log{}
	var title string
	var idss []string

	for page <= 25 && xc <= 20 {
		if !IsRun {
			time.Sleep(60 * time.Second)
			continue
		}
		if xc != 1 {
			time.Sleep(1 * time.Second)
		}
		xc++

		urll := ""
		if page != 1 {
			if strings.Contains(urlll, "?") {
				urll = urlll + "&page=" + strconv.Itoa(page)
			} else {
				urll = urlll + "?page=" + strconv.Itoa(page)
			}
		} else {
			urll = urlll
		}
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		req, _ := http.NewRequest("GET", urll, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "zh")
		req.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Set("Sec-Fetch-Dest", `document`)
		req.Header.Set("Sec-Fetch-Mode", `navigate`)
		req.Header.Set("Sec-Fetch-Site", `none`)
		req.Header.Set("Sec-Fetch-User", `?1`)
		req.Header.Set("Upgrade-Insecure-Requests", `1`)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		if config.IsC2 {
			req.Header.Set("Cookie", tools.GenerateRandomString(10))
		}
		response, err := client.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
				log.Println("代理IP无效，开始重试")
				lo.Msg = "验证代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				lo.Msg = err.Error()

			}
			lo.Classify = "activity"
			lo.Val = urll
			continue
		}
		defer response.Body.Close()
		result := ""
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("出现错误，重新开始：" + urll)
				lo.Classify = "activity"
				lo.Msg = "解析body错误：" + err.Error()
				lo.Val = urll
				continue
			}
			defer reader.Close()
			con, err := io.ReadAll(reader)
			if err != nil {
				log.Println("gzip解压错误，重新开始：" + urll)
				lo.Classify = "activity"
				lo.Msg = "gzip解压错误"
				lo.Val = urll
				continue
			}
			result = string(con)
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("请求超时，重新开始：" + urll)
				lo.Classify = "activity"
				lo.Msg = "请求超时"
				lo.Val = urll
				continue
			}
			result = string(dataBytes)
		}
		result = strings.ReplaceAll(result, `\u0026`, `&`)
		//log.Println(result)

		cw1 := regexp.MustCompile("(is not valid JSON)").FindAllStringSubmatch(result, -1)
		cw2 := regexp.MustCompile("(The requested URL was rejected. Please consult with your administrator)").FindAllStringSubmatch(result, -1)
		if len(cw1) > 0 || len(cw2) > 0 {
			log.Println("内容错误，跳过该活动")
			lo.Classify = "activity"
			lo.Msg = "内容错误，跳过该活动"
			lo.Val = urll
			break
		}
		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控，更换IP重新开始")
			fmt.Println(urll)
			lo.Classify = "activity"
			lo.Msg = "被风控:Robot or human?"
			lo.Val = urll
			config.IsC2 = !config.IsC2
			continue
		}
		if err != nil {
			log.Println("错误信息：" + err.Error())
		}

		resultStr := ""
		resultS := regexp.MustCompile("(items\":\\[.+?\\].+?layoutEnum)").FindAllStringSubmatch(result, -1)
		allString := regexp.MustCompile("(There were no search results for)").FindAllString(result, -1)
		if len(allString) > 0 {
			log.Println("活动:" + urll + " 第" + strconv.Itoa(page) + "页 无搜索结果")
			act := mode.Activity{Name: title, Date: time.Now().Format("2006-01-02")}
			for i := range idss {
				if i == 0 {
					act.Ids += idss[i]
					continue
				}
				act.Ids += "," + idss[i]
			}
			activity.AddActivity(act)
			return
		}
		if len(resultS) == 0 {
			log.Println("被风控，更换IP重新开始")
			lo.Classify = "activity"
			lo.Msg = "被风控:items\":\\[.+?\\].+?layoutEnum"
			lo.Val = urll
			continue
		} else {
			resultStr = resultS[0][1]
		}
		if title == "" {
			titles := regexp.MustCompile(`<title>([^>]+?)</title>`).FindAllStringSubmatch(result, -1)
			if len(titles) != 1 {
				titles = regexp.MustCompile(`"name":"(.+?)"}`).FindAllStringSubmatch(result, -1)
				if len(titles) != 1 {
					titles = regexp.MustCompile(`"name":"(.+?)"}`).FindAllStringSubmatch(result, -1)
				}
			}

			if len(titles) == 0 {
				log.Println("不存在活动标题")
				lo.Classify = "activity"
				lo.Msg = "不存在活动标题"
				lo.Val = urll
				return
			}
			replace := strings.Replace(strings.Replace(titles[0][1], ` - Walmart.com`, "", -1), `\u0026`, "&", -1)
			title = replace
		}
		id := regexp.MustCompile("usItemId\":\"([0-9]+?)\",\"[^c]").FindAllStringSubmatch(resultStr, -1)

		ids := make([]string, 0)
		for i := range id {
			ids = append(ids, id[i][1])
			idss = append(idss, id[i][1])
		}
		//最大分页
		var max int
		maxPage := regexp.MustCompile("\"maxPage\":([0-9]+?),").FindAllStringSubmatch(result, -1)
		if len(maxPage) != 0 {
			max, _ = strconv.Atoi(maxPage[0][1])
		}
		log.Println("活动:" + urll + " 第" + strconv.Itoa(page) + "页 完成 " + strconv.Itoa(len(id)) + "个")
		AddCrawlerId(ids)
		if page == max {
			log.Println("活动:" + urll + "到达页尾")
			act := mode.Activity{Name: title, Date: time.Now().Format("2006-01-02")}
			for i := range idss {
				if i == 0 {
					act.Ids += idss[i]
					continue
				}
				act.Ids += "," + idss[i]
			}
			activity.AddActivity(act)
			return
		}
		xc = 1
		page++
	}
	walLog.AddLog(lo)
}

func _CrawlerKey(keyw string) {
	TaskKeysNum++
	defer func() {
		TaskKeysNum--
		<-config.Ch
		config.KeyWg.Done()
	}()
	var page = 1
	var xc = 1
	lo := mode.Log{}
	var idss []string
	for page <= 25 && xc <= 16 {

		if func() bool {
			if !IsRun {
				time.Sleep(60 * time.Second)
				return false
			}
			if xc != 1 {
				time.Sleep(1 * time.Second)
			}
			xc++
			proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
			proxy, _ := url.Parse(proxy_str)
			var client *http.Client
			if config.IsC {
				client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
			} else {
				client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
			}
			var k = strings.Replace(url.QueryEscape(keyw), "%20", "+", -1)
			urll := ""
			if page != 1 {
				urll = k + "&page=" + strconv.Itoa(page)
			} else {
				urll = k
			}

			urll = strings.Replace(urll, "%20", "+", -1)
			request, err := http.NewRequest("GET", "https://www.walmart.com/search?q="+urll, nil)
			if err != nil {
				log.Println("请求错误:", err)
				lo.Classify = "key"
				lo.Msg = "请求错误:" + err.Error()
				lo.Val = keyw
				return false
			}

			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
			request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
			request.Header.Set("Accept-Encoding", "gzip, deflate, br")
			//request.Header.Set("Accept-Language", "zh")
			request.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
			request.Header.Set("Sec-Ch-Ua-Mobile", "?0")
			request.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
			request.Header.Set("Sec-Fetch-Dest", `document`)
			request.Header.Set("Sec-Fetch-Mode", `navigate`)
			request.Header.Set("Sec-Fetch-Site", `none`)
			request.Header.Set("Sec-Fetch-User", `?1`)
			request.Header.Set("Upgrade-Insecure-Requests", `1`)
			request.Header.Set("Accept-Encoding", "gzip, deflate, br")
			if config.IsC2 {
				request.Header.Set("Cookie", tools.GenerateRandomString(10))
			}
			response, err := client.Do(request)
			if err != nil {
				if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
					log.Println("代理IP无效，开始重试")
					lo.Msg = "验证代理IP无效"
				} else if strings.Contains(err.Error(), "441") {
					log.Println("代理超频！暂停10秒后继续...")
					time.Sleep(time.Second * 10)
					lo.Msg = "代理超频"
				} else if strings.Contains(err.Error(), "440") {
					log.Println("代理宽带超频！暂停5秒后继续...")
					time.Sleep(time.Second * 5)
					lo.Msg = "代理超频"
				} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
					log.Println("代理宽带超频！暂停5秒后继续...")
					time.Sleep(time.Second * 5)
					lo.Msg = "代理超频"
				} else {
					log.Println("错误信息：" + err.Error())
					lo.Msg = err.Error()
				}
				lo.Classify = "key"
				lo.Val = urll
				return false
			}

			defer response.Body.Close()
			//if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("出现错误，重新开始：" + keyw)
				lo.Classify = "key"
				lo.Msg = "解析body错误：" + err.Error()
				lo.Val = keyw
				return false
			}
			con, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				log.Println("gzip解压错误，重新开始：" + keyw)
				lo.Classify = "key"
				lo.Msg = "gzip解压错误"
				lo.Val = keyw
				return false
			}
			res := strings.Replace(string(con), `\u0026`, `&`, -1)
			log.Println("测试长度：", len(res))
			con = nil
			//} else {
			//	dataBytes, err := io.ReadAll(response.Body)
			//	if err != nil {
			//		log.Println("请求超时，重新开始：" + keyw)
			//		lo.Classify = "key"
			//		lo.Msg = "请求超时"
			//		lo.Val = keyw
			//		return false
			//	}
			//	result := string(dataBytes)
			//}
			//log.Println(result)
			cw1 := regexp.MustCompile("(is not valid JSON)").FindAllStringSubmatch(res, -1)
			cw2 := regexp.MustCompile("(The requested URL was rejected. Please consult with your administrator)").FindAllStringSubmatch(res, -1)
			if len(cw1) > 0 || len(cw2) > 0 {
				log.Println("搜索内容错误，跳过该标题：" + keyw)
				lo.Classify = "key"
				lo.Msg = "内容错误，跳过该标题"
				lo.Val = keyw
				return true
			}
			fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(res, -1)
			if len(fk) > 0 {
				log.Println("被风控，更换IP重新开始")
				lo.Classify = "key"
				lo.Msg = "被风控:Robot or human?"
				lo.Val = keyw
				config.IsC2 = !config.IsC2
				return false
			}
			//doc, err := htmlquery.Parse(strings.NewReader(result))
			if err != nil {
				log.Println("错误信息：" + err.Error())
			}
			//log.Println(keyw+page)
			//log.Println(result)
			var resultStr string
			resultS := regexp.MustCompile("(items\":\\[.+?\\].+?layoutEnum)").FindAllStringSubmatch(res, -1)
			allString := regexp.MustCompile("(There were no search results for)").FindAllString(res, -1)
			if len(allString) > 0 {
				log.Println("关键词:" + keyw + " 第" + strconv.Itoa(page) + "页 无搜索结果")
				k := mode.Keyword{Name: keyw}
				for i := range idss {
					if i == 0 {
						k.Ids += idss[i]
						return false
					}
					k.Ids += "," + idss[i]
				}
				keyword.AddKeyword(k)
				return true
			}
			if len(resultS) == 0 {
				log.Println("被风控，更换IP重新开始")
				lo.Classify = "key"
				lo.Msg = "被风控:items\":\\[.+?\\].+?layoutEnum"
				lo.Val = keyw
				return false
			} else {
				resultStr = resultS[0][1]
			}
			id := regexp.MustCompile("usItemId\":\"([0-9]+?)\",\"[^c]").FindAllStringSubmatch(resultStr, -1)

			ids := make([]string, 0)
			for i := range id {
				ids = append(ids, id[i][1])
				idss = append(idss, id[i][1])
			}
			//最大分页
			var max int
			maxPage := regexp.MustCompile("\"maxPage\":([0-9]+?),").FindAllStringSubmatch(res, -1)
			if len(maxPage) != 0 {
				max, _ = strconv.Atoi(maxPage[0][1])
			}
			log.Println("关键词:" + keyw + " 第" + strconv.Itoa(page) + "页 完成 " + strconv.Itoa(len(id)) + "个")
			AddCrawlerId(ids)
			if page == max {
				log.Println("关键词:" + keyw + "到达页尾")
				k := mode.Keyword{Name: keyw}
				for i := range idss {
					if i == 0 {
						k.Ids += idss[i]
						continue
					}
					k.Ids += "," + idss[i]
				}
				keyword.AddKeyword(k)
				return true
			}
			xc = 1
			page++
			return false
		}() {
			return
		}

	}
	walLog.AddLog(lo)
}
func CrawlerKey(keyw string) {
	TaskKeysNum++
	var result []string
	defer func() {
		TaskKeysNum--
		<-config.Ch
		config.KeyWg.Done()
	}()
	var page = 1
	var xc = 1
	lo := mode.Log{}
	var idss []string
	for page <= 25 && xc <= 16 {
		if !IsRun {
			time.Sleep(60 * time.Second)
			continue
		}
		if xc != 1 {
			time.Sleep(1 * time.Second)
		}
		xc++
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		var k = strings.Replace(url.QueryEscape(keyw), "%20", "+", -1)
		urll := ""
		if page != 1 {
			urll = k + "&page=" + strconv.Itoa(page)
		} else {
			urll = k
		}

		urll = strings.Replace(urll, "%20", "+", -1)
		request, err := http.NewRequest("GET", "https://www.walmart.com/search?q="+urll, nil)
		if err != nil {
			log.Println("请求错误:", err)
			lo.Classify = "key"
			lo.Msg = "请求错误:" + err.Error()
			lo.Val = keyw
			continue
		}

		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")
		//request.Header.Set("Accept-Language", "zh")
		request.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
		request.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		request.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		request.Header.Set("Sec-Fetch-Dest", `document`)
		request.Header.Set("Sec-Fetch-Mode", `navigate`)
		request.Header.Set("Sec-Fetch-Site", `none`)
		request.Header.Set("Sec-Fetch-User", `?1`)
		request.Header.Set("Upgrade-Insecure-Requests", `1`)
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")
		if config.IsC2 {
			request.Header.Set("Cookie", tools.GenerateRandomString(10))
		}
		response, err := client.Do(request)
		if err != nil {
			if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
				log.Println("代理IP无效，开始重试")
				lo.Msg = "验证代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				lo.Msg = err.Error()
			}
			lo.Classify = "key"
			lo.Val = urll
			continue
		}

		defer response.Body.Close()
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("出现错误，重新开始：" + keyw)
				lo.Classify = "key"
				lo.Msg = "解析body错误：" + err.Error()
				lo.Val = keyw
				continue
			}
			con, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				log.Println("gzip解压错误，重新开始：" + keyw)
				lo.Classify = "key"
				lo.Msg = "gzip解压错误"
				lo.Val = keyw
				continue
			}
			result = append(result, string(con))
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("请求超时，重新开始：" + keyw)
				lo.Classify = "key"
				lo.Msg = "请求超时"
				lo.Val = keyw
				continue
			}
			result = append(result, string(dataBytes))
		}
		//result = strings.ReplaceAll(result, `\u0026`, `&`)
		//log.Println(result)

		cw1 := regexp.MustCompile("(is not valid JSON)").FindAllStringSubmatch(result[0], -1)
		cw2 := regexp.MustCompile("(The requested URL was rejected. Please consult with your administrator)").FindAllStringSubmatch(result[0], -1)
		if len(cw1) > 0 || len(cw2) > 0 {
			log.Println("搜索内容错误，跳过该标题：" + keyw)
			lo.Classify = "key"
			lo.Msg = "内容错误，跳过该标题"
			lo.Val = keyw
			break
		}
		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result[0], -1)
		if len(fk) > 0 {
			log.Println("被风控，更换IP重新开始")
			lo.Classify = "key"
			lo.Msg = "被风控:Robot or human?"
			lo.Val = keyw
			config.IsC2 = !config.IsC2
			continue
		}
		//doc, err := htmlquery.Parse(strings.NewReader(result))
		if err != nil {
			log.Println("错误信息：" + err.Error())
		}
		//log.Println(keyw+page)
		//log.Println(result)
		resultStr := ""
		resultS := regexp.MustCompile("(items\":\\[.+?\\].+?layoutEnum)").FindAllStringSubmatch(result[0], -1)
		allString := regexp.MustCompile("(There were no search results for)").FindAllString(result[0], -1)
		if len(allString) > 0 {
			log.Println("关键词:" + keyw + " 第" + strconv.Itoa(page) + "页 无搜索结果")
			k := mode.Keyword{Name: keyw}
			for i := range idss {
				if i == 0 {
					k.Ids += idss[i]
					continue
				}
				k.Ids += "," + idss[i]
			}
			keyword.AddKeyword(k)
			return
		}
		if len(resultS) == 0 {
			log.Println("被风控，更换IP重新开始")
			lo.Classify = "key"
			lo.Msg = "被风控:items\":\\[.+?\\].+?layoutEnum"
			lo.Val = keyw
			continue
		} else {
			resultStr = resultS[0][1]
		}
		id := regexp.MustCompile("usItemId\":\"([0-9]+?)\",\"[^c]").FindAllStringSubmatch(resultStr, -1)

		ids := make([]string, 0)
		for i := range id {
			ids = append(ids, id[i][1])
			idss = append(idss, id[i][1])
		}
		//最大分页
		var max int
		maxPage := regexp.MustCompile("\"maxPage\":([0-9]+?),").FindAllStringSubmatch(result[0], -1)
		if len(maxPage) != 0 {
			max, _ = strconv.Atoi(maxPage[0][1])
		}
		log.Println("关键词:" + keyw + " 第" + strconv.Itoa(page) + "页 完成 " + strconv.Itoa(len(id)) + "个")
		AddCrawlerId(ids)
		if page == max {
			log.Println("关键词:" + keyw + "到达页尾")
			k := mode.Keyword{Name: keyw}
			for i := range idss {
				if i == 0 {
					k.Ids += idss[i]
					continue
				}
				k.Ids += "," + idss[i]
			}
			keyword.AddKeyword(k)
			return
		}
		xc = 1
		page++
	}
	walLog.AddLog(lo)
}

func CrawlerUrl(urlll string) {
	TaskUrlNum++
	defer func() {
		TaskUrlNum--
		<-config.Ch
		config.UrlWg.Done()
	}()
	var page = 1
	var xc = 1
	lo := mode.Log{}

	for page <= 25 && xc <= 20 {
		if !IsRun {
			time.Sleep(60 * time.Second)
			continue
		}

		if xc != 1 {
			time.Sleep(1 * time.Second)
		}
		xc++

		urll := ""
		if page != 1 {
			urll = urlll + "?page=" + strconv.Itoa(page)
		} else {
			urll = urlll
		}
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		req, _ := http.NewRequest("GET", urll, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "zh")
		req.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Set("Sec-Fetch-Dest", `document`)
		req.Header.Set("Sec-Fetch-Mode", `navigate`)
		req.Header.Set("Sec-Fetch-Site", `none`)
		req.Header.Set("Sec-Fetch-User", `?1`)
		req.Header.Set("Upgrade-Insecure-Requests", `1`)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		if config.IsC2 {
			req.Header.Set("Cookie", tools.GenerateRandomString(10))
		}
		response, err := client.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
				log.Println("代理IP无效，开始重试")
				lo.Msg = "验证代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				lo.Msg = err.Error()
			}
			lo.Classify = "url"
			lo.Val = urll
			continue
		}

		result := ""
		defer response.Body.Close()
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("出现错误，重新开始：" + urll)
				lo.Classify = "url"
				lo.Msg = "解析body错误：" + err.Error()
				lo.Val = urll
				continue
			}
			defer reader.Close()
			con, err := io.ReadAll(reader)
			if err != nil {
				log.Println("gzip解压错误，重新开始：" + urll)
				lo.Classify = "url"
				lo.Msg = "gzip解压错误"
				lo.Val = urll
				continue
			}
			result = string(con)
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println("请求超时，重新开始：" + urll)
				lo.Classify = "url"
				lo.Msg = "请求超时"
				lo.Val = urll
				continue
			}
			result = string(dataBytes)
		}

		//log.Println(result)
		result = strings.ReplaceAll(result, `\u0026`, `&`)
		cw1 := regexp.MustCompile("(is not valid JSON)").FindAllStringSubmatch(result, -1)
		cw2 := regexp.MustCompile("(The requested URL was rejected. Please consult with your administrator)").FindAllStringSubmatch(result, -1)
		if len(cw1) > 0 || len(cw2) > 0 {
			log.Println("内容错误，跳过该url：" + urll)
			lo.Classify = "url"
			lo.Msg = "内容错误，跳过该url"
			lo.Val = urll
			break
		}
		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控，更换IP重新开始")
			fmt.Println(urll)
			lo.Classify = "url"
			lo.Msg = "被风控:Robot or human?"
			lo.Val = urll
			config.IsC2 = !config.IsC2
			continue
		}
		//doc, err := htmlquery.Parse(strings.NewReader(result))
		if err != nil {
			log.Println("错误信息：" + err.Error())
		}
		//log.Println(keyword+page)
		//log.Println(result)
		resultStr := ""
		resultS := regexp.MustCompile("(items\":\\[.+?\\].+?layoutEnum)").FindAllStringSubmatch(result, -1)
		allString := regexp.MustCompile("(There were no search results for)").FindAllString(result, -1)
		if len(allString) > 0 {
			log.Println("url:" + urll + " 第" + strconv.Itoa(page) + "页 无搜索结果")
			return
		}
		if len(resultS) == 0 {
			log.Println("被风控，更换IP重新开始")
			lo.Classify = "url"
			lo.Msg = "被风控:items\":\\[.+?\\].+?layoutEnum"
			lo.Val = urll
			continue
		} else {
			resultStr = resultS[0][1]
		}
		id := regexp.MustCompile("usItemId\":\"([0-9]+?)\",\"[^c]").FindAllStringSubmatch(resultStr, -1)

		ids := make([]string, 0)
		for i := range id {
			ids = append(ids, id[i][1])
		}
		//最大分页
		var max int
		maxPage := regexp.MustCompile("\"maxPage\":([0-9]+?),").FindAllStringSubmatch(result, -1)
		if len(maxPage) != 0 {
			max, _ = strconv.Atoi(maxPage[0][1])
		}
		log.Println("url:" + urll + " 第" + strconv.Itoa(page) + "页 完成 " + strconv.Itoa(len(id)) + "个")
		AddCrawlerId(ids[:])
		if page == max {
			log.Println("url:" + urll + "到达页尾")
			return
		}
		xc = 1
		page++
	}
	walLog.AddLog(lo)
}

func CrawlerId(id string) {
	TaskIdsNum++
	var result string
	defer func() {
		result = ""
		<-config.Ch
		config.IdWg.Done()
		TaskIdsNum--
	}()
	var xc = 1
	lo := mode.Log{}
	for xc <= 16 {
		result = ""
		if xc != 1 {
			time.Sleep(2 * time.Second)
		}
		xc++
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		request, err := http.NewRequest("GET", "https://www.walmart.com/ip/"+id, nil)
		if err != nil {
			log.Println("请求错误：", err)
			lo.Classify = "id"
			lo.Msg = "请求错误:" + err.Error()
			lo.Val = id
			continue
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		request.Header.Set("Accept-Encoding", "gzip, deflate, br")
		//request.Header.Set("Accept-Language", "zh")
		request.Header.Set("Sec-Ch-Ua", `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`)
		request.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		request.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		request.Header.Set("Sec-Fetch-Dest", `document`)
		request.Header.Set("Sec-Fetch-Mode", `navigate`)
		request.Header.Set("Sec-Fetch-Site", `none`)
		request.Header.Set("Sec-Fetch-User", `?1`)
		request.Header.Set("Upgrade-Insecure-Requests", `1`)
		if config.IsC {
			request.Header.Set("Cookie", tools.GenerateRandomString(10))
		}
		response, err := client.Do(request)

		if err != nil {
			if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
				log.Println("代理IP无效,重新开始：" + id)
				lo.Msg = "代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				lo.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				log.Println("出现错误，如果同id连续出现请联系我，重新开始：" + id)
				lo.Msg = err.Error()
			}
			lo.Classify = "id"
			lo.Val = id
			continue
		}
		defer response.Body.Close()
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("解析body错误，重新开始：" + id)
				lo.Classify = "id"
				lo.Msg = "解析body错误：" + err.Error()
				lo.Val = id
				continue
			}
			con, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				log.Println("gzip解压错误，重新开始：" + id)
				lo.Classify = "id"
				lo.Msg = "gzip解压错误"
				lo.Val = id
				continue
			}
			result = string(con)
			con = nil
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
					log.Println("代理IP无效,重新开始：" + id)
					lo.Msg = "代理IP无效"
				} else {
					log.Println("错误信息：" + err.Error())
					log.Println("出现错误，如果同id连续出现请联系我，重新开始：" + id)
					lo.Msg = err.Error()
				}
				lo.Classify = "id"
				lo.Val = id
				continue
			}
			result = string(dataBytes)
		}
		result = strings.ReplaceAll(result, `\u0026`, `&`)

		pro := mode.ProductDetails{}
		pro.Id, _ = strconv.Atoi(id)
		if strings.Contains(result, "This page could not be found.") {
			log.Println("id:" + id + "商品不存在")
			pro.Remark = "商品不存在"
			UploadProductNoIs(pro)
			return
		}

		//upc与upc类型
		upc := regexp.MustCompile("upc\":\"(.{4,30}?)\"").FindAllStringSubmatch(result, -1)
		gtin := regexp.MustCompile("gtin13\":\"(.{4,30}?)\"").FindAllStringSubmatch(result, -1)
		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		pro.CodeType = "gtin"
		if len(gtin) > 0 {
			pro.Code = gtin[0][1]
		} else if len(upc) > 0 {
			pro.Code = upc[0][1]
		} else if len(fk) > 0 {
			log.Println("id:" + id + " 被风控,更换IP继续")
			lo.Classify = "id"
			lo.Msg = "被风控：Robot or human?"
			lo.Val = id
			config.IsC = !config.IsC
			continue
		} else {
			//log.Println("id:"+id+" 获取为空，默认为ean")
		}

		doc, err := htmlquery.Parse(strings.NewReader(result))
		if err != nil {
			log.Println("错误信息：" + err.Error())
			lo.Classify = "id"
			lo.Msg = "未知错误，直接停止：" + err.Error()
			lo.Val = id
			return
		}
		result = strings.Replace(result, "\\u0026", "&", -1)

		//品牌
		brand := regexp.MustCompile("\"brand\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
		if len(brand) == 0 {
			//log.Println("品牌获取失败id："+id)
		} else {
			pro.Brands = brand[0][1]
		}

		//图片
		img := regexp.MustCompile("<meta property=\"og:image\" content=\"(.*?)\"/>").FindAllStringSubmatch(result, -1)
		if len(img) > 0 {
			pro.Img = img[0][1]
		}

		//标签
		query, err := htmlquery.QueryAll(doc, "//*[@id=\"maincontent\"]//div[@data-testid=\"sticky-buy-box\"]//div/p//span//text()")
		if err != nil {
			//log.Println("无标签")
		} else {
			queryStr := ""
			for _, v := range query {
				text := htmlquery.InnerText(v)
				if !strings.Contains(queryStr, text) {
					queryStr += text + " "
				}
			}
			pro.Tags = queryStr
		}
		if pro.Tags == "" {
			queryf := regexp.MustCompile("2\" aria-hidden=\"false\">(.*?)</span>").FindAllStringSubmatch(result, -1)
			queryStr := ""
			for _, v := range queryf {
				queryStr += v[1] + " "
			}
			pro.Tags = queryStr
		}

		//标题
		title := regexp.MustCompile("\"productName\":\"(.+?)\",").FindAllStringSubmatch(result, -1)
		if len(title) == 0 {
			log.Println("获取失败id："+id, "重新请求")

			lo.Classify = "id"
			lo.Msg = "获取id失败:\"productName\":\"(.+?)\""
			lo.Val = id
			continue
		} else {
			pro.Title = title[0][1]
		}

		//评分
		score := regexp.MustCompile("[(]([\\d][.][\\d])[)]").FindAllStringSubmatch(result, -1)
		if len(score) == 0 {
			//log.Println("评分获取失败id："+id)
		} else {
			pro.Rating = score[0][1]
		}

		//评论数量
		review := regexp.MustCompile("\"totalReviewCount\":(\\d+)").FindAllStringSubmatch(result, -1)
		if len(review) == 0 {
			//log.Println("评论数量获取失败id："+id)
		} else {
			pro.Comments = review[0][1]
		}

		//价格
		price := regexp.MustCompile("<span itemprop=\"price\".*?.{0,20}\\$([.\\d]+).{0,20}?</span>").FindAllStringSubmatch(result, -1)
		if len(price) == 0 {
			//log.Println("价格获取失败id："+id)
		} else {
			pro.Price, _ = strconv.ParseFloat(price[0][1], 64)
		}

		//卖家与配送
		//fulfilled := regexp.MustCompile(">Fulfilled by (.*?)</div>|>Fulfilled by .*?>(.*?)</a>?").FindAllStringSubmatch(result, -1)
		//sold := regexp.MustCompile(">Sold by ([^<]*?)</div>|>Sold by .*?>(.*?)</a>?").FindAllStringSubmatch(result, -1)
		//shipped := regexp.MustCompile("<div>Sold and shipped by ([^/]*?)</div>|<div>Sold and shipped by.*?>(.*?)</a>?").FindAllStringSubmatch(result, -1)	if len(fulfilled) != 0 && len(sold) != 0 {
		//	pro.Sellers = sold[0][1]
		//	pro.Distribution = fulfilled[0][1]
		//} else if len(shipped) != 0 {
		//	pro.Sellers = shipped[0][1]
		//	pro.Distribution = shipped[0][1]
		//}

		//卖家与配送
		all, err := htmlquery.QueryAll(doc, "//div/div/span[@class=\"lh-title\"]//text()")
		//log.Println(result)
		if err != nil {
			//log.Println("卖家与配送获取失败")
		} else {
			for i, v := range all {
				sv := htmlquery.InnerText(v)
				if strings.Contains(sv, "Sold by") {
					pro.Sellers = htmlquery.InnerText(all[i+1])
					continue
				}
				if strings.Contains(sv, "Fulfilled by") {
					pro.Distribution = strings.Replace(sv, "Fulfilled by ", "", -1)
					if len(pro.Distribution) < 3 && len(all) > i+1 {
						pro.Distribution = htmlquery.InnerText(all[i+1])
					}
					continue
				}
				if strings.Contains(sv, "Sold and shipped by") {
					pro.Sellers = htmlquery.InnerText(all[i+1])
					pro.Distribution = pro.Sellers
					break
				}
			}
		}

		//都为空的情况下
		if pro.Sellers == "" {
			seller := regexp.MustCompile("\"sellerDisplayName\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
			if len(seller) > 0 {
				pro.Sellers = seller[0][1]
			}
		}
		//}

		//配送时间
		deliveryDate := regexp.MustCompile("\"fulfillmentText\":\"(.+?)\"").FindAllStringSubmatch(result, -1)
		if len(deliveryDate) == 0 {
			//log.Println("配送时间获取失败id："+id)
		} else {
			pro.ArrivalTime = deliveryDate[0][1]
		}

		//变体
		variant := regexp.MustCompile(":</span><span aria-hidden=\"true\" class=\"ml1\">(.*?)</span>").FindAllStringSubmatch(result, -1)
		if len(variant) == 0 {
			//log.Println("评论数量获取失败id："+id)
		} else if len(variant) == 1 {
			pro.Variants1 = variant[0][1]
		} else if len(variant) == 2 {
			pro.Variants1 = variant[0][1]
			pro.Variants2 = variant[1][1]
		}

		allString := regexp.MustCompile("\",\"usItemId\":\"([0-9]+?)\"").FindAllStringSubmatch(result, -1)
		for i := range allString {
			if i == 0 {
				pro.VariantsId = fmt.Sprintf(allString[i][1])
			} else {
				pro.VariantsId += fmt.Sprintf("," + allString[i][1])

			}
		}
		paths := regexp.MustCompile("\"categoryPathName\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
		if len(paths) > 0 {
			spl := strings.Split(paths[0][1], "/")
			spl = spl[1:]
			var upName string
			var names = make(map[string]bool)
			var spli int
			for i := range spl {
				ii := i - spli
				if upName == spl[i] {
					spli++
					continue
				}
				if names[spl[i]] {
					break
				}
				names[spl[i]] = true
				if pro.CategoryName == "" {
					pro.CategoryName = spl[i]
				} else {
					pro.CategoryName += " -> " + spl[i]
				}
				switch ii {
				case 0:
					pro.Category1 = spl[i]
				case 1:
					pro.Category2 = spl[i]
				case 2:
					pro.Category3 = spl[i]
				case 3:
					pro.Category4 = spl[i]
				case 4:
					pro.Category5 = spl[i]
				case 5:
					pro.Category6 = spl[i]
				case 6:
					pro.Category7 = spl[i]
				}
				upName = spl[i]
			}
		}

		//startingFrom := regexp.MustCompile(`>Starting from \$([^<]+)<`).FindAllStringSubmatch(result, -1)
		//if len(startingFrom) == 0 {
		//	//log.Println("配送时间获取失败id："+id)
		//} else {
		//	pro.StarFrom = startingFrom[0][1]
		//}
		moreSellerOptions := regexp.MustCompile(`"additionalOfferCount":(\d+),`).FindAllStringSubmatch(result, -1)
		if len(moreSellerOptions) == 0 {
		} else {
			pro.StarFrom = moreSellerOptions[0][1]
		}
		AddProduct(pro)
		log.Println("id:", pro.Id, "完成")
		return
	}
	walLog.AddLog(lo)

}
