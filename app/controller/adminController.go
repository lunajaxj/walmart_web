package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"walmart_web/app/activity"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/product"
	"walmart_web/app/productSales"
	"walmart_web/app/shopping"
	"walmart_web/app/storeInformation"
	"walmart_web/app/timedUploads"
	"walmart_web/app/tools"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func Home(c *gin.Context) {
	c.Redirect(http.StatusFound, "/admin/index")
}

func LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Table(c *gin.Context) {
	//marks := product.GETProductMark()
	acts := activity.GetActivity()
	var actName []string
	for i := range acts {
		actName = append(actName, acts[i].Name)
	}
	c.HTML(http.StatusOK, "product-table.html", gin.H{
		"keysNum":  len(product.TaskKeys),
		"idsNum":   len(product.TaskIds),
		"urlsNum":  len(product.TaskUrls),
		"actsNum":  len(product.TaskActs),
		"keysNum2": product.TaskKeysNum,
		"idsNum2":  product.TaskIdsNum,
		"urlsNum2": product.TaskUrlNum,
		"actsNum2": product.TaskActsNum,
		//"marks":    marks,
		"isRun":    product.IsRun,
		"actNames": tools.UniqueArr(actName),
	})
}

func Shopping(c *gin.Context) {
	msg := shopping.GetShoppingMsg()
	sellers := shopping.GetShoppingSeller()
	c.HTML(http.StatusOK, "product-shopping.html", gin.H{
		"msgCh":   shopping.MsgCh,
		"msgs":    msg,
		"sellers": sellers,
	})
}

func TimedUploads(c *gin.Context) {
	msg := timedUploads.GetMsg()               // 获取消息
	sellers := timedUploads.GetSeller()        // 获取卖家信息 (包含 PID 和 shop_name)
	file := timedUploads.GetFile()             // 获取文件列表
	genres := []string{"价格", "库存", "类别", "促销"} // 任务类型

	// 获取 PID 和店铺名称的映射，确保返回的是 map[string]string
	shopMapping := timedUploads.GetShopMapping()

	// 将数据传递给前端的 HTML 模板
	c.HTML(http.StatusOK, "timed-uploads.html", gin.H{
		"msgs":        msg,
		"sellers":     sellers, // 包含 PID 和 shop_name
		"files":       file,
		"genres":      genres,
		"shopMapping": shopMapping, // 将 PID 和店铺映射信息传递给前端
	})
}

func StockAvailability(c *gin.Context) {
	c.HTML(http.StatusOK, "stock-availability.html", nil)
}
func StoreInformation(c *gin.Context) {

	c.HTML(http.StatusOK, "store-information.html", gin.H{"num": fmt.Sprintf("更新任务 %d个", storeInformation.Num), "info": fmt.Sprintf("抓取任务 %d个", storeInformation.InfoNum)})
}

func EditStoreInformation(c *gin.Context) {
	sid := c.Query("sid")
	information, _ := storeInformation.GetStoreInformation(1, 1, "", sid, "", "", "")
	c.HTML(http.StatusOK, "edit-store-information.html", gin.H{"ShopID": information[0].ShopID, "PID": information[0].PID, "SID": information[0].SID, "AffiliatedCompany": information[0].AffiliatedCompany, "Note1": information[0].Note1, "Note2": information[0].Note2, "Note3": information[0].Note3, "Note4": information[0].Note4, "Note5": information[0].Note5})
}
func AddStoreInformation(c *gin.Context) {
	c.HTML(http.StatusOK, "add-store-information.html", nil)
}

func Mail(c *gin.Context) {
	c.HTML(http.StatusOK, "mail.html", nil)
}

func Edit(c *gin.Context) {
	id := c.Query("id")
	remark := c.Query("remark")
	mark := c.Query("mark")
	marks := product.GETProductMark()

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"id":     id,
		"remark": remark,
		"mark":   mark,
		"marks":  marks,
	})
}

func Chart(c *gin.Context) {
	query := c.DefaultQuery("brands", "")
	query2 := c.DefaultQuery("sellers", "")
	var brands []string

	var chart map[string]mode.ProductBrandss
	var sellers string
	if query != "" {
		sellers = "店铺："
		split := strings.Split(query, "\r\n")
		s := strings.Split(split[0], ",")
		var bfb1 int
		var bfb2 = 33
		var err error
		switch len(s) {
		case 1:
			bfb1, err = strconv.Atoi(s[0])
			if err != nil {
				c.HTML(http.StatusOK, "chart.html", gin.H{
					"brands":  brands,
					"chart":   chart,
					"sellers": "第一行必须是品牌占比,不能有小数点",
				})
				return
			}
		case 2:
			bfb1, err = strconv.Atoi(s[0])
			bfb2, err = strconv.Atoi(s[1])
			if err != nil {
				c.HTML(http.StatusOK, "chart.html", gin.H{
					"brands":  brands,
					"chart":   chart,
					"sellers": "第一行必须是品牌占比,不能有小数点",
				})
				return
			}

		}
		brands, chart = product.GetChart(split[1:], 1, bfb1, bfb2)
		//brands, chart = product.GetChart(split, 1, 100)
		sort.Strings(brands)
		for s := range chart {
			if sellers != "店铺：" {
				sellers += "  |  " + s
			} else {
				sellers += s
			}
		}
	}

	if query2 != "" {
		sellers = "店铺："
		split := strings.Split(query2, "\r\n")
		brands, chart = product.GetChartSellers(split)
		for s := range chart {
			if sellers != "店铺：" {
				sellers += "  |  " + s
			} else {
				sellers += s
			}
		}
	}
	//fmt.Println(chart)
	c.HTML(http.StatusOK, "chart.html", gin.H{
		"brands":  brands,
		"chart":   chart,
		"sellers": sellers,
	})
}

func CrawlerEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "crawler-edit.html", nil)
}

func EditShopping(c *gin.Context) {
	id := c.Query("id")
	seller := c.Query("seller")
	floorPrice := c.Query("floorPrice")
	xPrice := c.Query("xPrice")
	inventory := c.Query("inventory")
	shoppingCron := c.Query("shoppingCron")
	theShelvesCron := c.Query("theShelvesCron")
	inventoryCron := c.Query("inventoryCron")
	xinventoryCron := c.Query("xinventoryCron")
	statusCron1 := c.Query("statusCron1")
	statusCron2 := c.Query("statusCron2")
	statusCron3 := c.Query("statusCron3")
	statusCron4 := c.Query("statusCron4")
	statusCron5 := c.Query("statusCron5")

	sales := c.Query("sales")
	note := c.Query("note")
	name := c.Query("name")
	c.HTML(http.StatusOK, "edit-shopping.html", gin.H{
		"id":             id,
		"seller":         seller,
		"floorPrice":     floorPrice,
		"xPrice":         xPrice,
		"inventory":      inventory,
		"note":           note,
		"sales":          sales,
		"name":           name,
		"shoppingCron":   shoppingCron,
		"theShelvesCron": theShelvesCron,
		"inventoryCron":  inventoryCron,
		"xinventoryCron": xinventoryCron,
		"statusCron1":    statusCron1,
		"statusCron2":    statusCron2,
		"statusCron3":    statusCron3,
		"statusCron4":    statusCron4,
		"statusCron5":    statusCron5,
	})
}
func LoadShopMappings() []mode.ShopMapping {
	var shopMappings []mode.ShopMapping
	err := config.Db.Select(&shopMappings, "SELECT pid, shop_name FROM shop_mapping")
	if err != nil {
		log.Println("获取店铺映射失败:", err)
		return nil
	}
	log.Println("店铺映射关系LoadShopMappings：", shopMappings)
	return shopMappings
}
func EditTimedUploads(c *gin.Context) {
	id := c.Query("id")
	genre := c.Query("genre")
	msg := c.Query("msg")
	file := c.Query("file")
	cron := c.Query("cron")
	name := c.Query("name")
	genres := []string{"价格", "库存", "类别", "促销"}
	files := timedUploads.GetFile()
	seller := c.Query("seller")        // 从前端传递来的 PID
	shopMappings := LoadShopMappings() // 获取shop_mapping表的数据
	// 定义固定的店铺数组
	c.HTML(http.StatusOK, "edit-timed-uploads.html", gin.H{
		"id":      id,
		"genre":   genre,
		"genres":  genres,
		"msg":     msg,
		"file":    file,
		"files":   files,
		"cron":    cron,
		"name":    name,
		"seller":  seller,
		"sellers": shopMappings,
	})
}

func CreateTimedUploads(c *gin.Context) {
	genres := []string{"价格", "库存", "类别", "促销"}
	files := timedUploads.GetFile()
	// 定义固定的店铺数组
	shopMappings := LoadShopMappings()
	c.HTML(http.StatusOK, "create-timed-uploads.html", gin.H{
		"genres":  genres,
		"files":   files,
		"sellers": shopMappings, // 将 sellers 数组传递给模板
	})
}

func File(c *gin.Context) {
	c.HTML(http.StatusOK, "file.html", nil)
}

func ProductSales(c *gin.Context) {
	c.HTML(http.StatusOK, "product-sales.html", gin.H{"num": fmt.Sprintf("更新任务 %d", productSales.Num)})
}

func UpdateShopMapping(c *gin.Context) {
	var request struct {
		Pid      string `json:"pid"`
		ShopName string `json:"shopName"`
	}

	// 解析前端传递的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数错误", "code": 1})
		return
	}

	// 更新数据库中的店铺名称
	_, err := config.Db.Exec("UPDATE shop_mapping SET shop_name = ? WHERE pid = ?", request.ShopName, request.Pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新店铺名称失败", "code": 1})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "店铺名称修改成功", "code": 0})
}
