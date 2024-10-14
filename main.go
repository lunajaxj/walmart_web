package main

import (
	"embed"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"walmart_web/app/config"
	"walmart_web/app/controller"
)

//go:embed static/* templates/*
var f embed.FS

func main() {
	// 设置全局日志格式
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("欢迎使用 DISEN 沃尔玛系统")
	gin.SetMode(gin.ReleaseMode)
	// 禁止Gin的控制台输出
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	// 初始化session存储器
	store := cookie.NewStore([]byte("your-secret-key"))
	r.Use(sessions.Sessions("session-name", store))
	//if config.Mode != 0 {
	//	r.Use(config.RuleConfig)
	//}

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Delims("[[", "]]")
	r.StaticFS("/public", http.FS(f))
	r.SetHTMLTemplate(template.Must(template.New("").Delims("[[", "]]").ParseFS(f, "templates/*.html")))
	r.RedirectTrailingSlash = false
	//首页
	r.GET("/admin/index", controller.Index)
	//首页
	r.GET("/", controller.Home)
	//登录
	r.GET("/login/index", controller.LoginIndex)
	//商品页面
	r.GET("/admin/product/product-table", controller.Table)
	//商品修改页面
	r.GET("/admin/product/edit", controller.Edit)
	//图表页面
	r.GET("/admin/product/chart", controller.Chart)

	//抢购物车页面
	r.GET("/admin/shopping/product-shopping", controller.Shopping)
	//抢购物车修改页面
	r.GET("/admin/shopping/edit", controller.EditShopping)

	//固定任务页面
	r.GET("/admin/timedUploads/timed-uploads", controller.TimedUploads)
	// 更新店铺名称映射
	r.POST("/api/timedUploads/updateShopMapping", controller.UpdateShopMapping)
	//固定任务修改页面
	r.GET("/admin/timedUploads/edit", controller.EditTimedUploads)
	//固定任务创建页面
	r.GET("/admin/timedUploads/create", controller.CreateTimedUploads)
	//固定任务文件页面
	r.GET("/admin/timedUploads/file", controller.File)

	//备货页面
	r.GET("/admin/stockAvailability/stock-availability", controller.StockAvailability)

	//邮件页面
	r.GET("/admin/mail/mail-table", controller.Mail)

	//类目抓取页面
	r.GET("/admin/product/crawler-edit", controller.CrawlerEdit)

	//卖家信息页面
	r.GET("/admin/storeInformation/store-information", controller.StoreInformation)
	//修改卖家信息
	r.GET("/admin/storeInformation/edit", controller.EditStoreInformation)
	//添加卖家信息
	r.GET("/admin/storeInformation/add", controller.AddStoreInformation)

	//库存销量页面
	r.GET("/admin/productSales/product-sales", controller.ProductSales)

	//获取日志
	r.GET("/api/log/getLogs", controller.GetLogs)
	//获取商品
	r.GET("/api/product/getProduct", controller.GetProduct)
	//登录
	r.POST("/api/login", controller.Login)
	//注销
	r.GET("/api/logout", controller.Logout)
	//获取类目
	r.POST("/api/category/getCategory", controller.GetCategory)
	//抓取id
	r.POST("/api/product/crawlerId", controller.CrawlerId)
	//更新id
	r.POST("/api/product/crawlerIdUp", controller.CrawlerIdUp)
	//抓取key
	r.POST("/api/product/crawlerKey", controller.CrawlerKey)
	//抓取url
	r.POST("/api/product/crawlerUrl", controller.CrawlerUrl)
	//抓取关键词
	r.POST("/api/product/crawlerAct", controller.CrawlerAct)
	//抓取类目
	r.POST("/api/product/crawlerCategory", controller.CrawlerCategory)

	//设置代理频率
	r.POST("/api/config/setCh", controller.SetCh)
	//修改代理状态
	r.GET("/api/config/upIsRun", controller.UpdateIsRun)
	//删除缓存File
	r.GET("/api/config/delKeys", controller.DelFile)
	//删除缓存下载
	r.GET("/api/config/delDow", controller.DowRemove)
	//加入下载队列
	r.GET("/api/config/addDow", controller.DowAdd)
	//下载
	r.GET("/api/config/dow", controller.Download)

	//加入导出列表
	r.POST("/api/config/addDoww", controller.DowwAdd)
	//导出
	r.GET("/api/config/doww", controller.Dowwnload)

	//修改标记和备注
	r.POST("/api/product/editProduct", controller.EditProduct)
	//更新gtin
	r.POST("/api/product/uploadGtin", controller.UploadGtin)

	//获取购物车信息
	r.GET("/api/shopping/getShopping", controller.GetShopping)
	//修改购物车信息
	r.POST("/api/shopping/editShopping", controller.EditSho)
	//删除购物车
	r.POST("/api/shopping/remove", controller.ShoppingRemove)
	//文件上传购物车
	r.POST("/api/shopping/upload", controller.ShoUpload)
	//更新购物车信息
	r.POST("/api/shopping/upIds", controller.UpIds)
	//删除信息
	r.GET("/api/shopping/delMsg", controller.ShoDelMsg)

	//获取固定任务信息
	r.GET("/api/timedUploads/getTimedUploads", controller.GetTimedUploads)
	//修改固定任务信息
	r.POST("/api/timedUploads/editTimedUploads", controller.EditTu)
	//创建固定任务信息
	r.POST("/api/timedUploads/createTimedUploads", controller.CreateTu)
	//删除固定任务
	r.POST("/api/timedUploads/remove", controller.TimedUploadsRemove)
	//获取表格文件
	r.GET("/api/timedUploads/getFiles", controller.GetFile)
	//删除表格文件
	r.POST("/api/timedUploads/removeFile", controller.RemoveFile)
	//文件上传固定任务
	r.POST("/api/timedUploads/upload", controller.TuUpload)
	//删除固定任务信息
	r.GET("/api/timedUploads/delMsg", controller.TuDelMsg)

	//获取mail信息
	r.GET("/api/mail/getMails", controller.GetMails)
	//删除mail信息
	r.GET("/api/mail/delMsg", controller.MailDelMsg)
	//执行mail
	r.GET("/api/mail/run", controller.MailRun)

	//获取stockAvailability信息
	r.GET("/api/stockAvailability/getStockAvailability", controller.GetStockAvailability)
	//上传stockAvailability信息
	r.POST("/api/stockAvailability/upload", controller.SAUpload)
	r.GET("/api/stockAvailability/dowImg", controller.DowImg)

	//清除stockAvailability信息
	r.GET("/api/stockAvailability/dow", controller.SADow)
	//清除stockAvailability信息
	r.GET("/api/stockAvailability/delMsg", controller.SADel)
	//删除stockAvailability
	r.POST("/api/stockAvailability/remove", controller.SARemove)

	//获取storeInformation信息
	r.GET("/api/storeInformation/getStoreInformation", controller.GetStoreInformation)

	//上传storeInformation信息
	r.POST("/api/storeInformation/upload", controller.SIUpload)

	//修改storeInformation信息
	r.POST("/api/storeInformation/update", controller.EditSi)

	//添加storeInformation信息
	r.POST("/api/storeInformation/add", controller.AddSi)

	//更新storeInformation信息
	r.GET("/api/storeInformation/updateDate", controller.SIUpdateDate)

	//抓取storeInformation信息
	r.GET("/api/storeInformation/crawlerInfo", controller.SIInfo)

	//删除storeInformation
	r.POST("/api/storeInformation/remove", controller.SIRemove)

	//获取productSales信息
	r.GET("/api/productSales/getProductSales", controller.GetProductSales)

	//执行获取
	r.POST("/api/productSales/updateDate", controller.ProductSalesUpdateDate)
	//新增
	r.POST("/api/productSales/add", controller.AddProductSales)

	//删除productSales
	r.POST("/api/productSales/remove", controller.SalesRemove)

	go func() {
		log.Println(http.ListenAndServe("192.168.2.8:8092", nil))
	}()
	if config.Mode == 0 {
		r.Run(":80")
	}
	if config.Mode == 1 {
		r.Run(":80")
	}
	if config.Mode == 2 {
		r.Run(":8080")
	}

}
