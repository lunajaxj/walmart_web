package productSales

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/robfig/cron"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
	"walmart_web/app/walLog"
)

var Wg = sync.WaitGroup{}
var Mux sync.Mutex
var Num int
var cr = cron.New()

func init() {
	return
	if config.Mode == 2 {
		return
	}
	log.Println("加载库存销量定时任务...")
	err := cr.AddFunc("0 0 16 * * *", func() {
		saless, _ := SelectSales(1, 99999999, "")
		log.Printf("开始爬取库存更新销量 数量：%d", len(saless))
		Num = len(saless)
		for i := range saless {
			config.Ch <- 1
			Wg.Add(1)
			go CrawlerSales(saless[i])
		}
		Wg.Wait()
	})
	if err != nil {
		log.Println(err)
	} else {
		cr.Start()
	}
}

func CrawlerSales(se mode.ProductSales) {
	var result string
	defer func() {
		result = ""
		defer func() {
			<-config.Ch
			Num--
			Wg.Done()
		}()
	}()
	var xc = 1
	lo := mode.Log{}
	for xc <= 25 {
		result = ""
		if xc != 1 {
			time.Sleep(2 * time.Second)
		}
		xc++
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		id := strconv.Itoa(se.ITEMID)
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
		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控,更换IP继续")
			lo.Msg = "被风控"
			lo.Val = id
			config.IsC = !config.IsC
			continue
		}

		var sales mode.ProductSalesDay
		doc, err := htmlquery.Parse(strings.NewReader(result))
		if err != nil {
			log.Println("错误信息：" + err.Error())
			return
		}
		//卖家与配送
		all, err := htmlquery.QueryAll(doc, "//div/div/span[@class=\"lh-title\"]//text()")
		//log.Println(result)
		if err != nil {
			//log.Println("卖家与配送获取失败")
		} else {
			for i, v := range all {
				sv := htmlquery.InnerText(v)
				if strings.Contains(sv, "Sold by") {
					//sales.ShopName = htmlquery.InnerText(all[i+1])
					continue
				}
				if strings.Contains(sv, "Fulfilled by") {
					sales.Day0 = strings.Replace(sv, "Fulfilled by ", "", -1)
					if len(sales.Day0) < 3 && len(all) > i+1 {
						sales.Day0 = htmlquery.InnerText(all[i+1])
					}
					continue
				}
				if strings.Contains(sv, "Sold and shipped by") {
					sales.Day0 = htmlquery.InnerText(all[i+1])
					break
				}
			}
		}

		catalogSellerId := regexp.MustCompile("\"catalogSellerId\":(\\d+),").FindAllStringSubmatch(result, -1)
		if len(catalogSellerId) > 0 {
			sales.CatalogSellerId = catalogSellerId[0][1]
		}

		//图片
		var imgstr string
		img := regexp.MustCompile("<meta property=\"og:image\" content=\"(.*?)\"/>").FindAllStringSubmatch(result, -1)
		if len(img) > 0 {
			imgstr = img[0][1]
		}

		if se.CatalogSellerId != sales.CatalogSellerId {
			sales.Day0 = "有跟卖"
		} else if !strings.Contains(strings.ToLower(sales.Day0), "walmart") {
			sales.Day0 = "自发货"
		} else {
			inventory := regexp.MustCompile("availableQuantity\":(\\d+),").FindAllStringSubmatch(result, -1)
			if len(inventory) > 0 {
				sales.Day0 = inventory[0][1]
			}
		}
		log.Println("id:" + id + " 库存获取完成")
		sales.CatalogSellerId = se.CatalogSellerId
		sales.ITEMID = se.ITEMID
		Update(sales, imgstr)
		return
	}
	walLog.AddLog(lo)
}
