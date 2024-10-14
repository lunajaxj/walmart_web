package product

import (
	"log"
	"strconv"
	"time"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

var res []mode.ProductDetails

var TaskIds []string
var TaskKeys []string
var TaskUrls []string
var TaskActs []string

var TaskKeyIs = true
var TaskIdsIs = true
var TaskUrlsIs = true
var TaskActsIs = true

var IsRun = true

func init() {
	//定时判断是否需要获取id
	go func() {
		// 每 60 秒钟时执行一次
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			config.MuxId.Lock()
			if TaskIdsIs && len(TaskIds) > 0 && IsRun {
				TaskIdsIs = false
				config.MuxId.Unlock()
				log.Println("开始抓取id，当前数量：", len(TaskIds))
				go func() {
					defer func() {
						TaskIdsIs = true
					}()
					for {
						if !IsRun {
							break
						}
						id := GetTaskId()
						if id == "" {
							break
						}
						config.Ch <- 1
						config.IdWg.Add(1)
						go CrawlerId(id)
					}
					config.IdWg.Wait()
					log.Println("id任务全部完成")
					TaskIdsNum = 0
				}()
			} else {
				config.MuxId.Unlock()
			}
		}
	}()

	//定时判断是否需要获取key
	go func() {
		// 每 60 秒钟时执行一次
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			config.MuxKey.Lock()
			if TaskKeyIs && len(TaskKeys) > 0 && IsRun {
				TaskKeyIs = false
				config.MuxKey.Unlock()
				log.Println("开始抓取key，当前数量：", len(TaskKeys))
				go func() {
					defer func() {
						TaskKeyIs = true
					}()
					for {
						if !IsRun {
							break
						}
						key := GetTaskKey()
						if key == "" {
							break
						}
						config.Ch <- 1
						config.KeyWg.Add(1)
						go CrawlerKey(key)
					}
					config.KeyWg.Wait()
					log.Println("关键词任务全部完成")
				}()
			} else {
				config.MuxKey.Unlock()
			}
		}
	}()

	//定时判断是否需要获取url
	go func() {
		// 每 60 秒钟时执行一次
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			config.MuxUrl.Lock()
			if TaskUrlsIs && len(TaskUrls) > 0 && IsRun {
				TaskUrlsIs = false
				config.MuxUrl.Unlock()
				log.Println("开始抓取url，当前数量：", len(TaskUrls))
				go func() {
					defer func() {
						TaskUrlsIs = true
					}()
					for {
						if !IsRun {
							break
						}
						urll := GetTaskUrl()
						if urll == "" {
							break
						}
						config.Ch <- 1
						config.UrlWg.Add(1)
						go CrawlerUrl(urll)
					}
					config.UrlWg.Wait()
					log.Println("url任务全部完成")
				}()
			} else {
				config.MuxUrl.Unlock()
			}
		}
	}()

	//定时判断是否需要获取活动
	go func() {
		// 每 60 秒钟时执行一次
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			config.MuxAct.Lock()
			if TaskActsIs && len(TaskActs) > 0 && IsRun {
				TaskActsIs = false
				config.MuxAct.Unlock()
				log.Println("开始抓取活动，当前数量：", len(TaskActs))
				go func() {
					defer func() {
						TaskActsIs = true
					}()
					for {
						if !IsRun {
							break
						}
						act := GetTaskAct()
						if act == "" {
							break
						}
						config.Ch <- 1
						config.ActWg.Add(1)
						go CrawlerActivity(act)
					}
					config.ActWg.Wait()
					log.Println("活动任务全部完成")
				}()
			} else {
				config.MuxAct.Unlock()
			}
		}
	}()

}

func GetTaskId() (id string) {
	config.MuxId.Lock()
	defer config.MuxId.Unlock()
	if len(TaskIds) > 0 {
		TaskIds, id = tools.Remove(TaskIds, 0)
	} else {
		return ""
	}
	return id

}
func GetTaskKey() (key string) {
	config.MuxKey.Lock()
	defer config.MuxKey.Unlock()
	if len(TaskKeys) > 0 {
		TaskKeys, key = tools.Remove(TaskKeys, 0)
	} else {
		return ""
	}
	return key
}
func GetTaskUrl() (urll string) {
	config.MuxUrl.Lock()
	defer config.MuxUrl.Unlock()
	if len(TaskUrls) > 0 {
		TaskUrls, urll = tools.Remove(TaskUrls, 0)
	} else {
		return ""
	}
	return urll
}

func GetTaskAct() (acy string) {
	config.MuxAct.Lock()
	defer config.MuxAct.Unlock()
	if len(TaskActs) > 0 {
		TaskActs, acy = tools.Remove(TaskActs, 0)
	} else {
		return ""
	}
	return acy
}

func AddCrawlerId(ids []string) {
	config.MuxId.Lock()
	defer config.MuxId.Unlock()
	idsz := GetProductId(ids)
	var idsd []string
	for i := range idsz {
		idsd = append(idsd, strconv.Itoa(idsz[i].Id))
	}
	ids = tools.UniqueArrT(ids, idsd)
	TaskIds = tools.UniqueArr(tools.MergeArray(TaskIds, ids))

}

func AddCrawlerIdNo(ids []string) {
	config.MuxId.Lock()
	defer config.MuxId.Unlock()
	TaskIds = tools.UniqueArr(tools.MergeArray(TaskIds, ids))
}

func AddCrawlerKey(keys []string) {
	config.MuxKey.Lock()
	defer config.MuxKey.Unlock()
	TaskKeys = tools.UniqueArr(tools.MergeArray(TaskKeys, keys))
}

func AddCrawlerUrl(urls []string) {
	config.MuxUrl.Lock()
	defer config.MuxUrl.Unlock()
	TaskUrls = tools.UniqueArr(tools.MergeArray(TaskUrls, urls))
}

func AddCrawlerAct(acts []string) {
	config.MuxAct.Lock()
	defer config.MuxAct.Unlock()
	TaskActs = tools.UniqueArr(tools.MergeArray(TaskActs, acts))
}
