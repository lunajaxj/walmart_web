package storeInformation

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
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

var Num int
var Mux sync.Mutex
var MuxInfo sync.Mutex
var Wait sync.WaitGroup

var InfoNum int

type demo struct {
	CombinedColumnbrands string `db:"combined_columnbrands" comment:"combinedColumnbrands"`
}

func UpdateData(stos []mode.StoreInformation) {
	defer func() {
		Mux.Unlock()
		Num = 0
	}()
	//sql1 := `SELECT
	//   SUM(CASE WHEN sellers = ? THEN 1 ELSE 0 END) AS on_sale_product_count,
	//   SUM(CASE WHEN sellers = ? AND tags LIKE '%BEST SELLER%' THEN 1 ELSE 0 END) AS bsr_link,
	//   SUM(CASE WHEN sellers = ? AND tags LIKE '%POPPULUE%' THEN 1 ELSE 0 END) AS pp_link,
	//   SUM(CASE WHEN sellers = ? AND distribution LIKE '%walmart%' THEN 1 ELSE 0 END) AS wfs_delivery_count,
	//   SUM(CASE WHEN sellers = ? AND (price > 0 AND price <= 10) THEN 1 ELSE 0 END) AS price_0_10,
	//   SUM(CASE WHEN sellers = ? AND (price > 10 AND price <= 15) THEN 1 ELSE 0 END) AS price_10_15,
	//   SUM(CASE WHEN sellers = ? AND (price > 15 AND price <= 20) THEN 1 ELSE 0 END) AS price_15_20,
	//   SUM(CASE WHEN sellers = ? AND (price > 20 AND price <= 40) THEN 1 ELSE 0 END) AS price_20_40,
	//   SUM(CASE WHEN sellers = ? AND (price > 40 AND price <= 60) THEN 1 ELSE 0 END) AS price_40_60,
	//   SUM(CASE WHEN sellers = ? AND (price > 60) THEN 1 ELSE 0 END) AS price_60_above,
	//   SUM(CASE WHEN sellers = ? AND (comments > 0 AND comments <= 5) THEN 1 ELSE 0 END) AS reviews_0_5,
	//   SUM(CASE WHEN sellers = ? AND (comments > 5 AND comments <= 10) THEN 1 ELSE 0 END) AS reviews_5_10,
	//   SUM(CASE WHEN sellers = ? AND (comments > 10 AND comments <= 20) THEN 1 ELSE 0 END) AS reviews_10_20,
	//   SUM(CASE WHEN sellers = ? AND (comments > 20 AND comments <= 100) THEN 1 ELSE 0 END) AS reviews_20_100,
	//   SUM(CASE WHEN sellers = ? AND (comments > 100) THEN 1 ELSE 0 END) AS reviews_100_above
	//FROM
	//   product_details
	//WHERE
	//   sellers = ?;
	//`

	sql2 := `select  CONCAT(brands, '(', count(id), ')') AS combined_columnbrands from product_details where sellers = ? GROUP BY brands`
	Num = len(stos)
	for s := range stos {
		Num--
		var ls1 []mode.StoreInformation
		name := strings.Replace(stos[s].AccountShopName, `'`, `\'`, -1)
		err := config.Db.Select(&ls1, fmt.Sprintf(`SELECT
	   SUM(CASE WHEN sellers = '%s' THEN 1 ELSE 0 END) AS on_sale_product_count,
	   SUM(CASE WHEN sellers = '%s'  AND tags LIKE '%%BEST SELLER%%' THEN 1 ELSE 0 END) AS bsr_link,
	   SUM(CASE WHEN sellers = '%s' AND tags LIKE '%%Popular pick%%' THEN 1 ELSE 0 END) AS pp_link,
	   SUM(CASE WHEN sellers = '%s' AND distribution LIKE '%%walmart%%' THEN 1 ELSE 0 END) AS wfs_delivery_count,
	   SUM(CASE WHEN sellers = '%s' AND (price > 0 AND price <= 10) THEN 1 ELSE 0 END) AS price_0_10,
	   SUM(CASE WHEN sellers = '%s' AND (price > 10 AND price <= 15) THEN 1 ELSE 0 END) AS price_10_15,
	   SUM(CASE WHEN sellers = '%s' AND (price > 15 AND price <= 20) THEN 1 ELSE 0 END) AS price_15_20,
	   SUM(CASE WHEN sellers = '%s' AND (price > 20 AND price <= 40) THEN 1 ELSE 0 END) AS price_20_40,
	   SUM(CASE WHEN sellers = '%s' AND (price > 40 AND price <= 60) THEN 1 ELSE 0 END) AS price_40_60,
	   SUM(CASE WHEN sellers = '%s' AND (price > 60) THEN 1 ELSE 0 END) AS price_60_above,
	   SUM(CASE WHEN sellers = '%s' AND (comments > 0 AND comments <= 5) THEN 1 ELSE 0 END) AS reviews_0_5,
	   SUM(CASE WHEN sellers = '%s' AND (comments > 5 AND comments <= 10) THEN 1 ELSE 0 END) AS reviews_5_10,
	   SUM(CASE WHEN sellers = '%s' AND (comments > 10 AND comments <= 20) THEN 1 ELSE 0 END) AS reviews_10_20,
	   SUM(CASE WHEN sellers = '%s' AND (comments > 20 AND comments <= 100) THEN 1 ELSE 0 END) AS reviews_20_100,
	   SUM(CASE WHEN sellers = '%s' AND (comments > 100) THEN 1 ELSE 0 END) AS reviews_100_above
	FROM
	   product_details
	WHERE
	   sellers = '%s';`, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name))
		if err != nil {
			log.Println(name, s, err)
			continue
		}
		stos[s].OnSaleProductCount = ls1[0].OnSaleProductCount
		stos[s].BSRLink = ls1[0].BSRLink
		stos[s].PPLink = ls1[0].PPLink
		stos[s].HasTargetLink = ls1[0].PPLink + ls1[0].BSRLink
		stos[s].WFSDeliveryCount = ls1[0].WFSDeliveryCount
		stos[s].WFSDeliveryPercentage = float64(ls1[0].WFSDeliveryCount) / float64(ls1[0].OnSaleProductCount)
		stos[s].Price0_10 = ls1[0].Price0_10
		stos[s].Price10_15 = ls1[0].Price10_15
		stos[s].Price15_20 = ls1[0].Price15_20
		stos[s].Price20_40 = ls1[0].Price20_40
		stos[s].Price40_60 = ls1[0].Price40_60
		stos[s].Price60Above = ls1[0].Price60Above
		stos[s].Reviews0_5 = ls1[0].Reviews0_5
		stos[s].Reviews5_10 = ls1[0].Reviews5_10
		stos[s].Reviews10_20 = ls1[0].Reviews10_20
		stos[s].Reviews20_100 = ls1[0].Reviews20_100
		stos[s].Reviews100Above = ls1[0].Reviews100Above

		var ls2 []demo
		err = config.Db.Select(&ls2, sql2, stos[s].AccountShopName)
		if err != nil {
			log.Println(name, s, err)
			continue
		}
		for i := range ls2 {
			if i == 0 {
				stos[s].Brand = ls2[i].CombinedColumnbrands
			} else {
				stos[s].Brand += "," + ls2[i].CombinedColumnbrands
			}

		}
		UpdateStoreInformation([]mode.StoreInformation{stos[s]})

	}

}

func CrawlerInfo(id string) {
	var result string
	defer func() {
		result = ""
		defer func() {
			<-config.Ch
			InfoNum--
			Wait.Done()
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

		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		request, err := http.NewRequest("GET", "https://www.walmart.com/seller/"+id, nil)
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
		request.Header.Set("Cookie", tools.GenerateRandomString(10))
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

		sellerName := regexp.MustCompile(`"sellerName":"(.+?)"`).FindAllStringSubmatch(result, -1)
		store := mode.StoreInformation{}
		atoi, _ := strconv.Atoi(id)
		store.ShopID = atoi

		if len(sellerName) > 0 {
			store.SellerName = sellerName[0][1]
		}
		var add string
		address1 := regexp.MustCompile(`"address1":"(.+?)"`).FindAllStringSubmatch(result, -1)
		address2 := regexp.MustCompile(`"address2":"(.+?)"`).FindAllStringSubmatch(result, -1)
		city := regexp.MustCompile(`"city":"(.+?)"`).FindAllStringSubmatch(result, -1)
		state := regexp.MustCompile(`"state":"(.+?)"`).FindAllStringSubmatch(result, -1)
		postalCode := regexp.MustCompile(`"postalCode":"(.+?)"`).FindAllStringSubmatch(result, -1)
		country := regexp.MustCompile(`"country":"(.+?)"`).FindAllStringSubmatch(result, -1)

		plNum := regexp.MustCompile(`>\((\d+) reviews\)<`).FindAllStringSubmatch(result, -1)

		pf := regexp.MustCompile(`"averageOverallRating":(.*?),`).FindAllStringSubmatch(result, -1)
		if len(address1) > 0 {
			add += address1[0][1] + ","
		}
		if len(address2) > 0 {
			add += address2[0][1] + ","
		}
		if len(city) > 0 {
			add += city[0][1] + ","
		}
		if len(state) > 0 {
			add += state[0][1] + " "
		}
		if len(postalCode) > 0 {
			add += postalCode[0][1]
		}
		if len(country) > 0 {
			store.Country = country[0][1]
		}
		if len(plNum) > 0 {
			i, _ := strconv.Atoi(plNum[0][1])
			store.SellerReviewsNum = i
		}
		if len(pf) > 0 {
			float, _ := strconv.ParseFloat(pf[0][1], 64)
			store.Read = float
		}

		store.Address = add
		log.Println("id:" + id + " 卖家获取完成")
		UpdateInfo(store)
		return
	}
	walLog.AddLog(lo)
}
