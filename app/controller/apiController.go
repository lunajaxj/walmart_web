package controller

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"walmart_web/app/activity"
	"walmart_web/app/category"
	"walmart_web/app/config"
	"walmart_web/app/keyword"
	"walmart_web/app/mail"
	"walmart_web/app/productSales"
	"walmart_web/app/stockAvailability"
	"walmart_web/app/storeInformation"

	"walmart_web/app/mode"
	"walmart_web/app/product"
	"walmart_web/app/shopping"
	"walmart_web/app/timedUploads"
	"walmart_web/app/tools"
	"walmart_web/app/user"
	"walmart_web/app/walLog"
)

var sasMut = new(sync.Mutex)

func SetCh(c *gin.Context) {
	num, err := strconv.Atoi(c.PostForm("num"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "设置失败,请输入纯数字",
		})
		return
	}
	if len(config.Ch) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "设置失败,请在没有任务执行时修改",
		})
		return
	}
	if num > 20 {
		num = 20
	} else if num < 1 {
		num = 1
	}
	config.Ch = make(chan int, num)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "设置成功",
	})
}

func GetProduct(c *gin.Context) {
	var actid []string
	var keyid []string
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryTree := c.DefaultQuery("categoryName", "")
	id := c.DefaultQuery("id", "")
	title := c.DefaultQuery("title", "")
	sellers := c.DefaultQuery("sellers", "")
	sellersType := c.DefaultQuery("sellersType", "")
	tags := c.DefaultQuery("tags", "")
	brands := c.DefaultQuery("brands", "")

	actDate := c.DefaultQuery("actDate", "")
	remark := c.DefaultQuery("remark", "")
	mark := c.DefaultQuery("mark", "")
	keyName := c.DefaultQuery("keyName", "")
	actName := c.DefaultQuery("actName", "")

	rating := c.DefaultQuery("rating", "")
	comments := c.DefaultQuery("comments", "")
	price := c.DefaultQuery("price", "")

	if actName != "" || actDate != "" {
		acts := activity.GetActivityWhere(actName, actDate)
		var actids string
		for i := range acts {
			if i == 0 {
				actids += acts[i].Ids
				continue
			}
			actids += "," + acts[i].Ids
		}
		actid = strings.Split(actids, ",")
		actid = tools.UniqueArr(actid)
		if len(actid) == 0 {
			actid = []string{"0000000"}
		}
	}
	if keyName != "" {
		split := strings.Split(keyName, ",")
		keys := keyword.GetKeywordName(split)
		var keyids string
		for i := range keys {
			if i == 0 {
				keyids += keys[i].Ids
				continue
			}
			keyids += "," + keys[i].Ids
		}
		keyid = strings.Split(keyids, ",")
		keyid = tools.UniqueArr(keyid)
		if len(keyid) == 0 {
			keyid = []string{"0000000"}
		}
	}
	pors, count := product.GetProduct(page, limit, id, title, sellers, sellersType, tags, brands, remark, mark, categoryTree, rating, comments, price, actid, keyid)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  pors,
		"msg":   "成功",
	})
}

func UpdateIsRun(c *gin.Context) {
	product.IsRun = !product.IsRun

	var mess string
	if product.IsRun {
		mess = "任务开启成功"
	} else {
		mess = "任务停止成功，当前线程结束后将停止任务"
	}

	tools.SafeDeleteFile("ids.txt")
	fi, err := os.OpenFile("ids.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewWriter(fi) // 创建 Reader
	for i := range product.TaskIds {
		r.WriteString(product.TaskIds[i] + "\n")
	}
	r.Flush()

	tools.SafeDeleteFile("keyss.txt")
	fi, err = os.OpenFile("keyss.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r = bufio.NewWriter(fi) // 创建 Reader
	for i := range product.TaskKeys {
		r.WriteString(product.TaskKeys[i] + "\n")
	}
	r.Flush()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  mess,
	})
}
func GetCategory(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	//log.Println(string(data))
	cats := category.GetCategory(string(data))
	c.JSON(http.StatusOK, gin.H{
		"status": tools.Result{Code: 200, Message: "成功"},
		"data":   cats,
	})
}

func GetShopping(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	id := c.DefaultQuery("id", "")
	note := c.DefaultQuery("note", "")
	name := c.DefaultQuery("name", "")
	msg := c.DefaultQuery("msg", "")
	sales := c.DefaultQuery("sales", "")
	seller := c.DefaultQuery("seller", "")
	status1 := c.DefaultQuery("status1", "")
	status2 := c.DefaultQuery("status2", "")
	status3 := c.DefaultQuery("status3", "")
	status4 := c.DefaultQuery("status4", "")
	status5 := c.DefaultQuery("status5", "")
	j99991 := c.DefaultQuery("j99991", "")
	j99992 := c.DefaultQuery("j99992", "")
	j99993 := c.DefaultQuery("j99993", "")
	j99994 := c.DefaultQuery("j99994", "")
	j99995 := c.DefaultQuery("j99995", "")
	shos, count := shopping.GetShopping(page, limit, id, note, name, msg, seller, sales, status1, status2, status3, status4, status5, j99991, j99992, j99993, j99994, j99995)
	for i := range shos {
		shos[i].Status = fmt.Sprintf("%s|%s|%s|%s|%s", shos[i].Status1, shos[i].Status2, shos[i].Status3, shos[i].Status4, shos[i].Status5)
		shos[i].StatusCron = fmt.Sprintf("%s|%s|%s|%s|%s", shos[i].StatusCron1, shos[i].StatusCron2, shos[i].StatusCron3, shos[i].StatusCron4, shos[i].StatusCron5)
		shos[i].J9999 = fmt.Sprintf("%s|%s|%s|%s|%s", shos[i].J99991, shos[i].J99992, shos[i].J99993, shos[i].J99994, shos[i].J99995)

	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  shos,
		"msg":   "成功",
	})
}

func MailRun(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	getMail := mail.GetMail(id)
	if id == "" || getMail.MId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "失败",
		})
		return
	}
	if !mail.IsRun[getMail.Seller] {
		go mail.Run(getMail)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "正在执行中,请勿重复执行",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "开始执行",
		})
	}

}
func MailDelMsg(c *gin.Context) {
	mail.DelMsg()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
	})
}

func GetMails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	shos, count := mail.GetMails(page, limit)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  shos,
		"msg":   "成功",
	})
}

func GetTimedUploads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	genre := c.DefaultQuery("genre", "")
	name := c.DefaultQuery("name", "")
	msg := c.DefaultQuery("msg", "")
	file := c.DefaultQuery("file", "")
	seller := c.DefaultQuery("seller", "")
	// 打印 API 调用的输入参数
	//log.Printf("获取定时上传: page=%d, limit=%d, genre=%s, name=%s, msg=%s, file=%s, seller=%s", page, limit, genre, name, msg, file, seller)

	// 调用 GetTimedUploadsWithShopName 来获取包含 ShopName 的数据
	shos, count := timedUploads.GetTimedUploadsWithShopName(page, limit, genre, name, msg, file, seller)
	// 返回结果给前端时记录输出
	//log.Printf("返回结果: %+v, 总数: %d", shos, count)
	// 返回数据，data 现在包含了 ShopName 字段
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  shos,
		"msg":   "成功",
	})
}

func GetStockAvailability(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	itemId := c.DefaultQuery("itemId", "")

	sas, count := stockAvailability.GetStockAvailability(page, limit, itemId)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  sas,
		"msg":   "成功",
	})
}
func SARemove(c *gin.Context) {
	ids := c.PostForm("ids")
	count := stockAvailability.Remove(ids)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}

func GetStoreInformation(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("AccountShopName", "")
	affiliatedCompany := c.DefaultQuery("AffiliatedCompany", "")
	hasTargetLink := c.DefaultQuery("HasTargetLink", "")
	onSaleProductCount := c.DefaultQuery("OnSaleProductCount", "")
	sas, count := storeInformation.GetStoreInformation(page, limit, name, "", affiliatedCompany, hasTargetLink, onSaleProductCount)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  sas,
		"msg":   "成功",
	})
}

func SIRemove(c *gin.Context) {
	names := c.PostForm("ids")
	count := storeInformation.Remove(names)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "删除成功",
	})
}

func CrawlerId(c *gin.Context) {
	ids := c.PostForm("ids")
	split := strings.Split(ids, "\n")
	split = tools.UniqueArr(split)
	product.AddCrawlerId(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}
func CrawlerIdUp(c *gin.Context) {
	ids := c.PostForm("ids")
	split := strings.Split(ids, "\n")
	split = tools.UniqueArr(split)
	product.AddCrawlerIdNo(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}

func CrawlerKey(c *gin.Context) {
	keys := c.PostForm("keys")
	split := strings.Split(keys, "\n")
	split = tools.UniqueArr(split)
	split = tools.UniqueArrT(split, config.KeysFile)
	config.OutFileKeys(split)
	product.AddCrawlerKey(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}
func CrawlerUrl(c *gin.Context) {
	urls := c.PostForm("urls")
	split := strings.Split(urls, "\n")
	split = tools.UniqueArr(split)
	split = tools.UniqueArrT(split, config.UrlsFile)
	config.OutFileUrls(split)
	product.AddCrawlerUrl(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}

func CrawlerAct(c *gin.Context) {
	acts := c.PostForm("acts")
	split := strings.Split(acts, "\n")
	split = tools.UniqueArr(split)
	product.AddCrawlerAct(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}

func CrawlerCategory(c *gin.Context) {
	rating := c.DefaultPostForm("rating", "")
	comments := c.DefaultPostForm("comments", "")
	price := c.DefaultPostForm("price", "")
	categoryTree := c.DefaultPostForm("categoryName", "")
	sellersType := c.DefaultPostForm("sellersType", "")
	if categoryTree == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprint("关键词类目不能为空"),
		})
		return
	}

	pors, count := product.GetProductById(sellersType, categoryTree, rating, comments, price)
	var ids []string
	for i := range pors {
		ids = append(ids, pors[i])
	}
	log.Println(len(ids))
	product.AddCrawlerIdNo(ids)
	log.Println(len(product.TaskIds))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  fmt.Sprint("提交成功,共", count, "个id任务"),
	})
}

func GetLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": walLog.GetLogs(),
		"msg":  "提交成功",
	})
}
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	// 删除 Session 对应的 Cookie
	cookie := &http.Cookie{
		Name:   "session-name",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(c.Writer, cookie)
	c.Redirect(http.StatusFound, "/login/index")
}

func Login(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	if len(username) == 0 && len(password) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": "",
			"msg":  "账号或密码不能为空",
		})
	}
	u := user.GetUser(username)
	if len(u) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": "",
			"msg":  "账号或密码错误",
		})
		return
	}
	if u[0].Password != password {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": "",
			"msg":  "账号或密码错误",
		})
		return
	}
	session := sessions.Default(c)
	session.Set("username", username)
	session.Set("rule", u[0].Rule)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": "",
		"msg":  "登录成功",
	})
}

func DelFile(c *gin.Context) {
	config.MuxKey.Lock()
	defer config.MuxKey.Unlock()
	tools.SafeDeleteFile("keys.txt")
	tools.SafeDeleteFile("urls.txt")
	config.KeysFile = []string{}
	config.UrlsFile = []string{}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "提交成功",
	})
}

func EditProduct(c *gin.Context) {
	id := c.DefaultPostForm("id", "")
	remark := c.DefaultPostForm("remark", "")
	mark := c.DefaultPostForm("mark", "")
	xmark := c.DefaultPostForm("xmark", "")
	if xmark != "" {
		mark = xmark
	}
	count := product.EditProductMark(id, remark, mark)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}

func EditSho(c *gin.Context) {
	id := c.DefaultPostForm("id", "")
	seller := c.DefaultPostForm("seller", "")
	floorPrice := c.DefaultPostForm("floorPrice", "")
	xPrice := c.DefaultPostForm("xPrice", "")
	inventory := c.DefaultPostForm("inventory", "")
	note := c.DefaultPostForm("note", "")
	name := c.DefaultPostForm("name", "")
	sales := c.DefaultPostForm("sales", "")
	shoppingCron := c.DefaultPostForm("shoppingCron", "")
	theShelvesCron := c.DefaultPostForm("theShelvesCron", "")
	inventoryCron := c.DefaultPostForm("inventoryCron", "")
	xinventoryCron := c.DefaultPostForm("xinventoryCron", "")
	statusCron1 := c.DefaultPostForm("statusCron1", "")
	statusCron2 := c.DefaultPostForm("statusCron2", "")
	statusCron3 := c.DefaultPostForm("statusCron3", "")
	statusCron4 := c.DefaultPostForm("statusCron4", "")
	statusCron5 := c.DefaultPostForm("statusCron5", "")
	atoi, _ := strconv.Atoi(id)
	sale, err := strconv.Atoi(strings.TrimSpace(sales))
	if err != nil {
		log.Println(err)
	}
	count := shopping.UploadSho(mode.Shopping{PrId: atoi, Seller: seller, FloorPrice: floorPrice, XPrice: xPrice, Inventory: inventory, Note: note, Name: name, Sales: sale, ShoppingCron: shoppingCron, TheShelvesCron: theShelvesCron, InventoryCron: inventoryCron, XInventoryCron: xinventoryCron, StatusCron1: statusCron1, StatusCron2: statusCron2, StatusCron3: statusCron3, StatusCron4: statusCron4, StatusCron5: statusCron5})
	shopping.CronRun()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}

func EditTu(c *gin.Context) {
	id := c.DefaultPostForm("id", "")
	seller := c.DefaultPostForm("seller", "")
	genre := c.DefaultPostForm("genre", "")
	cron := c.DefaultPostForm("cron", "")
	file := c.DefaultPostForm("file", "")
	msg := c.DefaultPostForm("msg", "")
	name := c.DefaultPostForm("name", "")
	atoi, _ := strconv.Atoi(id)
	count := timedUploads.UploadTimedUploads(mode.TimedUploads{
		TuId:   atoi,
		Name:   name,
		Genre:  genre,
		Seller: seller,
		Cron:   cron,
		File:   file,
		Msg:    msg,
	})
	timedUploads.CronRun()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}
func CreateTu(c *gin.Context) {
	seller := c.DefaultPostForm("seller", "")
	genre := c.DefaultPostForm("genre", "")
	cron := c.DefaultPostForm("cron", "")
	file := c.DefaultPostForm("file", "")
	name := c.DefaultPostForm("name", "")
	msg := c.DefaultPostForm("msg", "")
	timedUploads.AddTimedUploads(mode.TimedUploads{
		Name:   name,
		Genre:  genre,
		Seller: seller,
		Cron:   cron,
		Msg:    msg,
		File:   file,
	})
	timedUploads.CronRun()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})
}

func ShoppingRemove(c *gin.Context) {
	ids := c.PostForm("ids")
	count := shopping.Remove(ids)
	shopping.CronRun()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}

func TimedUploadsRemove(c *gin.Context) {
	ids := c.PostForm("ids")
	count := timedUploads.Remove(ids)
	timedUploads.CronRun()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "提交成功",
	})
}

func GetFile(c *gin.Context) {
	files := timedUploads.GetFile2()
	var fs []map[string]string
	for i := range files {
		fs = append(fs, map[string]string{"File": files[i]["name"], "UpdateTime": files[i]["time"]})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": len(files),
		"data":  fs,
		"msg":   "成功",
	})
}

func RemoveFile(c *gin.Context) {
	files := c.PostForm("files")
	split := strings.Split(files, "，")
	for i := range split {
		tools.SafeDeleteFile("file/" + split[i])
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nil,
		"msg":  "删除成功",
	})
}

func UploadGtin(c *gin.Context) {
	go func() {
		gtin := c.PostForm("gtin")
		split := strings.Split(gtin, "\n")
		split = tools.UniqueArr(split)
		for i := range split {
			product.ChGtin <- 1
			sp := strings.Split(strings.TrimSpace(split[i]), ",")
			go product.UploadProductGtin(sp[0], sp[1])
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})

}

func UpIds(c *gin.Context) {
	ids := c.PostForm("ids")
	split := strings.Split(ids, "\n")
	split = tools.UniqueArr(split)
	//go shopping.UpIds(split)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "提交成功",
	})
}

func Download(c *gin.Context) {
	tools.SafeDeleteFile("./file.xlsx")
	var pros []mode.ProductDetails
	if len(config.DowIds) != 0 {
		pros = product.GetProductId(tools.UniqueArr(config.DowIds))
	}
	xlsx := excelize.NewFile()
	num := 2
	if err := xlsx.SetSheetRow("Sheet1", "A1", &[]interface{}{"标记", "备注", "id", "图片", "商品码类型", "商品码值", "品牌", "标签", "标题", "评分", "评论数量", "价格", "卖家", "配送", "变体1", "变体2", "变体id", "到达时间", "类目1", "类目2", "类目3", "类目4", "类目5", "类目6", "类目7", "创建时间", "更新时间"}); err != nil {
		log.Println(err)
	}
	for _, v := range pros {
		if err := xlsx.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{v.Mark, v.Remark, v.Id, v.Img, v.CodeType, v.Code, v.Brands, v.Tags, v.Title, v.Rating, v.Comments, v.Price, v.Sellers, v.Distribution, v.Variants1, v.Variants2, v.VariantsId, v.ArrivalTime, v.Category1, v.Category2, v.Category3, v.Category4, v.Category5, v.Category6, v.Category7, v.CreateTime, v.UpdateTime}); err != nil {
			log.Println(err)
		}
		num++

	}
	fileName := "./file.xlsx"
	xlsx.SaveAs(fileName)
	//c.Header("Content-Length", "-1")
	//c.Header("Transfer-Encoding", "true")
	//c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName) // 用来指定下载下来的文件名
	//c.Header("Content-Transfer-Encoding", "binary")
	c.File("./file.xlsx")
}
func Dowwnload(c *gin.Context) {
	tools.SafeDeleteFile("./file.xlsx")
	var pros []mode.ProductDetails
	if len(config.DowwIds) != 0 {
		pros = product.GetProductId(tools.UniqueArr(config.DowwIds))
	}
	xlsx := excelize.NewFile()
	num := 2
	if err := xlsx.SetSheetRow("Sheet1", "A1", &[]interface{}{"标记", "备注", "id", "图片", "商品码类型", "商品码值", "品牌", "标签", "标题", "评分", "评论数量", "价格", "卖家", "配送", "变体1", "变体2", "变体id", "到达时间", "类目1", "类目2", "类目3", "类目4", "类目5", "类目6", "类目7", "创建时间", "更新时间"}); err != nil {
		log.Println(err)
	}
	for _, v := range pros {
		if err := xlsx.SetSheetRow("Sheet1", "A"+strconv.Itoa(num), &[]interface{}{v.Mark, v.Remark, v.Id, v.Img, v.CodeType, v.Code, v.Brands, v.Tags, v.Title, v.Rating, v.Comments, v.Price, v.Sellers, v.Distribution, v.Variants1, v.Variants2, v.VariantsId, v.ArrivalTime, v.Category1, v.Category2, v.Category3, v.Category4, v.Category5, v.Category6, v.Category7, v.CreateTime, v.UpdateTime}); err != nil {
			log.Println(err)
		}
		num++

	}
	fileName := "./file.xlsx"
	xlsx.SaveAs(fileName)
	//c.Header("Content-Length", "-1")
	//c.Header("Transfer-Encoding", "true")
	//c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName) // 用来指定下载下来的文件名
	//c.Header("Content-Transfer-Encoding", "binary")
	c.File("./file.xlsx")
}

func DowAdd(c *gin.Context) {
	id, _ := c.GetQuery("id")
	config.MuxDowIds.Lock()
	config.DowIds = append(config.DowIds, id)
	config.MuxDowIds.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})

}
func DowwAdd(c *gin.Context) {
	ids := c.PostForm("ids")
	split := strings.Split(ids, ",")
	config.MuxDowIds.Lock()
	config.DowwIds = []string{}
	for i := range split {
		config.DowwIds = append(config.DowwIds, split[i])
	}
	config.MuxDowIds.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})

}

func DowRemove(c *gin.Context) {
	config.MuxDowIds.Lock()
	config.DowIds = []string{}
	config.MuxDowIds.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})

}
func ShoDelMsg(c *gin.Context) {
	shopping.DelMsg()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})
}
func TuDelMsg(c *gin.Context) {
	timedUploads.DelMsg()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "提交成功",
	})
}

func ShoUpload(c *gin.Context) {
	var code int
	var msg string

	file, err := c.FormFile("file")
	msg = "上传成功 " + file.Filename
	if err != nil {
		log.Println(err)
		code = 1
		msg = err.Error()
	} else {
		open, err := file.Open()
		if err != nil {
			log.Println(err)
			code = 1
			msg = err.Error()
		} else {
			read := csv.NewReader(transform.NewReader(open, simplifiedchinese.GBK.NewDecoder()))
			for i := 0; true; i++ {
				split, err := read.Read()
				if i == 0 {
					continue
				}
				if len(split) > 0 {
					for i2 := range split {
						split[i2] = strings.Replace(strings.Replace(strings.Replace(split[i2], "/r", "", -1), "/n", "", -1), "，", ",", -1)
					}
					id, err := strconv.Atoi(split[0])
					if err != nil {
						msg = strconv.Itoa(i) + err.Error()
						log.Println(err)
						break
					}
					var space int
					if split[7] == "" {
						space = 159057
					} else {
						space, err = strconv.Atoi(strings.TrimSpace(split[7]))
						if err != nil {
							msg = strconv.Itoa(i) + err.Error()
							log.Println(err)
							break
						}
					}
					go shopping.AddShopping(mode.Shopping{PrId: id, Sku: split[1], FloorPrice: strings.TrimSpace(split[2]), XPrice: strings.TrimSpace(split[3]), Inventory: split[4], CenterId: split[5], Seller: split[6], Sales: space, Note: split[8], Name: split[9], ShoppingCron: split[10], TheShelvesCron: split[11], InventoryCron: split[12], XInventoryCron: split[13], StatusCron1: split[14], StatusCron2: split[15], StatusCron3: split[16], StatusCron4: split[17], StatusCron5: split[18]})
				}
				if err != nil {
					break
				}
			}
		}
	}
	err = shopping.CronRun()
	if err != nil {
		msg = "错误：" + err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": 0,
		"msg":  msg,
	})
}

func TuUpload(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// Create a destination file
	dst, err := os.Create("file/" + file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()

	// Copy the file contents to the destination
	_, err = io.Copy(dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "上传成功",
	})
}

func SAUpload(c *gin.Context) {
	var code int
	var msg string
	sasMut.Lock()
	defer func() {
		sasMut.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"data": 0,
			"msg":  msg,
		})
	}()
	file, err := c.FormFile("file")
	msg = "上传成功 " + file.Filename
	if err != nil {
		log.Println(err)
		code = 1
		msg = err.Error()
		return
	}
	open, err := file.Open()
	if err != nil {
		log.Println(1, err)
		code = 1
		msg = err.Error()
		return
	}
	xls, err := excelize.OpenReader(open)
	if err != nil {
		log.Println(2, err)
		code = 1
		msg = err.Error()
		return
	}
	// 读取工作表内容
	rows, err := xls.GetRows("Sheet1")
	if err != nil {
		log.Println(3, err)
		code = 1
		msg = err.Error()
		return
	}

	switch file.Filename {
	case "在库数量+发货仓和销售账号.xlsx":
		sas := make(map[string]mode.StockAvailability)
		rows2, err := xls.GetRows("Sheet2")
		if err != nil {
			log.Println(err)
			code = 1
			msg = err.Error()
			return
		}

		for i, row := range rows2 {
			if i == 0 {
				continue
			}
			if row[0] != "" {
				var remarks1, remarks2 string
				switch len(row) {
				case 4:
					remarks1 = row[3]
				case 5:
					remarks1 = row[3]
					remarks2 = row[4]
				}
				sas[row[0]] = mode.StockAvailability{ItemId: row[0], Warehouse: row[1], SalesUser: row[2], Remarks1: remarks1, Remarks2: remarks2}
			} else {
				break
			}
		}
		for i, row := range rows {
			if i == 0 {
				continue
			}
			salesUser := row[0]
			gtin := row[2]
			itemId := row[3]
			ptSKU := row[4]
			libraryNum, err := strconv.Atoi(strings.TrimSpace(row[8]))
			if err != nil {
				log.Println(i, err)
				libraryNum = 0
			}
			sa, ok := sas[itemId]
			if ok {
				sa.LibraryNum += int64(libraryNum)
				if sa.SalesUser == salesUser {
					sa.PtSku = ptSKU
					//log.Println(ptSKU)
				}
				sa.Gtin = gtin
				sas[itemId] = sa
			} else {
				sas[itemId] = mode.StockAvailability{ItemId: itemId, Gtin: gtin, LibraryNum: int64(libraryNum)}
			}
		}

		all, _ := stockAvailability.GetStockAvailabilityAll()
		var upsas []mode.StockAvailability
		var crsas []mode.StockAvailability
	ca:
		for i := range sas {
			for i2 := range all {
				if all[i2].ItemId == sas[i].ItemId {
					all[i2].Gtin = sas[i].Gtin
					all[i2].PtSku = sas[i].PtSku
					all[i2].LibraryNum = sas[i].LibraryNum
					all[i2].SalesUser = sas[i].SalesUser
					all[i2].Warehouse = sas[i].Warehouse
					all[i2].Remarks1 = sas[i].Remarks1
					all[i2].Remarks2 = sas[i].Remarks2
					upsas = append(upsas, all[i2])
					continue ca
				}
			}
			crsas = append(crsas, sas[i])
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "在库数量+发货仓和销售账号，更新数据失败！！！"
			return
		}
		is = stockAvailability.Install(crsas)
		if !is {
			msg = "在库数量+发货仓和销售账号，创建数据失败！！！"
			return
		}
		log.Println("在库数量+发货仓和销售账号完成")

	case "在途数量.xlsx":
		sas := make(map[string]mode.StockAvailability)
		for i, row := range rows {
			if i == 0 {
				continue
			}
			itemId := row[0]
			transitNum, err := strconv.Atoi(strings.TrimSpace(row[1]))
			if err != nil {
				log.Println(err)
				transitNum = 0
			}
			sa, ok := sas[itemId]
			if ok {
				sa.TransitNum += int64(transitNum)
				sas[itemId] = sa
			} else {
				sas[itemId] = mode.StockAvailability{ItemId: itemId, TransitNum: int64(transitNum)}
			}
		}
		all, _ := stockAvailability.GetStockAvailabilityAll()
		var upsas []mode.StockAvailability
		var crsas []mode.StockAvailability
	ca2:
		for i := range sas {
			for i2 := range all {
				if all[i2].ItemId == sas[i].ItemId {
					all[i2].TransitNum = sas[i].TransitNum
					upsas = append(upsas, all[i2])
					continue ca2
				}
			}
			crsas = append(crsas, sas[i])
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "在途数量，更新数据失败！！！"
			return
		}
		is = stockAvailability.Install(crsas)
		if !is {
			msg = "在途数量，创建数据失败！！！"
			return
		}
		log.Println("在途数量完成")

	case "15天30天60天销量.xlsx":
		sas := make(map[string]map[string]int64)
		for i, row := range rows {
			if i == 0 {
				continue
			}
			ts := strings.TrimSpace(row[0])
			itemId := row[1]
			num, err := strconv.Atoi(strings.TrimSpace(row[15]))
			if err != nil {
				log.Println(err)
				num = 0
			}
			_, ok := sas[itemId]
			if ok {
				_, ok2 := sas[itemId][ts]
				if ok2 {
					sas[itemId][ts] += int64(num)
				} else {
					sas[itemId][ts] = int64(num)
				}
			} else {
				sas[itemId] = make(map[string]int64)
				sas[itemId][ts] = int64(num)
			}

		}
		all, _ := stockAvailability.GetStockAvailabilityAll()
		average, counts := stockAvailability.WeightedDailyAverage(sas)

		var upsas []mode.StockAvailability
		var crsas []mode.StockAvailability
	ca3:
		for i := range average {
			for i2 := range all {
				if all[i2].ItemId == i {
					all[i2].Weighted = average[i]
					all[i2].Counts = counts[i]
					upsas = append(upsas, all[i2])
					continue ca3
				}
			}
			crsas = append(crsas, mode.StockAvailability{ItemId: i, Weighted: average[i], Counts: counts[i]})
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "15天30天60天销量，更新数据失败！！！"
			return
		}
		is = stockAvailability.Install(crsas)
		if !is {
			msg = "15天30天60天销量，创建数据失败！！！"
			return
		}
		log.Println("15天30天60天销量完成")

	case "三方仓-产品编码.xlsx":
		sas := make(map[string]mode.StockAvailability)
		for i, row := range rows {
			if i == 0 {
				continue
			}
			cySku := row[2]
			gtin := row[4]
			sas[gtin] = mode.StockAvailability{Gtin: gtin, CySku: cySku}
		}
		all, _ := stockAvailability.GetStockAvailabilityAll()
		var upsas []mode.StockAvailability
	ca4:
		for i := range sas {
			for i2 := range all {
				if all[i2].Gtin == sas[i].Gtin {
					all[i2].CySku = sas[i].CySku
					upsas = append(upsas, all[i2])
					continue ca4
				}
			}
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "三方仓-产品编码，更新数据失败！！！"
			return
		}
		log.Println("三方仓-产品编码完成")

	case "易仓库存.xlsx":
		sas := make(map[string]mode.StockAvailability)
		for i, row := range rows {
			if i == 0 {
				continue
			}
			cySku := row[0]
			cyName := row[1]
			declaration := row[2]
			sas[cySku] = mode.StockAvailability{CySku: cySku, CyName: cyName, Declaration: declaration}
		}

		all, _ := stockAvailability.GetStockAvailabilityAll()
		var upsas []mode.StockAvailability
		for i := range sas {
			//fmt.Println(sas[i].CySku)
			for i2 := range all {
				if all[i2].CySku == sas[i].CySku {
					all[i2].CyName = sas[i].CyName
					all[i2].Declaration = sas[i].Declaration
					upsas = append(upsas, all[i2])
				}
			}
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "易仓库存，更新数据失败！！！"
			return
		}
		log.Println("易仓库存完成")

	case "备货天数.xlsx":
		sas := make(map[string]mode.StockAvailability)
		for i, row := range rows {
			if i == 0 {
				continue
			}
			itemId := row[0]
			leadTime, err := strconv.Atoi(strings.TrimSpace(row[1]))
			if err != nil {
				log.Println(err)
				leadTime = 0
			}
			sa, ok := sas[itemId]
			if ok {
				sa.LeadTime += int64(leadTime)
				sas[itemId] = sa
			} else {
				sas[itemId] = mode.StockAvailability{ItemId: itemId, LeadTime: int64(leadTime)}
			}
		}

		all, _ := stockAvailability.GetStockAvailabilityAll()
		var upsas []mode.StockAvailability
		var crsas []mode.StockAvailability
	ca6:
		for i := range sas {
			for i2 := range all {
				if all[i2].ItemId == sas[i].ItemId {
					all[i2].LeadTime = sas[i].LeadTime
					upsas = append(upsas, all[i2])
					continue ca6
				}
			}
			crsas = append(crsas, sas[i])
		}
		is := stockAvailability.UpdateStockAvailability(stockAvailability.CalculateStockingQuantity(upsas))
		if !is {
			msg = "备货天数，更新数据失败！！！"
			return
		}
		is = stockAvailability.Install(crsas)
		if !is {
			msg = "备货天数，创建数据失败！！！"
			return
		}
		log.Println("备货天数完成")
	}

}

func SADel(c *gin.Context) {
	stockAvailability.DelTockAvailabilityMsg()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "清除成功",
	})
}
func DowImg(c *gin.Context) {
	all, _ := stockAvailability.GetStockAvailabilityAll()
	var list []mode.StockAvailability
	for i := range all {
		if all[i].Img == "" {
			list = append(list, all[i])
		}
	}
	go stockAvailability.DowImg(list)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  "执行成功",
	})
}

func SADow(c *gin.Context) {
	// 创建一个内存缓冲区
	buffer := new(bytes.Buffer)

	// 创建一个新的XLSX文件
	file := excelize.NewFile()

	// 在第一个工作表中写入一些数据
	file.SetCellValue("Sheet1", "A1", "销售账号")
	file.SetCellValue("Sheet1", "B1", "ITEM ID")
	file.SetCellValue("Sheet1", "C1", "图片")
	file.SetCellValue("Sheet1", "D1", "易仓SKU")
	file.SetCellValue("Sheet1", "E1", "产品名（中文）")
	file.SetCellValue("Sheet1", "F1", "GTIN")
	file.SetCellValue("Sheet1", "G1", "平台SKU")
	file.SetCellValue("Sheet1", "H1", "产品名（英文）")
	file.SetCellValue("Sheet1", "I1", "数量")
	file.SetCellValue("Sheet1", "J1", "发货仓库")
	file.SetCellValue("Sheet1", "K1", "备注1")
	file.SetCellValue("Sheet1", "L1", "备注2")
	file.SetCellValue("Sheet1", "M1", "备货天数")
	file.SetCellValue("Sheet1", "N1", "在途数量")
	file.SetCellValue("Sheet1", "O1", "在库数量")
	file.SetCellValue("Sheet1", "P1", "加权日均")

	all, _ := stockAvailability.GetStockAvailabilityAll()
	for i := range all {
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), all[i].SalesUser)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), all[i].ItemId)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), all[i].Img)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), all[i].CySku)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), all[i].CyName)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), all[i].Gtin)
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), all[i].PtSku)
		file.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), all[i].Declaration)
		file.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), all[i].Num)
		file.SetCellValue("Sheet1", "J"+strconv.Itoa(i+2), all[i].Warehouse)
		file.SetCellValue("Sheet1", "K"+strconv.Itoa(i+2), all[i].Remarks1)
		file.SetCellValue("Sheet1", "L"+strconv.Itoa(i+2), all[i].Remarks2)
		file.SetCellValue("Sheet1", "M"+strconv.Itoa(i+2), all[i].LeadTime)
		file.SetCellValue("Sheet1", "N"+strconv.Itoa(i+2), all[i].TransitNum)
		file.SetCellValue("Sheet1", "O"+strconv.Itoa(i+2), all[i].LibraryNum)
		file.SetCellValue("Sheet1", "P"+strconv.Itoa(i+2), all[i].Weighted)
	}

	// 将XLSX文件写入内存缓冲区
	if err := file.Write(buffer); err != nil {
		c.String(http.StatusInternalServerError, "写入XLSX文件到缓冲区失败!")
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename=算备货销量导出数据.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

func SIUpload(c *gin.Context) {
	var code int
	var msg string
	sasMut.Lock()
	defer func() {
		sasMut.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"data": 0,
			"msg":  msg,
		})
	}()
	file, err := c.FormFile("file")
	msg = "上传成功 " + file.Filename
	if err != nil {
		log.Println(err)
		code = 1
		msg = err.Error()
		return
	}
	open, err := file.Open()
	if err != nil {
		log.Println(1, err)
		code = 1
		msg = err.Error()
		return
	}
	xls, err := excelize.OpenReader(open)
	if err != nil {
		log.Println(2, err)
		code = 1
		msg = err.Error()
		return
	}
	// 读取工作表内容
	rows, err := xls.GetRows("Sheet1")
	if err != nil {
		log.Println(3, err)
		code = 1
		msg = err.Error()
		return
	}

	switch file.Filename {
	case "店铺信息ID.xlsx":
		var sis []mode.StoreInformation
		for i, row := range rows {
			if i == 0 {
				continue
			}
			sa := mode.StoreInformation{}
			if len(row) < 1 {
				continue
			}
			sa.AccountShopName = strings.TrimSpace(row[0])
			if len(row) >= 2 {
				sa.ShopID, err = strconv.Atoi(strings.TrimSpace(row[1]))
				if err != nil {
					log.Println(err)
					sa.ShopID = 0
					code = 1
					msg = err.Error()
				}
			}
			if len(row) >= 3 {
				sa.PID, err = strconv.Atoi(strings.TrimSpace(row[2]))
				if err != nil {
					log.Println(err)
					sa.PID = 0
					code = 1
					msg = err.Error()
				}
			}

			sis = append(sis, sa)
		}
		//fmt.Println(sis)
		storeInformation.Install(sis)

	case "店铺信息备注.xlsx":
		var sis []mode.StoreInformation
		for i, row := range rows {
			if i == 0 {
				continue
			}
			if len(row) < 1 {
				continue
			}
			sa := mode.StoreInformation{}
			if len(row) >= 1 && row[0] != "" {
				sa.ShopID, err = strconv.Atoi(strings.TrimSpace(row[0]))
				if err != nil {
					log.Println(err)
					sa.ShopID = 0
					code = 1
					msg = err.Error()
				}
			}
			if len(row) >= 2 {
				sa.AffiliatedCompany = row[1]
			}
			if len(row) >= 3 {
				sa.Note1 = row[2]
			}
			if len(row) >= 4 {
				sa.Note2 = row[3]
			}
			if len(row) >= 5 {
				sa.Note3 = row[4]
			}
			if len(row) >= 6 {
				sa.Note4 = row[5]
			}
			if len(row) >= 7 {
				sa.Note5 = row[6]
			}
			sis = append(sis, sa)
		}
		storeInformation.UpdateStoreInformationNoteAndAffiliatedCompany(sis)

	}

}

func SIUpdateDate(c *gin.Context) {
	all, _ := storeInformation.GetStoreInformationAll()
	msg := "开始执行更新"
	if storeInformation.Mux.TryLock() {
		go storeInformation.UpdateData(all)
	} else {
		msg = "更新执行中"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  msg,
	})
}

func SIInfo(c *gin.Context) {
	msg := "开始获取卖家信息"
	if storeInformation.MuxInfo.TryLock() {
		go func() {
			defer storeInformation.MuxInfo.Unlock()
			all, i := storeInformation.GetStoreInformationAll()
			storeInformation.InfoNum = i
			for i := range all {

				if all[i].Country == "" && all[i].ShopID != 0 {
					storeInformation.Wait.Add(1)
					config.Ch <- 1
					go storeInformation.CrawlerInfo(strconv.Itoa(all[i].ShopID))
				} else {
					storeInformation.InfoNum--
				}
			}
			storeInformation.Wait.Wait()

		}()
	} else {
		msg = "卖家信息获取中执行中"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": 0,
		"msg":  msg,
	})
}

func EditSi(c *gin.Context) {
	sid := c.DefaultPostForm("Sid", "")
	pid := c.DefaultPostForm("PID", "")
	shopID := c.DefaultPostForm("ShopID", "")
	AffiliatedCompany := c.DefaultPostForm("AffiliatedCompany", "")

	Note1 := c.DefaultPostForm("Note1", "")
	Note2 := c.DefaultPostForm("Note2", "")
	Note3 := c.DefaultPostForm("Note3", "")
	Note4 := c.DefaultPostForm("Note4", "")
	Note5 := c.DefaultPostForm("Note5", "")
	atoi, _ := strconv.Atoi(sid)
	atoi2, _ := strconv.Atoi(shopID)
	atoi3, _ := strconv.Atoi(pid)
	count := storeInformation.UpdateStoreInformationNote([]mode.StoreInformation{{PID: atoi3, ShopID: atoi2, SID: atoi, AffiliatedCompany: AffiliatedCompany, Note1: Note1, Note2: Note2, Note3: Note3, Note4: Note4, Note5: Note5}})
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "修改成功",
	})
}
func AddSi(c *gin.Context) {
	accountShopName := c.DefaultPostForm("AccountShopName", "")
	if accountShopName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": 0,
			"msg":  "店铺名称不能为空",
		})
		return
	}
	_, i := storeInformation.GetStoreInformation(1, 1, accountShopName, "", "", "", "")
	if i > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": 0,
			"msg":  "店铺名称已存在",
		})
		return
	}
	pid := c.DefaultPostForm("PID", "")
	shopID := c.DefaultPostForm("ShopID", "")
	AffiliatedCompany := c.DefaultPostForm("AffiliatedCompany", "")

	Note1 := c.DefaultPostForm("Note1", "")
	Note2 := c.DefaultPostForm("Note2", "")
	Note3 := c.DefaultPostForm("Note3", "")
	Note4 := c.DefaultPostForm("Note4", "")
	Note5 := c.DefaultPostForm("Note5", "")
	atoi2, _ := strconv.Atoi(shopID)
	atoi3, _ := strconv.Atoi(pid)
	count := storeInformation.Install([]mode.StoreInformation{{PID: atoi3, ShopID: atoi2, AccountShopName: accountShopName, AffiliatedCompany: AffiliatedCompany, Note1: Note1, Note2: Note2, Note3: Note3, Note4: Note4, Note5: Note5}})
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "添加成功",
	})
}

func GetProductSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	itemId := c.DefaultQuery("itemId", "")

	sas, count := productSales.SelectSales(page, limit, itemId)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": count,
		"data":  sas,
		"msg":   "成功",
	})
}

func AddProductSales(c *gin.Context) {
	sales := c.PostForm("sales")
	sales = strings.TrimSpace(sales)
	// 获取当前时间
	now := time.Now()
	// 将时间戳截断为天
	truncatedTime := now.Truncate(24 * time.Hour)
	// 输出精确到天的时间
	day := truncatedTime.Format("2006-01-02")
	split := strings.Split(sales, "\n")
	var saless []mode.ProductSales
	for i := range split {
		i2 := strings.Split(split[i], ",")
		atoi, _ := strconv.Atoi(i2[0])
		saless = append(saless, mode.ProductSales{ITEMID: atoi, CatalogSellerId: i2[1], CreateDate: day})
	}
	productSales.AddSales(saless)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
	})
}

func SalesRemove(c *gin.Context) {
	names := c.PostForm("ids")
	count := productSales.RemoveSales(names)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": count,
		"msg":  "删除成功",
	})
}

func ProductSalesUpdateDate(c *gin.Context) {
	ids := c.PostForm("ids")
	ids = strings.TrimSpace(ids)
	saless, _ := productSales.SelectSales(1, 1000000, ids)
	if productSales.Mux.TryLock() {
		go func() {
			log.Printf("开始爬取库存更新销量 数量：%d", len(saless))
			defer productSales.Mux.Unlock()
			productSales.Num = len(saless)
			for i := range saless {
				config.Ch <- 1
				productSales.Wg.Add(1)
				go productSales.CrawlerSales(saless[i])
			}
			productSales.Wg.Wait()
		}()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "执行中",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "开始执行",
	})
}
