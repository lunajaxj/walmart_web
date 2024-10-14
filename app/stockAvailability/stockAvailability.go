package stockAvailability

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/shopspring/decimal"
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
)

// 计算备货数量
func CalculateStockingQuantity(sas []mode.StockAvailability) []mode.StockAvailability {
	for i := range sas {
		weighted := decimal.NewFromFloat(sas[i].Weighted)
		leadTime := decimal.NewFromInt(sas[i].LeadTime)
		sas[i].Num = weighted.Mul(leadTime).IntPart() - sas[i].LibraryNum - sas[i].TransitNum
	}
	return sas
}

// 计算加权日均
func WeightedDailyAverage(sas map[string]map[string]int64) (map[string]float64, map[string]int64) {
	res := make(map[string]float64)
	res2 := make(map[string]int64)
	for k := range sas {
		x15 := decimal.NewFromInt(sas[k]["15"])
		x30 := decimal.NewFromInt(sas[k]["30"])
		x60 := decimal.NewFromInt(sas[k]["60"])
		g15 := decimal.NewFromInt(15)
		g30 := decimal.NewFromInt(30)
		g60 := decimal.NewFromInt(60)
		g03 := decimal.NewFromFloat(0.3)
		g05 := decimal.NewFromFloat(0.5)
		g02 := decimal.NewFromFloat(0.2)
		//15天销量/15*0.3+30天销量/30*0.5+60天销量/60*0.2
		f, err := strconv.ParseFloat(x15.Div(g15).Mul(g03).Add(x30.Div(g30).Mul(g05)).Add(x60.Div(g60).Mul(g02)).String(), 64)
		//fmt.Printf("id:%s,%s/15*0.3+%s/30*0.5+%s/60*0.2=%s\n\n", k, x15, x30, x60, f)
		if err != nil {
			log.Println("计算加权日均错误", err)
		}
		res[k] = f
		res2[k] = sas[k]["15"] + sas[k]["30"] + sas[k]["60"]
	}
	return res, res2
}

func DowImg(list []mode.StockAvailability) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var wg = sync.WaitGroup{}
	var ch = make(chan int, 5)
	for o := range list {
		ch <- 1
		wg.Add(1)
		go func(n int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			for i := 0; i < 16; i++ {
				if i != 0 {
					time.Sleep(time.Second * 1)
				}
				proxyUrl, _ := url.Parse("http://l752.kdltps.com:15818")
				tr.Proxy = http.ProxyURL(proxyUrl)
				basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("t16545052065610:wsad123456"))
				tr.ProxyConnectHeader = http.Header{}
				tr.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)

				client := &http.Client{Timeout: 10 * time.Second, Transport: tr}
				request, _ := http.NewRequest("GET", "https://www.walmart.com/ip/"+list[n].ItemId, nil)

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
					if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") {
						log.Println("代理IP无效，自动切换中")
						log.Println("连续出现代理IP无效请联系我，重新开始：" + list[n].ItemId)
						continue
					} else if strings.Contains(err.Error(), "441") {
						log.Println("代理超频！暂停10秒后继续...")
						time.Sleep(time.Second * 10)
						continue
					} else if strings.Contains(err.Error(), "440") {
						log.Println("代理宽带超频！暂停5秒后继续...")
						time.Sleep(time.Second * 5)
						continue
					} else {
						log.Println("错误信息：" + err.Error())
						log.Println("出现错误，如果同id连续出现请联系我，重新开始：" + list[n].ItemId)
						continue
					}
				}
				dataBytes, err := io.ReadAll(response.Body)
				if err != nil {
					if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") {
						log.Println("代理IP无效，自动切换中")
						log.Println("连续出现代理IP无效请联系我，重新开始：" + list[n].ItemId)
						continue
					} else {
						log.Println("错误信息：" + err.Error())
						log.Println("出现错误，如果同id连续出现请联系我，重新开始：" + list[n].ItemId)
						continue
					}
				}
				defer response.Body.Close()
				result := string(dataBytes)
				if strings.Contains(result, "This page could not be found.") {
					log.Println("id:" + list[n].ItemId + "商品不存在")
					return
				}

				fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)

				if len(fk) > 0 {
					log.Println("id:" + list[n].ItemId + " 被风控,更换IP继续")
					config.IsC = !config.IsC
					continue
				}
				img := regexp.MustCompile(`media-thumbnail"><img loading="lazy" srcset="([^,^"^?]+)`).FindAllStringSubmatch(result, -1)
				for i := range img {
					list[n].Img = img[i][1]
					break
				}
				log.Println(list[n].ItemId, "图片完成")

			}
		}(o)
	}
	wg.Wait()
	log.Println("图片全部获取完成")
	UpdateStockAvailabilityImg(list)

}
