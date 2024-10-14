package shopping

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/gammazero/workerpool"
	"github.com/robfig/cron"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/context"
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

// test
// var users = map[string][]string{"demo": {"xcitnpkj@126.com", "*IcI4$fz152NFjI"}}

// 生产环境
// 添加用户，就在这里添加
var users = map[string][]string{"Online mini-mart": {"yumumandzsw@163.com", "Dyumuman123!"}, "GUBIN": {"guibinwalmart@163.com", "DGuibin123!"}, "Money Saving Center": {"ketianwalmart@163.com", "DKetian123!"}, "Money Saving World": {"ziyuewlkj@163.com", "Dziyue123!"}, "Kiaote Center": {"qiaotewalmart@163.com", "Dqiaote123!"}, "Moremuma": {"moremuma@163.com", "$Moremuma778"}, "Juno Kael": {"keyshinewalmart@163.com", "&Keyshine788"}}
var se = map[string]string{"GUBIN": "8333", "Money Saving Center": "9333", "Money Saving World": "7333", "Kiaote Center": "10222", "Moremuma": "11222", "Juno Kael": "12222"}
var IsRun = map[string]bool{"Online mini-mart": false, "GUBIN": false, "Money Saving Center": false, "Money Saving World": false, "Kiaote Center": false, "Moremuma": false, "Juno Kael": false}
var Mutexs = map[string]*sync.Mutex{"Online mini-mart": new(sync.Mutex), "GUBIN": new(sync.Mutex), "Money Saving Center": new(sync.Mutex), "Money Saving World": new(sync.Mutex), "Kiaote Center": new(sync.Mutex), "Moremuma": new(sync.Mutex), "Juno Kael": new(sync.Mutex)}
var wg = map[string]*sync.WaitGroup{"Online mini-mart": {}, "GUBIN": {}, "Money Saving Center": {}, "Money Saving World": {}, "Kiaote Center": {}, "Moremuma": {}, "Juno Kael": {}}
var wg2 = map[string]*sync.WaitGroup{"Online mini-mart": {}, "GUBIN": {}, "Money Saving Center": {}, "Money Saving World": {}, "Kiaote Center": {}, "Moremuma": {}, "Juno Kael": {}}
var wg3 = map[string]*sync.WaitGroup{"Online mini-mart": {}, "GUBIN": {}, "Money Saving Center": {}, "Money Saving World": {}, "Kiaote Center": {}, "Moremuma": {}, "Juno Kael": {}}
var crs = map[string]*cron.Cron{"Online mini-mart": cron.New(), "GUBIN": cron.New(), "Money Saving Center": cron.New(), "Money Saving World": cron.New(), "Kiaote Center": cron.New(), "Moremuma": cron.New(), "Juno Kael": cron.New()}
var IsRun2 = map[string]bool{"Online mini-mart": false, "GUBIN": false, "Money Saving Center": false, "Money Saving World": false, "Kiaote Center": false, "Moremuma": false, "Juno Kael": false}
var ch = make(chan int, 6)
var MsgCh string
var isCh bool

//var isPro bool //是否活动

var wp = workerpool.New(1)

func init() {
	CronRun()
}

// 定时任务执行购物车
func CronRun() error {
	if config.Mode == 1 {
		return nil
	}
	log.Println("购物车cron加载...")
	for cs := range crs {
		crs[cs].Stop()
		time.Sleep(1 * time.Second)
		crs[cs] = cron.New()
		shoppingCron := GetShoppingCron("shopping_cron", cs)
		for _, v := range shoppingCron {
			//抢购物车
			err := func(v, cs string) error {
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("shopping_cron", v, cs)
					if IsRun[cs] {
						for i2 := range byCron {
							byCron[i2].Msg = "撞车了，抢购物车排队中..."
							UploadShopping(byCron[i2])
						}
						log.Println("撞车了，抢购物车排队中...")
					}
					wp.Submit(func() {
						IsRun[cs] = true
						hoppingRun(byCron)
						IsRun[cs] = false
					})
				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("inventoryCron", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		theShelvesCron := GetShoppingCron("the_shelves_cron", cs)
		for _, v := range theShelvesCron {
			err := func(v, cs string) error {
				//下架购物车
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("the_shelves_cron", v, cs)

					if IsRun[cs] {
						for i2 := range byCron {
							byCron[i2].Msg = "撞车了，下购物车排队中..."
							UploadShopping(byCron[i2])
						}
						log.Println("撞车了，下购物车排队中...")
					}
					wp.Submit(func() {
						IsRun[cs] = true
						theShelvesRun(byCron)
						IsRun[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("inventoryCron", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		inventoryCron := GetShoppingCron("inventory_cron", cs)
		for _, v := range inventoryCron {
			err := func(v, cs string) error {
				//加库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("inventory_cron", v, cs)
					if IsRun[cs] {
						for i2 := range byCron {
							byCron[i2].Msg = "撞车了，加库存排队中..."
							UploadShopping(byCron[i2])
						}
						log.Println("撞车了，加库存排队中...")
					}
					wp.Submit(func() {
						IsRun[cs] = true
						inventoryRun(byCron)
						IsRun[cs] = false
					})
				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("inventoryCron", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		xInventoryCron := GetShoppingCron("xinventory_cron", cs)
		for _, v := range xInventoryCron {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("xinventory_cron", v, cs)
					if IsRun[cs] {
						for i2 := range byCron {
							byCron[i2].Msg = "撞车了，下库存排队中..."
							UploadShopping(byCron[i2])
						}
						log.Println("撞车了，下库存排队中...")
					}
					wp.Submit(func() {
						IsRun[cs] = true
						xinventoryRun(byCron)
						IsRun[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("xInventoryCron", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		statusCron1 := GetShoppingCron("status_cron1", cs)
		for _, v := range statusCron1 {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("status_cron1", v, cs)
					if IsRun2[cs] {
						//for i2 := range byCron {
						//byCron[i2].Msg = "撞车了，状态排队中..."
						//UploadShopping(byCron[i2])
						//}
						log.Println("撞车了，状态排队中...")
					}
					wp.Submit(func() {
						IsRun2[cs] = true
						statusRun(byCron, 1)
						IsRun2[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("status_cron1", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		statusCron2 := GetShoppingCron("status_cron2", cs)
		for _, v := range statusCron2 {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("status_cron2", v, cs)
					if IsRun2[cs] {
						//for i2 := range byCron {
						//	byCron[i2].Msg = "撞车了，状态排队中..."
						//	UploadShopping(byCron[i2])
						//}
						log.Println("撞车了，状态排队中...")
					}
					wp.Submit(func() {
						IsRun2[cs] = true
						statusRun(byCron, 2)
						IsRun2[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("status_cron2", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		statusCron3 := GetShoppingCron("status_cron3", cs)
		for _, v := range statusCron3 {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("status_cron3", v, cs)
					if IsRun2[cs] {
						//for i2 := range byCron {
						//	byCron[i2].Msg = "撞车了，状态排队中..."
						//	UploadShopping(byCron[i2])
						//}
						log.Println("撞车了，状态排队中...")
					}
					wp.Submit(func() {
						IsRun2[cs] = true
						statusRun(byCron, 3)
						IsRun2[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("status_cron3", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		statusCron4 := GetShoppingCron("status_cron4", cs)
		for _, v := range statusCron4 {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("status_cron4", v, cs)
					if IsRun2[cs] {
						//for i2 := range byCron {
						//	byCron[i2].Msg = "撞车了，状态排队中..."
						//	UploadShopping(byCron[i2])
						//}
						log.Println("撞车了，状态排队中...")
					}
					wp.Submit(func() {
						IsRun2[cs] = true
						statusRun(byCron, 4)
						IsRun2[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("status_cron4", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}

		statusCron5 := GetShoppingCron("status_cron5", cs)
		for _, v := range statusCron5 {
			err := func(v, cs string) error {
				//下库存
				err := crs[cs].AddFunc(v, func() {
					byCron := GetShoppingByCron("status_cron5", v, cs)
					if IsRun2[cs] {
						//for i2 := range byCron {
						//	byCron[i2].Msg = "撞车了，状态排队中..."
						//	UploadShopping(byCron[i2])
						//}
						log.Println("撞车了，状态排队中...")
					}
					wp.Submit(func() {
						IsRun2[cs] = true
						statusRun(byCron, 5)
						IsRun2[cs] = false
					})

				})
				if err != nil {
					if !strings.Contains(err.Error(), "Expected 5 to 6 fields, found 0") && !strings.Contains(err.Error(), "Empty spec string") {
						log.Println("status_cron5", err)
						return err
					}
				}
				return nil
			}(v, cs)
			if err != nil {
				return err
			}
		}
		crs[cs].Start()

	}

	return nil
}

// 抢购物车
func hoppingRun(shs []mode.Shopping) {
	if len(shs) == 0 {
		return
	}
	log.Println("开始抢购物车", shs[0].ShoppingCron)
	//sPro = true
	//isPro = false
	for i := range shs {
		shs[i].IsActive = false
		shs[i].IsInventory = false
		shs[i].Msg = "抢购物车前,验证中..."
		UploadShopping(shs[i])
	}

	for i := range shs {
		//直接改价
		//shs[i].IsActive = true
		//shs[i].PromotionsStatus = "Active"
		ch <- 1
		wg[shs[0].Seller].Add(1)
		go crawlers(&shs[i], true)

	}
	wg[shs[0].Seller].Wait()

	//// 获取世界标准当前时间
	//now := time.Now().UTC()
	//// 加上5分钟
	//now1 := now.Add(5 * time.Minute)
	//now2 := now.Add(8760 * time.Hour)
	//date1 := now1.Format("2006/1/2 3:04:05")
	//date2 := now2.Format("2006/1/2 3:04:05")
	//第一次去除无需更新的产品
	var shos1 []mode.Shopping
	for i := range shs {
		if shs[i].IsUp {
			//shs[i].PromoStartDate = date1
			//shs[i].PromoEndDate = date2
			shs[i].Msg = "抢购物车,准备上传..."
			shos1 = append(shos1, shs[i])
			UploadShopping(shs[i])
		}
	}
	log.Printf("抓取完成，开始上传 %s 卖家购物车\n", shs[0].Seller)
	//生成上传用的文件
	//log.Println(shos1)
	if len(shos1) == 0 {
		log.Printf("%s 全部抢购成功\n", shs[0].Seller)
		return
	}
	shoppings(shos1, fmt.Sprintf("PriceUp%s.xlsm", se[shs[0].Seller]), true)
	//上传
	isCh = true
	for iio := 0; iio < 4; iio++ {
		chHoppingup(shos1[0])
		if isCh {
			log.Println(MsgCh)
			break
		}
		log.Println("开始重试：", iio+1)
	}
	log.Println("上传结束")
	if !isCh {
		log.Println(MsgCh)
		for i := range shos1 {
			shos1[i].Msg = MsgCh
			UploadShopping(shos1[i])
		}
		return
	}
	//log.Println("上传完成，等待15分钟后，查看是否抢购物车完成")
	//time.Sleep(time.Minute * 15)
	for i := range shos1 {
		shos1[i].Msg = "抢购物车完成,等待5分钟后验证..."
		UploadShopping(shos1[i])
	}

	log.Printf("%s 上传完成，等待5分钟后，查看是否抢购成功\n", shs[0].Seller)
	go func() {
		time.Sleep(time.Minute * 5)
		for i := range shos1 {
			shos1[i].Msg = "抢购物车后,验证中..."
			UploadShopping(shos1[i])
		}
		//log.Println("开始验证活动抢购情况")
		log.Printf("%s 开始验证抢购情况\n ", shs[0].Seller)
		for i := range shos1 {
			ch <- 1
			wg[shs[0].Seller].Add(1)
			go crawlers(&shos1[i], false)

		}
		wg[shs[0].Seller].Wait()
		var shos2 []mode.Shopping
		//更新信息
		for i := range shos1 {
			if shos1[i].IsUp {
				shos1[i].Msg = "抢购物车失败"
				shos2 = append(shos2, shos1[i])
			}
			UploadShopping(shos1[i])
		}
		if len(shos2) == 0 {
			log.Printf("%s 全部抢购物车成功\n ", shs[0].Seller)
		} else {
			log.Println(shs[0].Seller, "抢购物车结束", "有", len(shos2), "个商品没抢到购物车")
		}
	}()

	////对没有抢购成功的提取出来，先取消活动
	//var shos2 []mode.Shopping
	//for i := range shos1 {
	//	if shos1[i].IsUp {
	//		shos1[i].Msg = "活动抢购物车失败，开始无活动抢购物车"
	//		shos1[i].PromotionsStatus = "Delete All"
	//		shos2 = append(shos2, shos1[i])
	//	}
	//	UploadShopping(shos1[i])
	//}
	////生成上传用的文件
	//log.Println(shos2)
	//if len(shos2) == 0 {
	//	log.Println("全部抢购成功")
	//	return
	//}
	//log.Println("对活动抢购失败的商品，取消活动")
	//shoppings(shos2)
	////上传
	////cho()
	////if !isCh {
	////	log.Println(MsgCh)
	////	return
	////}
	////等待1分钟
	//log.Println("上传完成，等待1分钟，开始无活动抢购")
	//time.Sleep(time.Minute * 1)
	//log.Println("开始进行无活动抢购物车")
	////上传无活动的抢购
	//for i := range shos2 {
	//	shos2[i].IsActive = false
	//	shos2[i].Price = shos2[i].PromoPrice
	//}
	//isPro = false
	////生成上传用的文件
	//shoppings(shos2)
	//log.Println(shos2)
	////上传
	////cho()
	////if !isCh {
	////	log.Println(MsgCh)
	////	return
	////}
	////等待1分钟
	//log.Println("上传完成，等待10分钟，开始验证无活动抢购情况")
	//time.Sleep(time.Minute * 10)
	//log.Println("开始验证无活动抢购情况")
	//for i := range shos2 {
	//	ch <- 1
	//	wg.Add(1)
	//	crawlers(&shos2[i])
	//}
	//wg.Wait()
	//var shos3 []mode.Shopping
	////更新信息
	//for i := range shos2 {
	//	if shos2[i].IsUp {
	//		shos2[i].Msg = "抢购物车失败"
	//		shos3 = append(shos3, shos2[i])
	//	}
	//	UploadShopping(shos2[i])
	//}
	//log.Println("共有", len(shs), "个商品，其中有", len(shos3), "个商品没抢到购物车")
	//log.Println(shos3)
}

// 下架购物车
func theShelvesRun(shos []mode.Shopping) {
	if len(shos) == 0 {
		return
	}
	log.Println("开始下购物车", shos[0].TheShelvesCron)
	var shs []mode.Shopping
	for i := range shos {
		shos[i].IsActive = false
		shos[i].IsInventory = false
		shos[i].Price = shos[i].XPrice
		shos[i].IsTakeDown = true
		shos[i].Msg = "下购物车,准备上传..."
		shs = append(shs, shos[i])
		UploadShopping(shos[i])
	}

	sel := shs[0].Seller
	shoppings(shs, fmt.Sprintf("PriceUp%s.xlsm", se[sel]), false)
	//上传
	isCh = true
	for iio := 0; iio < 4; iio++ {
		chHoppingup(shos[0])
		if isCh {
			log.Println(MsgCh)
			break
		}
		log.Println("开始重试：", iio+1)
	}
	log.Println("上传结束")
	if !isCh {
		log.Println(MsgCh)
		for i := range shs {
			shs[i].Msg = MsgCh
			UploadShopping(shs[i])
		}
		return
	}
	for i := range shs {
		shs[i].Msg = "下购物车完成,等待5分钟后验证..."
		UploadShopping(shs[i])
	}
	log.Printf("%s 上传完成，等待5分钟后，查看是否下架成功\n", shos[0].Seller)
	go func() {
		//检测是否下架成功
		time.Sleep(time.Minute * 5)
		log.Printf("%s 开始检测是否下架成功\n", shos[0].Seller)
		for i := range shs {
			shs[i].Msg = "下架购物车后,验证中..."
			UploadShopping(shs[i])
		}
		for i2 := range shs {
			ch <- 1
			wg2[shos[0].Seller].Add(1)
			go crawlersSell(&shs[i2])
		}
		wg2[shos[0].Seller].Wait()
		var shs1 []mode.Shopping
		for i2 := range shs {
			if shs[i2].Msg != "下购物车成功" && shs[i2].Msg != "商品不存在" {
				shs1 = append(shs1, shs[i2])
			}
		}
		shs = shs1
		if len(shs) > 0 {
			log.Printf("%s 下架购物车结束,有%d个商品下架失败,\n", sel, len(shs))
		} else {
			log.Printf("%s 全部下架购物车成功\n", sel)
		}
	}()

}

// 加库存
func inventoryRun(shs []mode.Shopping) {
	if len(shs) == 0 {
		return
	}
	log.Println("开始加库存", shs[0].InventoryCron)
	var shs1 []mode.Shopping
	for i := range shs {
		shs[i].IsInventory = true
		shs[i].Msg = "加库存,准备上传..."
		shs1 = append(shs1, shs[i])
		UploadShopping(shs[i])
	}
	Inventorys(shs1, fmt.Sprintf("InventoryUp%s.xlsx", se[shs[0].Seller]), true)
	//上传
	isCh = true
	for iio := 0; iio < 4; iio++ {
		chHoppingup(shs[0])
		if isCh {
			log.Println(MsgCh)
			break
		}
		log.Println("开始重试：", iio+1)
	}

	log.Println("上传结束")
	if !isCh {
		log.Println(MsgCh)
		for i := range shs1 {
			shs1[i].Msg = MsgCh
			UploadShopping(shs1[i])
		}
		return
	}
	for i := range shs1 {
		shs1[i].Msg = "加库存完成"
		UploadShopping(shs1[i])
	}
	log.Printf("%s 库存添加完成\n", shs[0].Seller)

}

// 下库存
func xinventoryRun(shs []mode.Shopping) {
	if len(shs) == 0 {
		return
	}
	log.Println("开始下库存", shs[0].XInventoryCron)
	var shs1 []mode.Shopping
	for i := range shs {
		shs[i].IsInventory = true
		shs[i].Msg = "下库存,准备上传..."
		shs1 = append(shs1, shs[i])
		UploadShopping(shs[i])
	}
	Inventorys(shs1, fmt.Sprintf("InventoryUp%s.xlsx", se[shs[0].Seller]), false)
	//上传
	isCh = true
	for iio := 0; iio < 4; iio++ {
		chHoppingup(shs[0])
		if isCh {
			log.Println(MsgCh)
			break
		}
		log.Println("开始重试：", iio+1)
	}
	log.Println("上传结束")
	if !isCh {
		log.Println(MsgCh)
		for i := range shs1 {
			shs1[i].Msg = MsgCh
			UploadShopping(shs1[i])
		}
		return
	}
	for i := range shs1 {
		shs1[i].Msg = "下库存完成"
		UploadShopping(shs1[i])
	}
	log.Printf("%s 下库存完成\n", shs[0].Seller)

}

func statusRun(shs []mode.Shopping, in int) {
	if len(shs) == 0 {
		return
	}

	var cron string
	switch in {
	case 1:
		cron = shs[0].StatusCron1
	case 2:
		cron = shs[0].StatusCron2
	case 3:
		cron = shs[0].StatusCron3
	case 4:
		cron = shs[0].StatusCron4
	case 5:
		cron = shs[0].StatusCron5
	}
	log.Println("开始状态检测", cron)
	UpStatus(shs, in)
	log.Println("状态检测结束", cron)
}

// 根据当前价格，当前商家，做出处理
func processor(shopping *mode.Shopping, price, seller string) *mode.Shopping {
	//无需更新
	FloorPricef, err := strconv.ParseFloat(strings.TrimSpace(shopping.FloorPrice), 64)
	if err != nil {
		log.Println(err)
	}
	pricef, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println(err)
	}
	if shopping.Seller == seller {
		shopping.Msg = "已抢到购物车"
		shopping.IsUp = false
		return shopping
	} else if FloorPricef > pricef-0.01 {
		shopping.Msg = "已达到最低价格,放弃购物车"
		shopping.Price = price
		shopping.IsUp = false
		return shopping
	} else {
		shopping.Msg = "抢购物车,验证完成,等待其他完成..."
		shopping.Price = price
		shopping.IsUp = true
	}
	if shopping.IsActive {
		shopping.Price = strconv.FormatFloat(pricef+20, 'f', 2, 64)
		shopping.PromoPrice = strconv.FormatFloat(pricef-0.01, 'f', 2, 64)
	} else {
		shopping.Price = strconv.FormatFloat(pricef-0.01, 'f', 2, 64)
	}
	return shopping
}

// 获取商家和价格，并处理数据,需要更新的，IsUp会变成true
func crawlers(shopping *mode.Shopping, isUp bool) *mode.Shopping {
	defer func() {
		<-ch
		wg[shopping.Seller].Done()
		if isUp {
			if shopping.Msg == "抢购物车前,验证中..." {
				shopping.Msg = "莫名的错误"
			}
			UploadShopping(*shopping)
		}
	}()
	for i := 0; i < 20; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		request, _ := http.NewRequest("GET", "https://www.walmart.com/ip/"+strconv.Itoa(shopping.PrId), nil)

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
				log.Println("代理IP无效，开始重试")
				shopping.Msg = "验证代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				shopping.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				shopping.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				shopping.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				shopping.Msg = err.Error()
			}
			continue
		}
		result := ""
		defer response.Body.Close()
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("解析body错误，重新开始")
				shopping.Msg = "解析错误"
				continue
			}
			defer reader.Close()
			con, err := io.ReadAll(reader)
			if err != nil {
				log.Println("gzip解压错误，重新开始")
				shopping.Msg = "gzip解压错误"
				continue
			}
			result = string(con)
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
					log.Println("代理IP无效,重新开始")
					shopping.Msg = "代理IP无效"
				} else {
					log.Println("错误信息：" + err.Error())
					log.Println("出现错误，如果同id连续出现请联系我，重新开始")
					shopping.Msg = "解析错误"
				}
				continue
			}
			result = string(dataBytes)
		}

		if err != nil {
			if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
				log.Println("代理IP无效，开始重试")
				shopping.Msg = "验证代理IP无效"
			} else {
				log.Println("错误信息：" + err.Error())
				shopping.Msg = "验证请求失败"
			}
			continue
		}

		if strings.Contains(result, "This page could not be found.") {
			log.Println("商品不存在")
			shopping.Msg = "商品不存在"
			shopping.IsUp = false
			return shopping
		}

		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控,更换IP继续")
			shopping.Msg = "验证被风控"
			config.IsC = !config.IsC
			continue
		}

		//图片
		img := regexp.MustCompile("<meta property=\"og:image\" content=\"(.*?)\"/>").FindAllStringSubmatch(result, -1)
		if len(img) > 0 {
			shopping.Img = img[0][1]
		}
		//价格
		var price string
		//prices := regexp.MustCompile("<span itemprop=\"price\".*?.{0,20}\\$([.\\d]+).{0,20}?</span>").FindAllStringSubmatch(result, -1)
		//if len(prices) > 0 {
		//	price = prices[0][1]
		//}
		// 使用具体的XPath查找目标span标签
		price1 := regexp.MustCompile(`"best[^{]+?,"priceDisplay":"([^"]+)"`)
		price2 := price1.FindAllString(result, -1)
		if len(price2) > 0 {
			//log.Println(result)
			// Check if the matched string contains "priceDisplay"
			if strings.Contains(price2[0], `"priceDisplay":"`) {
				// Split the string to isolate the part after "priceDisplay":"
				parts := strings.Split(price2[0], `"priceDisplay":"`)
				if len(parts) > 1 {
					// Further split to get just the value before the closing quote
					valueParts := strings.Split(parts[1], `"`)
					if len(valueParts) > 0 {
						//fmt.Println("Extracted Value:", valueParts[0]) // Should print "Now $16.99"
						reg := regexp.MustCompile(`[^\d.]`)
						numericValue := reg.ReplaceAllString(valueParts[0], "")
						fmt.Println("Numeric Value:", numericValue) // Should print "16.99"
						price = numericValue
					} else {
						fmt.Println("No value extracted after priceDisplay")
					}
				} else {
					fmt.Println("No priceDisplay part found in string")
				}
			} else {
				fmt.Println("String does not contain priceDisplay")
			}
		} else {
			fmt.Println("No matches found or result is empty")
			price = "" // 如果是空的，赋值空字符串
		}
		var selle string
		//卖家与配送
		selles := regexp.MustCompile("\"sellerDisplayName\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
		if len(selles) > 0 {
			selle = selles[0][1]
		}
		if price != "" && selle != "" {
			shopping = processor(shopping, price, selle)
		} else {
			shopping.IsUp = false
			shopping.Msg = "未发现价格或卖家"
			i = 19
			continue
		}
		return shopping
	}
	shopping.IsUp = false
	return shopping
}

// 获取商家
func crawlersSell(shopping *mode.Shopping) *mode.Shopping {
	defer func() {
		<-ch
		wg2[shopping.Seller].Done()
		UploadShopping(*shopping)
	}()
	for i := 0; i < 20; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		request, _ := http.NewRequest("GET", "https://www.walmart.com/ip/"+strconv.Itoa(shopping.PrId), nil)

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
				log.Println("代理IP无效，开始重试")
				shopping.Msg = "代理IP无效"
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				shopping.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				shopping.Msg = "代理超频"
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				shopping.Msg = "代理超频"
			} else {
				log.Println("错误信息：" + err.Error())
				shopping.Msg = "错误：" + err.Error()
			}
			continue
		}
		result := ""
		defer response.Body.Close()
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body) // gzip解压缩
			if err != nil {
				log.Println("解析body错误，重新开始")
				continue
			}
			defer reader.Close()
			con, err := io.ReadAll(reader)
			if err != nil {
				log.Println("gzip解压错误，重新开始")
				continue
			}
			result = string(con)
		} else {
			dataBytes, err := io.ReadAll(response.Body)
			if err != nil {
				if strings.Contains(err.Error(), "Proxy Bad Serve") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Service Unavailable") {
					log.Println("代理IP无效,重新开始")

				} else {
					log.Println("错误信息：" + err.Error())
					log.Println("出现错误，如果同id连续出现请联系我，重新开始")

				}

				continue
			}
			result = string(dataBytes)
		}

		if strings.Contains(result, "This page could not be found.") {
			log.Println("商品不存在")
			shopping.Msg = "商品不存在"
			return shopping
		}

		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控,更换IP继续")
			shopping.Msg = "被风控"
			config.IsC = !config.IsC
			continue
		}
		var selle string
		//卖家与配送
		selles := regexp.MustCompile("\"sellerDisplayName\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
		if len(selles) > 0 {
			selle = selles[0][1]
		}
		if selle != "" {
			if selle != shopping.Seller {
				shopping.Msg = "下架购物车成功"
			} else {
				shopping.Msg = "下架购物车失败"
			}
		} else {
			shopping.Msg = "未发现卖家"
		}
		//价格
		prices := regexp.MustCompile("<span itemprop=\"price\".*?.{0,20}\\$([.\\d]+).{0,20}?</span>").FindAllStringSubmatch(result, -1)
		if len(prices) > 0 {
			shopping.Price = prices[0][1]
		}
		return shopping
	}
	shopping.Msg = "下架购物车失败"
	return shopping
}

// 获取商家
func crawlersUpSell(shopping *mode.Shopping, in int) *mode.Shopping {
	defer func() {
		<-ch
		wg3[shopping.Seller].Done()
		switch in {
		case 1:
			if shopping.Status1 == "正在获取" {
				shopping.Status1 = "莫名错误"
			}
		case 2:
			if shopping.Status2 == "正在获取" {
				shopping.Status2 = "莫名错误"
			}
		case 3:
			if shopping.Status3 == "正在获取" {
				shopping.Status3 = "莫名错误"
			}
		case 4:
			if shopping.Status4 == "正在获取" {
				shopping.Status4 = "莫名错误"
			}
		case 5:
			if shopping.Status5 == "正在获取" {
				shopping.Status5 = "莫名错误"
			}
		}
		UploadShoppingStatus(*shopping)
	}()
	for i := 0; i < 20; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		proxy_str := fmt.Sprintf("http://%s:%s@%s", config.ProxyUser, config.ProxyPass, config.ProxyUrl)
		proxy, _ := url.Parse(proxy_str)
		var client *http.Client
		if config.IsC {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
		} else {
			client = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		}
		request, _ := http.NewRequest("GET", "https://www.walmart.com/ip/"+strconv.Itoa(shopping.PrId), nil)

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
				log.Println("代理IP无效，开始重试")
				switch in {
				case 1:
					shopping.Status1 = "获取失败"
				case 2:
					shopping.Status2 = "获取失败"
				case 3:
					shopping.Status3 = "获取失败"
				case 4:
					shopping.Status4 = "获取失败"
				case 5:
					shopping.Status5 = "获取失败"
				}
			} else if strings.Contains(err.Error(), "441") {
				log.Println("代理超频！暂停10秒后继续...")
				time.Sleep(time.Second * 10)
				switch in {
				case 1:
					shopping.Status1 = "获取失败"
				case 2:
					shopping.Status2 = "获取失败"
				case 3:
					shopping.Status3 = "获取失败"
				case 4:
					shopping.Status4 = "获取失败"
				case 5:
					shopping.Status5 = "获取失败"
				}
			} else if strings.Contains(err.Error(), "440") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				switch in {
				case 1:
					shopping.Status1 = "获取失败"
				case 2:
					shopping.Status2 = "获取失败"
				case 3:
					shopping.Status3 = "获取失败"
				case 4:
					shopping.Status4 = "获取失败"
				case 5:
					shopping.Status5 = "获取失败"
				}
			} else if strings.Contains(err.Error(), "Request Rate Over Limit") {
				log.Println("代理宽带超频！暂停5秒后继续...")
				time.Sleep(time.Second * 5)
				switch in {
				case 1:
					shopping.Status1 = "获取失败"
				case 2:
					shopping.Status2 = "获取失败"
				case 3:
					shopping.Status3 = "获取失败"
				case 4:
					shopping.Status4 = "获取失败"
				case 5:
					shopping.Status5 = "获取失败"
				}
			} else {
				log.Println("错误信息：" + err.Error())
				switch in {
				case 1:
					shopping.Status1 = "获取失败"
				case 2:
					shopping.Status2 = "获取失败"
				case 3:
					shopping.Status3 = "获取失败"
				case 4:
					shopping.Status4 = "获取失败"
				case 5:
					shopping.Status5 = "获取失败"
				}
			}
			continue
		}
		result := ""
		reader, err := gzip.NewReader(response.Body) // gzip解压缩
		if err != nil {
			log.Println("解析body错误，重新开始")
			switch in {
			case 1:
				shopping.Status1 = "获取失败"
			case 2:
				shopping.Status2 = "获取失败"
			case 3:
				shopping.Status3 = "获取失败"
			case 4:
				shopping.Status4 = "获取失败"
			case 5:
				shopping.Status5 = "获取失败"
			}
			continue
		}
		defer reader.Close()
		con, err := io.ReadAll(reader)
		if err != nil {
			log.Println("gzip解压错误，重新开始")
			switch in {
			case 1:
				shopping.Status1 = "获取失败"
			case 2:
				shopping.Status2 = "获取失败"
			case 3:
				shopping.Status3 = "获取失败"
			case 4:
				shopping.Status4 = "获取失败"
			case 5:
				shopping.Status5 = "获取失败"
			}
			continue
		}
		defer response.Body.Close()
		result = string(con)

		if strings.Contains(result, "This page could not be found.") {
			log.Println("商品不存在")
			switch in {
			case 1:
				shopping.Status1 = "获取失败"
			case 2:
				shopping.Status2 = "获取失败"
			case 3:
				shopping.Status3 = "获取失败"
			case 4:
				shopping.Status4 = "获取失败"
			case 5:
				shopping.Status5 = "获取失败"
			}
			return shopping
		}

		fk := regexp.MustCompile("(Robot or human?)").FindAllStringSubmatch(result, -1)
		if len(fk) > 0 {
			log.Println("被风控,更换IP继续")
			config.IsC = !config.IsC
			switch in {
			case 1:
				shopping.Status1 = "获取失败"
			case 2:
				shopping.Status2 = "获取失败"
			case 3:
				shopping.Status3 = "获取失败"
			case 4:
				shopping.Status4 = "获取失败"
			case 5:
				shopping.Status5 = "获取失败"
			}
			continue
		}
		var selle string
		//卖家与配送
		selles := regexp.MustCompile("\"sellerDisplayName\":\"(.*?)\"").FindAllStringSubmatch(result, -1)
		if len(selles) > 0 {
			selle = selles[0][1]
		}
		if selle != "" {
			if selle != shopping.Seller {
				switch in {
				case 1:
					shopping.Status1 = "无购物车"
				case 2:
					shopping.Status2 = "无购物车"
				case 3:
					shopping.Status3 = "无购物车"
				case 4:
					shopping.Status4 = "无购物车"
				case 5:
					shopping.Status5 = "无购物车"
				}
			} else {
				switch in {
				case 1:
					shopping.Status1 = "有购物车"
				case 2:
					shopping.Status2 = "有购物车"
				case 3:
					shopping.Status3 = "有购物车"
				case 4:
					shopping.Status4 = "有购物车"
				case 5:
					shopping.Status5 = "有购物车"
				}
			}
		} else {
			switch in {
			case 1:
				shopping.Status1 = "无购物车"
			case 2:
				shopping.Status2 = "无购物车"
			case 3:
				shopping.Status3 = "无购物车"
			case 4:
				shopping.Status4 = "无购物车"
			case 5:
				shopping.Status5 = "无购物车"
			}

		}
		switch in {
		case 1:
			shopping.J99991 = "无"
		case 2:
			shopping.J99992 = "无"
		case 3:
			shopping.J99993 = "无"
		case 4:
			shopping.J99994 = "无"
		case 5:
			shopping.J99995 = "无"
		}

		//最低跟卖价格
		startingFrom := regexp.MustCompile(`"priceType":.{0,20},"priceString":"\$([^<]+?)",`).FindAllStringSubmatch(result, -1)
		if len(startingFrom) > 0 {
			replace := strings.Replace(startingFrom[0][1], ",", "", -1)
			replace = strings.Replace(replace, ".", "", -1)
			if strings.Contains(replace, "9999") {
				switch in {
				case 1:
					shopping.J99991 = startingFrom[0][1]
				case 2:
					shopping.J99992 = startingFrom[0][1]
				case 3:
					shopping.J99993 = startingFrom[0][1]
				case 4:
					shopping.J99994 = startingFrom[0][1]
				case 5:
					shopping.J99995 = startingFrom[0][1]
				}

			}
		}

		//图片
		img := regexp.MustCompile("<meta property=\"og:image\" content=\"(.*?)\"/>").FindAllStringSubmatch(result, -1)
		if len(img) > 0 {
			shopping.Img = img[0][1]
		}
		//价格
		prices := regexp.MustCompile("<span itemprop=\"price\".*?.{0,20}\\$([.\\d]+).{0,20}?</span>").FindAllStringSubmatch(result, -1)
		if len(prices) > 0 {
			shopping.Price = prices[0][1]
		}
		return shopping
	}
	switch in {
	case 1:
		shopping.Status1 = "获取失败"
	case 2:
		shopping.Status2 = "获取失败"
	case 3:
		shopping.Status3 = "获取失败"
	case 4:
		shopping.Status4 = "获取失败"
	case 5:
		shopping.Status5 = "获取失败"
	}
	return shopping
}

// 根据数据生成对应改价表格，用于上传后台
func shoppings(shoppingss []mode.Shopping, name string, isUp bool) {
	tools.SafeDeleteFile(`C:\Users\Administrator\Desktop\` + name)
	shoppingsx, err := excelize.OpenFile(`C:\Users\Administrator\Desktop\Price.xlsm`)
	if err != nil {
		return
	}
	num := 5
	for i := range shoppingss {
		var pr string
		if isUp {
			pr = shoppingss[i].Price
		} else {
			pr = shoppingss[i].XPrice
		}
		if shoppingss[i].IsActive {
			if err := shoppingsx.SetSheetRow("Offer", "C"+strconv.Itoa(num), &[]interface{}{strings.TrimSpace(shoppingss[i].Sku), strings.TrimSpace(pr), shoppingss[i].PromotionsStatus, strings.TrimSpace(shoppingss[i].PromoPrice), "Reduced", nil, shoppingss[i].PromoStartDate, shoppingss[i].PromoEndDate}); err != nil {
				log.Println(err)
			}
		} else {
			if err := shoppingsx.SetSheetRow("Offer", "C"+strconv.Itoa(num), &[]interface{}{strings.TrimSpace(shoppingss[i].Sku), strings.TrimSpace(pr)}); err != nil {
				log.Println(err)
			}
		}
		num++
	}

	fileName := `C:\Users\Administrator\Desktop\` + name

	shoppingsx.SaveAs(fileName)

}

// 根据数据生成对应库存表格，用于上传后台
func Inventorys(shoppingss []mode.Shopping, name string, isAdd bool) {
	tools.SafeDeleteFile(`C:\Users\Administrator\Desktop\` + name)
	shoppingsx, err := excelize.OpenFile(`C:\Users\Administrator\Desktop\Inventory.xlsx`)
	if err != nil {
		return
	}
	num := 3

	for i := range shoppingss {
		var inv string
		//是否添加
		if isAdd {
			inv = shoppingss[i].Inventory
		} else {
			inv = "0"
		}
		if err := shoppingsx.SetSheetRow("ExampleSheet", "A"+strconv.Itoa(num), &[]interface{}{shoppingss[i].Sku, inv, shoppingss[i].CenterId}); err != nil {
			log.Println(err)
		}
		num++
	}

	fileName := `C:\Users\Administrator\Desktop\` + name

	shoppingsx.SaveAs(fileName)

}

// 操作浏览器
func chHoppingup(sho mode.Shopping) {
	Mutexs[sho.Seller].Lock()
	defer Mutexs[sho.Seller].Unlock()

	// 创建一个分配器来连接到已经打开的浏览器
	allocator, _ := chromedp.NewRemoteAllocator(context.Background(), fmt.Sprintf("ws://127.0.0.1:%s/devtools/browser", se[sho.Seller]))

	// 创建一个新的浏览器上下文
	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()
	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	chromedp.Run(ctx,
		startup(sho),
	)

}

// 控制器
func startup(sho mode.Shopping) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if loginup(ctx, sho) {
			if sho.IsTakeDown {
				//循环下架5次
				for i := 0; i < 2; i++ {
					MsgCh, isCh = uploadup(ctx, sho)
					log.Printf("%s 等待30秒，再次下架，当前第%d次\n", sho.Seller, i+1)
					chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(ctx)
					time.Sleep(time.Second * 30)
				}
			} else {
				MsgCh, isCh = uploadup(ctx, sho)
			}
		} else {
			MsgCh, isCh = "登录失败", false
		}
		return err
	}
}

// 检测是否登录，未登录就登录
func loginup(ctx context.Context, sho mode.Shopping) bool {
	timeout, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout)
	for i := 0; i < 2; i++ {
		timeout0, cancel0 := context.WithTimeout(ctx, 30*time.Second)
		defer cancel0()

		err := chromedp.WaitVisible(`input[id="radioFulfillmentSF"]`).Do(timeout0)
		if err != nil {
			timeout02, cancel02 := context.WithTimeout(ctx, 20*time.Second)
			defer cancel02()
			err := chromedp.WaitVisible(`input[data-automation-id="uname"]`).Do(timeout02)
			if err != nil {
				log.Println("页面加载失败，重新开始加载")
				timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
				defer cancel01()
				chromedp.Stop().Do(timeout01)
				chromedp.Sleep(time.Second * 1).Do(timeout)
				chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout01)
			} else {
				log.Println("未登录状态")
				break
			}
		} else {
			log.Println("已是登录状态")
			return true

		}
	}
	log.Println("开始登录")
	chromedp.Sleep(time.Second * 2).Do(timeout)
	chromedp.SendKeys(`input[data-automation-id="uname"]`, users[sho.Seller][0]).Do(timeout)
	chromedp.Sleep(time.Second * 2).Do(timeout)
	chromedp.SendKeys(`input[data-automation-id="pwd"]`, users[sho.Seller][1]+kb.Enter).Do(timeout)
	chromedp.Sleep(time.Second * 2).Do(timeout)
	timeout01, cancel01 := context.WithTimeout(ctx, 2*time.Second)
	defer cancel01()
	chromedp.SendKeys(`input[data-automation-id="pwd"]`, kb.Enter).Do(timeout01)
	chromedp.Sleep(time.Second * 30).Do(timeout)
	chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout)
	chromedp.Sleep(time.Second * 10).Do(timeout)
	timeout03, cancel03 := context.WithTimeout(ctx, 30*time.Second)
	defer cancel03()
	err := chromedp.WaitVisible(`input[id="radioFulfillmentSF"]`).Do(timeout03)
	if err != nil {
		log.Println("登录失败")
		chromedp.Stop().Do(ctx)
		return false
	}
	log.Println("登录成功")
	return true
}

// 修改价格
func uploadup(ctx context.Context, sho mode.Shopping) (string, bool) {
	var msg string
	var iso bool
	for i := 0; i < 2; i++ {
		if i != 0 {
			log.Println("开始重试：", i)
			timeout00, cance00 := context.WithTimeout(ctx, 10*time.Second)
			defer cance00()
			chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout00)

		}
		var i string
		timeout, cancel := context.WithTimeout(ctx, 6*time.Minute)
		defer cancel()
		for i := 0; i < 2; i++ {
			timeout0, cancel0 := context.WithTimeout(ctx, 30*time.Second)
			defer cancel0()
			timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel01()
			err := chromedp.WaitVisible(`input[id="radioFulfillmentSF"]`).Do(timeout0)
			if err != nil {
				log.Println("重新加载页面")
				chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout01)
			} else {
				break
			}
		}
		//橫幅
		chromedp.EvaluateAsDevTools(`document.getElementsByClassName("h-48")[0].style.display = "none"`, &i).Do(timeout)
		//右下角弹窗
		chromedp.EvaluateAsDevTools(`
			var elements = document.querySelectorAll('[data-vertical-alignment="Bottom Right Aligned"]');
			for (var i = 0; i < elements.length; i++) {
				elements[i].style.display = 'none';
		}`, &i).Do(timeout)
		chromedp.Click(`input[id="radioFulfillmentSF"]`).Do(timeout)
		chromedp.Sleep(time.Second * 2).Do(timeout)
		chromedp.Click(`input[id="radioMPTemplate"]`).Do(timeout)
		chromedp.Sleep(time.Second * 2).Do(timeout)
		log.Println("开始上传")

		if !sho.IsInventory {
			//改价
			chromedp.EvaluateAsDevTools(`document.getElementsByTagName("input")[5].style= ""`, &i).Do(timeout)
			chromedp.Sleep(time.Second * 1).Do(timeout)
			chromedp.SendKeys(`input[type="file"]`, fmt.Sprintf(`C:\Users\Administrator\Desktop\PriceUp%s.xlsm`, se[sho.Seller])).Do(timeout)
			chromedp.Sleep(time.Second * 2).Do(timeout)
			ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if !isExe(ctx1, `.css-1hwfws3`) {
				log.Println("没有找到下拉框，启用第二方案(价格)")
				log.Println("点击第一个下拉框")
				chromedp.MouseClickXY(420, 820).Do(timeout)
				chromedp.Sleep(time.Second * 3).Do(timeout)
				log.Println("选中第一个下拉框")
				chromedp.MouseClickXY(420, 820).Do(timeout)
				chromedp.Sleep(time.Second * 3).Do(timeout)
				log.Println("开始获取标签内的值")
				var source string
				// 设置超时时间为5秒
				ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				err := chromedp.Run(ctx1,
					chromedp.Text("body > div.js-content > div:nth-child(1) > section > div._2qDzh > div > div:nth-child(6) > div._2e6uJ.Q7vB2 > div > div:nth-child(3) > div.section-wrapper > div:nth-child(3) > div > div > div.css-1hwfws3 > div.css-1hxmfou-singleValue", &source),
				)
				if err != nil {
					log.Printf("获取标签内的值失败: %v", err)
					msg = "获取标签内的值失败"
					iso = false
					return msg, iso
				}
				log.Printf("标签内的值: %s", source)
				if source != "Bulk Price and Promotion Update" {
					log.Println("标签值不匹配，终止任务")
					msg = "标签值不匹配/获取失败，终止任务"
					iso = false
					return msg, iso
				} else {
					log.Println("点击提交")
					chromedp.MouseClickXY(660, 730).Do(timeout)
					chromedp.Sleep(time.Second * 5).Do(timeout)
				}
			} else {
				log.Println("点击第一个下拉框")
				err := chromedp.Run(ctx,
					chromedp.Click(`.css-1hwfws3`),
					chromedp.Sleep(time.Second*2),
					chromedp.Click(`#react-select-2-option-1`),
					chromedp.Sleep(time.Second*2),
				)
				if err != nil {
					log.Printf("选择第一个下拉框失败: %v", err)
					msg = "选择第一个下拉框失败，终止任务"
					iso = false
					return msg, iso
				}
				log.Println("开始获取标签内的值")
				var source string
				ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				err = chromedp.Run(ctx1,
					chromedp.Click(`body > div.js-content > div:nth-child(1) > section > div._2qDzh > div > div:nth-child(6) > div._2e6uJ.Q7vB2 > div > div:nth-child(3) > div.section-wrapper > div:nth-child(3) > div > div > div.css-1hwfws3 > div.css-1hxmfou-singleValue`),
					chromedp.Text("body > div.js-content > div:nth-child(1) > section > div._2qDzh > div > div:nth-child(6) > div._2e6uJ.Q7vB2 > div > div:nth-child(3) > div.section-wrapper > div:nth-child(3) > div > div > div.css-1hwfws3 > div.css-1hxmfou-singleValue", &source),
				)
				//if err != nil {
				//	log.Printf("获取标签内的值失败: %v", err)
				//	msg = "标签值不匹配/获取失败，终止任务"
				//	iso = false
				//	continue
				//}
				if err != nil {
					log.Printf("获取标签内的值失败: %v", err)
					msg = "获取标签内的值失败"
					iso = false
					return msg, iso
				}

				log.Printf("标签内的值: %s", source)
				if source != "Bulk Price and Promotion Update" {
					log.Println("标签值不匹配，终止任务")
					msg = "标签值不匹配/获取失败，终止任务"
					iso = false
					return msg, iso
				} else {
					log.Println("点击提交")
					err = chromedp.Run(ctx,
						chromedp.Click(`body > div.js-content > div:nth-child(1) > section > div._2qDzh > div > div:nth-child(6) > div._2e6uJ.Q7vB2 > div > div:nth-child(3) > div.section-wrapper > div:nth-child(3) > button`),
						chromedp.Sleep(time.Second*5),
					)
				}

				if err != nil {
					log.Printf("点击提交按钮失败: %v", err)
					msg = "点击提交按钮失败"
					iso = false
					return msg, iso
				}
			}
		} else {
			//改库存
			chromedp.EvaluateAsDevTools(`document.getElementsByTagName("input")[5].style= ""`, &i).Do(timeout)
			chromedp.Sleep(time.Second * 1).Do(timeout)
			chromedp.SendKeys(`input[type="file"]`, fmt.Sprintf(`C:\Users\Administrator\Desktop\InventoryUp%s.xlsx`, se[sho.Seller])).Do(timeout)
			chromedp.Sleep(time.Second * 2).Do(timeout)
			if !isExe(ctx, `.css-1hwfws3`) {
				log.Println("没有找到下拉框，启用第二方案4444444444")
				//return "没有找到第一个下拉框", false
				log.Println("点击第一个下拉框")
				chromedp.MouseClickXY(420, 828).Do(timeout)
				chromedp.Sleep(time.Second * 5).Do(timeout)
				log.Println("选中第一个下拉框")
				chromedp.MouseClickXY(420, 751).Do(timeout)
				chromedp.Sleep(time.Second * 5).Do(timeout)
			} else {
				log.Println("点击第一个下拉框")
				chromedp.Click(`.css-1hwfws3`).Do(timeout)
				chromedp.Sleep(time.Second * 2).Do(timeout)
				log.Println("选中第一个下拉框")
				chromedp.Click(`#react-select-2-option-0`).Do(timeout)
				chromedp.Sleep(time.Second * 2).Do(timeout)
			}

			var source string
			if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools("document.documentElement.outerHTML", &source)); err != nil {
				log.Println("获取源码失败")
				msg = "获取源码失败"
				iso = false
				continue
			}
			if !strings.Contains(source, "Bulk Inventory Update") {
				log.Println("选择出现错误")
				msg = "选择出现错误"
				iso = false
				continue
			}
		}

		log.Println("提交")

		if !isExe(ctx, `._3_kwe`) {
			log.Println("没有找到提交")
			msg = "没有找到提交"
			iso = false
			continue
		}

		err := chromedp.Click(`._3_kwe`).Do(timeout)
		if err != nil {
			log.Println("提交失败")
			msg = "提交失败"
			iso = false
			continue
		}
		log.Println("提交成功")
		time.Sleep(time.Second * 10)
		chromedp.Sleep(time.Second * 60).Do(ctx)
		return "成功", true
	}
	return msg, iso

}

// 模拟操作更新后台价格等
//func cho() {
//	// 配置
//	//var example string
//	options := append(chromedp.DefaultExecAllocatorOptions[:],
//		chromedp.NoDefaultBrowserCheck, //不检查默认浏览器
//		chromedp.Flag("headless", false),
//		chromedp.Flag("blink-settings", "imagesEnabled=true"), //开启图像界面,重点是开启这个
//		chromedp.Flag("ignore-certificate-errors", true),      //忽略错误
//		chromedp.Flag("disable-web-security", true),           //禁用网络安全标志
//		chromedp.Flag("disable-extensions", true),             //开启插件支持
//		chromedp.Flag("disable-default-apps", true),
//		// chromedp.Flag("disable-gpu", true), //开启gpu渲染
//		//chromedp.WindowSize(1920, 1080), // 设置浏览器分辨率（窗口大小）
//		chromedp.Flag("hshopping.Ide-scrollbars", true),
//		chromedp.Flag("mute-audio", true),
//		chromedp.Flag("no-sandbox", true),
//		chromedp.Flag("no-default-browser-check", true),
//		//chromedp.NoFirstRun,                                                                                                                       //设置网站不是首次运行
//		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"), //设置UserAgent
//	)
//	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
//	defer cancel()
//	// 初始化chromedp上下文，后续这个页面都使用这个上下文进行操作
//	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
//	defer cancel()
//	// 设置超时时间
//	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
//	defer cancel()
//	log.Println("开始登录")
//	err := chromedp.Run(ctx,
//		// 设置webdriver检测反爬
//		chromedp.ActionFunc(func(cxt context.Context) error {
//			_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => false, });").Do(cxt)
//			return err
//		}),
//		//loadCookies(),
//		start(),
//
//		// 停止网页加载
//		chromedp.Stop(),
//	)
//	if err != nil {
//		log.Println(err)
//	}
//}

//// 控制器
//func start() chromedp.ActionFunc {
//	return func(ctx context.Context) (err error) {
//		login(ctx)
//		MsgCh, isCh = upload(ctx)
//		return err
//	}
//}
//
//// 登录
//func login(ctx context.Context) {
//	timeout, cancel := context.WithTimeout(ctx, 180*time.Second)
//	defer cancel()
//	chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout)
//	for i := 0; i < 2; i++ {
//		timeout0, cancel0 := context.WithTimeout(ctx, 60*time.Second)
//		defer cancel0()
//		timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
//		defer cancel01()
//		err := chromedp.WaitVisible(`input[data-automation-id="uname"]`).Do(timeout0)
//		if err != nil {
//			log.Println("重新加载登录頁面")
//			chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout01)
//		} else {
//			break
//		}
//	}
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	chromedp.SendKeys(`input[data-automation-id="uname"]`, username).Do(timeout)
//	chromedp.Sleep(time.Second * 1).Do(timeout)
//	chromedp.SendKeys(`input[data-automation-id="pwd"]`, passwrod+kb.Enter).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("登录结束")
//	return
//
//}
//
//// 上传
//func upload(ctx context.Context) (string, bool) {
//	var i string
//	var err error
//	timeout, cancel := context.WithTimeout(ctx, 180*time.Second)
//	defer cancel()
//	//chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout)
//	for i := 0; i < 2; i++ {
//		timeout0, cancel0 := context.WithTimeout(ctx, 60*time.Second)
//		defer cancel0()
//		timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
//		defer cancel01()
//		err := chromedp.WaitVisible(`input[id="radioFulfillmentSF"]`).Do(timeout0)
//		if err != nil {
//			log.Println("重新加载页面")
//			chromedp.Navigate("https://seller.walmart.com/items-and-inventory/bulk-updates?returnUrl=%2Fitems-and-inventory%2Fmanage-items").Do(timeout01)
//		} else {
//			break
//		}
//	}
//
//	log.Println("登录检测通过")
//	log.Println("开始改价")
//	chromedp.Click(`input[id="radioFulfillmentSF"]`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	chromedp.Click(`input[id="radioMPTemplate"]`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("开始上传")
//	chromedp.EvaluateAsDevTools(`document.getElementsByTagName("input")[5].style= ""`, &i).Do(timeout)
//	chromedp.Sleep(time.Second * 1).Do(timeout)
//	chromedp.SendKeys(`input[type="file"]`, `C:\Users\Administrator\Desktop\out.xlsm`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("点击第一个下拉框")
//	if !isExe(ctx, `.css-1hwfws3`) {
//		return "没有找到第一个下拉框", false
//	}
//	chromedp.Click(`.css-1hwfws3`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("选中第一个下拉框")
//	chromedp.Click(`#react-select-2-option-1`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("点击第二个下拉框")
//	if !isExe(ctx, `.css-xqe4hi-control .css-1hwfws3`) {
//		return "没有找到第二个下拉框", false
//	}
//	chromedp.Click(`.css-xqe4hi-control .css-1hwfws3`).Do(timeout)
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	log.Println("选中第二个下拉框")
//	var str string
//	if isPro {
//		str = "Promotional Price Update"
//		chromedp.Click(`#react-select-3-option-1`).Do(timeout)
//	} else {
//		str = "Price Updates"
//		chromedp.Click(`#react-select-3-option-0`).Do(timeout)
//	}
//	chromedp.Sleep(time.Second * 2).Do(timeout)
//	var source string
//	if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools("document.documentElement.outerHTML", &source)); err != nil {
//		log.Println("获取源码失败")
//		return "获取源码失败", false
//	}
//	if !strings.Contains(source, "Bulk Pricing and Promotion Update") && !(strings.Contains(source, str)) {
//		log.Println("选择出现错误")
//		return "选择出现错误", false
//	}
//	log.Println("提交")
//	if !isExe(ctx, `._3_kwe`) {
//		return "没有找到提交", false
//	}
//	err = chromedp.Click(`._3_kwe`).Do(timeout)
//	if err != nil {
//		return "提交失败", false
//	}
//	chromedp.Sleep(time.Second * 60).Do(timeout)
//	return "成功", true
//
//}

//func UpIds(ids []string) {
//	byIds := GetShoppingByIds(ids)
//	for i := range byIds {
//		byIds[i].Msg = "获取购物车状态中..."
//		UploadShopping(byIds[i])
//	}
//	for i2 := range byIds {
//		ch <- 1
//		wg3[byIds[0].Seller].Add(1)
//		go crawlersUpSell(&byIds[i2], in)
//	}
//	wg3[byIds[0].Seller].Wait()
//}

func UpStatus(byIds []mode.Shopping, in int) {
	for i := range byIds {
		switch in {
		case 1:
			byIds[i].Status1 = "正在获取"
		case 2:
			byIds[i].Status2 = "正在获取"
		case 3:
			byIds[i].Status3 = "正在获取"
		case 4:
			byIds[i].Status4 = "正在获取"
		case 5:
			byIds[i].Status5 = "正在获取"

		}
		UploadShoppingStatus(byIds[i])
	}
	for i2 := range byIds {
		ch <- 1
		wg3[byIds[0].Seller].Add(1)
		go crawlersUpSell(&byIds[i2], in)
	}
	wg3[byIds[0].Seller].Wait()
}

func isExe(ctx context.Context, ex string) bool {
	timeout2, cancel2 := context.WithTimeout(ctx, 2*time.Second)
	defer cancel2()
	err := chromedp.WaitVisible(ex).Do(timeout2)
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}
