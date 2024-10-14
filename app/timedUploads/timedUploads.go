package timedUploads

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/gammazero/workerpool"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	_ "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net"
	"time"
	"walmart_web/app/config"
	"walmart_web/app/mode"
)

// 创建支持秒字段的 cron 解析器
var cronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

var (
	// 将 pid 作为关键字来映射用户信息
	users = map[string][]string{
		"10001468582": {"yumumandzsw@163.com", "Dyumuman123!"},
		"10001295723": {"guibinwalmart@163.com", "DGuibin123!"},
		"10001447777": {"ziyuewlkj@163.com", "Dziyue123!"},
		"10001659525": {"keyshinewalmart@163.com", "&Keyshine788"},
		"10001691058": {"moremuma@163.com", "$Moremuma778"},
	}

	// 将 pid 作为关键字来映射端口号
	se = map[string]string{
		"10001468582": "6222",  // Online mini-mart
		"10001447777": "7222",  // Money Saving World
		"10001295723": "8222",  // GUBIN
		"10001659525": "12222", // Juno Kael
		"10001691058": "11222", // Moremuma
	}

	// 初始化 cron.Cron 实例时，指定使用自定义解析器
	crs = map[string]*cron.Cron{
		"10001468582": cron.New(cron.WithSeconds()), // Online mini-mart
		"10001295723": cron.New(cron.WithSeconds()), // GUBIN
		"10001447777": cron.New(cron.WithSeconds()), // Money Saving World
		"10001659525": cron.New(cron.WithSeconds()), // Juno Kael
		"10001691058": cron.New(cron.WithSeconds()), // Moremuma
	}

	// 任务执行状态
	IsRun = map[string]bool{
		"10001468582": false, // Online mini-mart
		"10001295723": false, // GUBIN
		"10001447777": false, // Money Saving World
		"10001659525": false, // Juno Kael
		"10001691058": false, // Moremuma
	}
	//ch    = make(chan int, 6)
	MsgCh string
	isCh  bool

	wp = workerpool.New(1)
)

func init() {
	LoadShopMappings()                       // 加载店铺映射
	InitAndCheckPendingTasks()               // 启动时检查未完成任务
	CronRun()                                // 加载并启动定时任务 	// 启动任务执行器
	StartPollingFailedTasks(1 * time.Minute) // 启动轮询任务，检查并重新执行排队中的任务
	log.Println("端口映射状态:", se)               // 打印se字典
}

// LoadShopMappings 加载店铺映射
func LoadShopMappings() map[string]string {
	var shopMappings []mode.ShopMapping
	shopMap := make(map[string]string)

	// 从数据库中获取 pid 和 shop_name
	err := config.Db.Select(&shopMappings, "SELECT pid, shop_name FROM shop_mapping")
	if err != nil {
		log.Println("获取店铺映射失败:", err)
		return shopMap
	}

	// 将 pid 和 shop_name 映射
	for _, mapping := range shopMappings {
		shopMap[mapping.PID] = mapping.ShopName
		log.Printf("店铺映射: PID=%s, ShopName=%s", mapping.PID, mapping.ShopName)
		// 打印出每个店铺的cron表达式
		printCronExpressions(mapping.PID, mapping.ShopName)
	}

	return shopMap
}

// 打印出每个店铺的cron表达式
func printCronExpressions(pid, shopName string) {
	// 假设有一个获取 cron 表达式的方法 GetTimedUploadsCron
	genres := []string{"价格", "库存", "类别"} // 可以根据实际任务类型调整
	for _, genre := range genres {
		cronExprs := GetTimedUploadsCron(genre, shopName)
		for _, cronExpr := range cronExprs {
			log.Printf("店铺 %s (PID: %s) 的 %s 任务 cron 表达式: %s", shopName, pid, genre, cronExpr)
		}
	}
}

func isPortOpen(host string, port string) bool {
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func InitAndCheckPendingTasks() {
	var pendingTasks []mode.TimedUploads
	err := config.Db.Select(&pendingTasks, "SELECT * FROM timed_uploads WHERE msg NOT LIKE '%完成%'")
	if err != nil {
		log.Println("启动时查询未完成任务失败：", err)
		return
	}

	now := time.Now().Truncate(time.Second)

	for _, task := range pendingTasks {
		schedule, err := cronParser.Parse(task.Cron) // 使用自定义解析器解析 cron 表达式
		if err != nil {
			log.Printf("解析 cron 表达式失败: %s, 错误: %v", task.Cron, err)
			continue
		}

		// 直接比较当前时间和 cron 表达式
		if isTimeMatch(schedule, now) { // 当前时间是否匹配 cron 表达式
			task.Msg = "处理中..."
			UploadTimedUploads(task)

			wp.Submit(func() {
				executeTask(task.Genre, task.Cron, task.Seller)
			})
		} else {
			log.Printf("任务: %s 的当前时间 (%s) 与 cron 表达式 (%s) 不匹配，因此不执行", task.Name, now.Format(time.RFC3339), task.Cron)
		}
	}
}

// 轮询数据库，处理排队中的任务
func StartPollingFailedTasks(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			checkAndExecutePendingTasks()
		}
	}()
}

func CronRun() error {
	if config.Mode == 1 {
		return nil // 如果模式为 1，跳过任务加载
	}
	// 使用支持可选秒字段的解析器来初始化 cron 实例

	// 从数据库加载动态店铺映射
	shopMap := LoadShopMappings()

	log.Println("定时任务cron加载...")

	// 遍历每个店铺的 pid 和 shop_name
	for pid, shopName := range shopMap {
		// 停止当前店铺的定时任务
		crs[pid].Stop()
		time.Sleep(1 * time.Second) // 延迟 1 秒以避免冲突

		loadAndScheduleTasksForShop("价格", pid, shopName)
		loadAndScheduleTasksForShop("库存", pid, shopName)
		loadAndScheduleTasksForShop("类别", pid, shopName)

		// 启动 cron 任务
		crs[pid].Start()
		log.Printf("启动定时任务: 店铺 %s (PID: %s)", shopName, pid)
	}

	log.Println("定时任务加载完成。")
	return nil
}

func loadAndScheduleTasksForShop(genre, pid, shopName string) {
	// 获取该店铺的 cron 表达式列表
	cronExprs := GetTimedUploadsCron(genre, pid)
	if len(cronExprs) == 0 {
		log.Printf("店铺 %s 的 %s 任务没有找到任何 Cron 表达式", shopName, genre)
		return
	}

	// 添加每个任务的定时调度
	for _, cronExpr := range cronExprs {
		cronExprCopy := cronExpr // 创建一个局部变量的副本
		log.Printf("为店铺 %s 的 %s 任务添加 cron 表达式: %s", shopName, genre, cronExprCopy)

		// 调度任务
		entryID, err := crs[pid].AddFunc(cronExprCopy, func() {
			//log.Println("cronExpr=", cronExprCopy) // 使用局部变量的副本
			wp.Submit(func() { executeCronTask(genre, cronExprCopy, shopName) })
		})

		if err != nil {
			log.Printf("添加 %s 任务的 cron 表达式时出错: %v", genre, err)
		} else {
			log.Printf("任务已成功添加，任务 ID: %d,cronExpr:%s,shopName:%s", entryID, cronExprCopy, shopName)
		}
	}

	log.Printf("店铺 %s 的 %s 任务全部添加完成", shopName, genre)
}

func executeCronTask(genre, cronExpr, shopName string) {
	// 使用支持秒解析的 cron 调度器
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析 cron 表达式
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		log.Printf("解析 cron 表达式失败: %v", err)
		return
	}

	// 获取当前时间
	now := time.Now().Truncate(time.Second)
	log.Println("cronExpr=", cronExpr, "now=", now, "schedule=", schedule, "isTimeMatch=", isTimeMatch(schedule, now))

	// 检查当前时间是否符合 cron 表达式
	if isTimeMatch(schedule, now) {
		log.Printf("到时间执行任务: %s", cronExpr)
		executeTaskWithRetry(genre, cronExpr, shopName)
	} else {

		log.Printf("任务时间未到，任务不执行: %s, %s", cronExpr, shopName)
	}
}

func isTimeMatch(schedule cron.Schedule, now time.Time) bool {
	// 获取当前时间的各个字段
	second, minute, hour, day, month, weekday := now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month()), int(now.Weekday())

	// 将 cron 表达式解析为 SpecSchedule
	spec := schedule.(*cron.SpecSchedule)

	// 检查是否每个字段匹配
	if (1<<uint(second))&spec.Second == 0 {
		log.Println("秒不匹配")
		return false
	}
	if (1<<uint(minute))&spec.Minute == 0 {
		log.Println("分钟不匹配")
		return false
	}
	if (1<<uint(hour))&spec.Hour == 0 {
		log.Println("小时不匹配")
		return false
	}
	if (1<<uint(day))&spec.Dom == 0 {
		log.Println("日不匹配")
		return false
	}
	if (1<<uint(month))&spec.Month == 0 {
		log.Println("月不匹配")
		return false
	}
	if (1<<uint(weekday))&spec.Dow == 0 {
		log.Println("星期不匹配")
		return false
	}

	// 如果所有字段都匹配，则返回 true
	return true
}

func executeTaskWithRetry(genre, cronExpr, shopName string) {
	log.Printf("开始执行任务，店铺: %s, 类型: %s, cron 表达式: %s", shopName, genre, cronExpr)

	// 假设任务重试次数为 2 次
	for attempt := 1; attempt < 3; attempt++ {
		log.Printf("开始第 %d 次尝试，店铺: %s, 任务: %s", attempt, shopName, genre)

		isCh := executeTask(genre, cronExpr, shopName) // 修改为函数返回值更新 isCh

		// 如果任务成功，退出重试循环
		if isCh {
			log.Printf("任务执行成功，店铺: %s, 任务: %s", shopName, genre)
			break
		}

		// 如果没有成功，记录重试
		log.Printf("任务执行失败，正在重试第 %d 次，店铺: %s, 任务: %s", attempt, shopName, genre)
	}

	if !isCh {
		log.Printf("任务执行失败，已达到最大重试次数，店铺: %s, 任务: %s", shopName, genre)
	} else {
		log.Printf("任务执行完成，店铺: %s, 任务: %s", shopName, genre)
	}
}

// 检查并执行排队中的任务
func checkAndExecutePendingTasks() {
	var pendingTasks []mode.TimedUploads
	err := config.Db.Select(&pendingTasks, "SELECT * FROM timed_uploads WHERE msg LIKE '撞车了%' ORDER BY cron")
	if err != nil {
		log.Println("查询排队中的任务失败:", err)
		return
	}

	for _, task := range pendingTasks {
		if IsRun[task.Seller] {
			// 如果当前商家有任务在运行，则跳过
			continue
		}

		// 使用自定义解析器来解析 cron 表达式
		schedule, err := cronParser.Parse(task.Cron)
		if err != nil {
			log.Printf("解析任务的 cron 表达式失败: %s, 错误: %v", task.Cron, err)
			continue
		}

		nextRun := schedule.Next(time.Now())
		if nextRun.After(time.Now()) && nextRun.Sub(time.Now()) > 30*time.Minute {
			IsRun[task.Seller] = true
			wp.Submit(func() {
				defer func() {
					IsRun[task.Seller] = false
					checkAndExecutePendingTasks() // 任务完成后检查下一个排队任务
				}()
				executeTask(task.Genre, task.Cron, task.Seller)
			})
			break // 只执行一个任务
		} else {
			log.Printf("任务 %s 将跳过，因为下一个任务时间小于30分钟", task.Name)
		}
	}
}

// 执行任务的逻辑
func executeTask(genre, cronExpr, shopName string) bool {
	log.Printf("开始处理任务: %s for 店铺: %s", genre, shopName)

	// 获取店铺名对应的 PID
	pid := GetPIDByShopName(shopName)
	if pid == "" {
		log.Printf("未找到店铺 %s 对应的 PID", shopName)
		return false
	}

	// 根据 pid 和 cron 获取任务
	tasks := GetTimedUploadsByCron(genre, cronExpr, pid)
	if len(tasks) == 0 {
		log.Printf("未找到任何任务 for 店铺: %s (PID: %s)", shopName, pid)
		return false
	}

	// 打印任务信息
	log.Printf("开始处理任务: %s, 店铺: %s (PID: %s)", tasks[0].Name, shopName, pid)

	// 获取端口号
	port, ok := se[pid]
	if !ok || !isPortOpen("127.0.0.1", port) {
		// 如果端口不存在或未打开，记录日志并更新任务状态
		tasks[0].Msg = fmt.Sprintf("浏览器端口号 %s 未获取到", port)
		log.Printf("端口 %s 未打开，任务 %s 无法执行", port, tasks[0].Name)
		UploadTimedUploads(tasks[0])
		return false
	}

	// 执行相应的任务
	var success bool
	switch genre {
	case "价格":
		success = priceRun([]mode.TimedUploads{tasks[0]})
	case "库存":
		success = inventoryRun([]mode.TimedUploads{tasks[0]})
	case "类别":
		success = categoryRun([]mode.TimedUploads{tasks[0]})
	default:
		log.Printf("未知的任务类型: %s", genre)
		return false
	}

	// 检查任务是否成功
	if success {
		// 更新任务状态为完成
		tasks[0].Msg = "任务执行完成"
		UploadTimedUploads(tasks[0])

		// 将其余任务标记为撞车
		for i := 1; i < len(tasks); i++ {
			tasks[i].Msg = fmt.Sprintf("撞车了，%s排队中...", genre)
			UploadTimedUploads(tasks[i])
		}
		log.Printf("任务 %s 执行成功 for 店铺: %s (PID: %s)", tasks[0].Name, shopName, pid)
		return true
	} else {
		log.Printf("任务 %s 执行失败 for 店铺: %s (PID: %s)", tasks[0].Name, shopName, pid)
		return false
	}
}

func GetPIDByShopName(shopName string) string {
	var pid string
	err := config.Db.Get(&pid, "SELECT pid FROM shop_mapping WHERE shop_name = ?", shopName)
	if err != nil {
		log.Println("获取 PID 失败 for 店铺:", shopName, "错误信息:", err)
		return ""
	}
	return pid
}

// 价格
func priceRun(tus []mode.TimedUploads) bool {
	if len(tus) == 0 {
		log.Println("价格出现空情况")
		return false
	}
	log.Println("开始上传价格", tus[0].Cron, tus[0].Seller)
	for i := range tus {
		tus[i].Msg = "开始上传价格"
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	success := true
	for i := range tus {
		isCh = true
		for iio := 0; iio < 4; iio++ {
			chHoppingup(tus[i])
			if isCh {
				log.Println(MsgCh)
				break
			}
			log.Println("开始重试：", iio+1)
		}
		if isCh {
			tus[i].Msg = "上传价格完成"
		} else {
			tus[i].Msg = MsgCh
			success = false
		}

		log.Printf("任务 %s 开始更新状态为: %s", tus[i].Name, tus[i].Msg)
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	log.Println("上传价格完成")
	return success
}

// 库存
func inventoryRun(tus []mode.TimedUploads) bool {
	if len(tus) == 0 {
		log.Println("库存出现空情况")
		return false
	}
	log.Println("开始上传库存", tus[0].Cron, tus[0].Seller)
	for i := range tus {
		tus[i].Msg = "开始上传库存"
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	success := true
	for i := range tus {
		isCh = true
		for iio := 0; iio < 4; iio++ {
			chHoppingup(tus[i])
			if isCh {
				log.Println(MsgCh)
				break
			}
			log.Println("开始重试：", iio+1)
		}
		if isCh {
			tus[i].Msg = "上传库存完成"
		} else {
			tus[i].Msg = MsgCh
			success = false
		}

		log.Printf("任务 %s 开始更新状态为: %s", tus[i].Name, tus[i].Msg)
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	log.Println("上传库存完成")
	return success
}

// 类别
func categoryRun(tus []mode.TimedUploads) bool {
	if len(tus) == 0 {
		log.Println("类别出现空情况")
		return false
	}
	log.Println("开始上传类别", tus[0].Cron, tus[0].Seller)
	for i := range tus {
		tus[i].Msg = "开始上传类别"
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	success := true
	for i := range tus {
		isCh = true
		for iio := 0; iio < 4; iio++ {
			chHoppingup(tus[i])
			if isCh {
				log.Println(MsgCh)
				break
			}
			log.Println("开始重试：", iio+1)
		}
		if isCh {
			tus[i].Msg = "上传类别完成"
		} else {
			tus[i].Msg = MsgCh
			success = false
		}

		log.Printf("任务 %s 开始更新状态为: %s", tus[i].Name, tus[i].Msg)
		affected := UploadTimedUploads(tus[i])
		if affected == 0 {
			log.Printf("任务 %s 状态更新失败", tus[i].Name)
		} else {
			log.Printf("任务 %s 状态更新成功, 受影响的行数: %d", tus[i].Name, affected)
		}
	}
	log.Println("上传类别完成")
	return success
}

// 操作浏览器
func chHoppingup(tu mode.TimedUploads) {
	port := se[tu.Seller]
	if !isPortOpen("127.0.0.1", port) {
		log.Printf("端口 %s 未打开，任务 %s 无法执行", port, tu.Name)
		tu.Msg = fmt.Sprintf("端口 %s 未打开，无法执行", port)
		UploadTimedUploads(tu)
		return
	}

	allocator, _ := chromedp.NewRemoteAllocator(context.Background(), fmt.Sprintf("ws://127.0.0.1:%s/devtools/browser", port))

	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	chromedp.Run(ctx,
		startup(tu),
	)
}

// 控制器
func startup(tu mode.TimedUploads) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if loginup(ctx, tu) {
			MsgCh, isCh = uploadup(ctx, tu)
		} else {
			MsgCh, isCh = "登录失败", false
		}
		return err
	}
}

// 检测是否登录，未登录就登录
func loginup(ctx context.Context, tu mode.TimedUploads) bool {
	timeout, cancel := context.WithTimeout(ctx, 80*time.Second)
	defer cancel()
	chromedp.Navigate("https://seller.walmart.com/catalog/list-items").Do(timeout)
	for i := 0; i < 3; i++ {
		timeout0, cancel0 := context.WithTimeout(ctx, 30*time.Second)
		defer cancel0()

		err := chromedp.WaitVisible(`button[data-automation-id="itemListPageHeaderManageItemsButtonNode"]`).Do(timeout0)
		if err != nil {
			timeout02, cancel02 := context.WithTimeout(ctx, 30*time.Second)
			defer cancel02()
			err := chromedp.WaitVisible(`input[data-automation-id="uname"]`).Do(timeout02)
			if err != nil {
				log.Println("页面加载失败，重新开始加载")
				timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
				defer cancel01()
				chromedp.Stop().Do(timeout01)
				chromedp.Sleep(time.Second * 1).Do(timeout)
				chromedp.Navigate("https://seller.walmart.com/catalog/list-items").Do(timeout01)
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
	chromedp.SendKeys(`input[data-automation-id="uname"]`, users[tu.Seller][0]).Do(timeout)
	chromedp.Sleep(time.Second * 2).Do(timeout)
	chromedp.SendKeys(`input[data-automation-id="pwd"]`, users[tu.Seller][1]+kb.Enter).Do(timeout)
	chromedp.Sleep(time.Second * 2).Do(timeout)
	timeout01, cancel01 := context.WithTimeout(ctx, 2*time.Second)
	defer cancel01()
	chromedp.SendKeys(`input[data-automation-id="pwd"]`, kb.Enter).Do(timeout01)
	chromedp.Sleep(time.Second * 30).Do(timeout)
	chromedp.Navigate("https://seller.walmart.com/catalog/list-items").Do(timeout)
	chromedp.Sleep(time.Second * 10).Do(timeout)
	timeout03, cancel03 := context.WithTimeout(ctx, 30*time.Second)
	defer cancel03()
	err := chromedp.WaitVisible(`button[data-automation-id="itemListPageHeaderManageItemsButtonNode"]`).Do(timeout03)
	if err != nil {
		log.Println("登录失败")
		chromedp.Stop().Do(ctx)
		return false
	}
	log.Println("登录成功")
	return true

}

// 上传函数

func uploadup(ctx context.Context, tu mode.TimedUploads) (string, bool) {
	var msg string
	var iso bool

	for i := 0; i < 2; i++ {
		if i != 0 {
			log.Println("开始重试：", i)
			timeout00, cancel00 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel00()
			err := chromedp.Run(timeout00,
				chromedp.Navigate("https://seller.walmart.com/catalog/list-items"),
			)
			if err != nil {
				log.Printf("重试导航失败: %v", err)
				continue
			}
		}

		timeout, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		for j := 0; j < 2; j++ {
			timeout0, cancel0 := context.WithTimeout(ctx, 20*time.Second)
			defer cancel0()
			err := chromedp.Run(timeout0,
				chromedp.WaitVisible(`button[data-automation-id="itemListPageHeaderManageItemsButtonNode"]`),
			)
			if err != nil {
				log.Println("重新加载页面")
				timeout01, cancel01 := context.WithTimeout(ctx, 20*time.Second)
				defer cancel01()
				err := chromedp.Run(timeout01,
					chromedp.Navigate("https://seller.walmart.com/catalog/list-items"),
				)
				if err != nil {
					log.Printf("导航失败: %v", err)
					continue
				}
			} else {
				break
			}
		}

		log.Println("开始上传")
		// 隐藏横幅和右下角弹窗
		//err := chromedp.Run(timeout,
		//	chromedp.EvaluateAsDevTools(`document.getElementsByClassName("h-48")[0].style.display = "none"`, &tmp),
		//	chromedp.EvaluateAsDevTools(`
		//		var elements = document.querySelectorAll('[data-vertical-alignment="Bottom Right Aligned"]');
		//		for (var k = 0; k < elements.length; k++) {
		//			elements[k].style.display = 'none';
		//		}`, &tmp),
		//)
		//if err != nil {
		//	log.Printf("隐藏横幅和弹窗失败: %v", err)
		//	continue
		//}
		switch tu.Genre {
		case "价格":

		case "库存":
			// 开始 chromedp 操作
			err := chromedp.Run(timeout,
				chromedp.Sleep(time.Second*2),
				// 点击"Manage Items"按钮
				chromedp.Click(`button[data-automation-id="itemListPageHeaderManageItemsButtonNode"]`),
				chromedp.Sleep(time.Second*2),
				// 点击第一个 td 元素
				chromedp.Click(`(//td[@class="Options-module_cell__nel-r flex-grow"])[1]`),
				chromedp.Sleep(time.Second*2),
				// 等待第一个下拉框出现并点击
				chromedp.WaitVisible(`//select[contains(@id, 'ld_select_')][1]`),
				chromedp.Click(`//select[contains(@id, 'ld_select_')][1]`),
				chromedp.Sleep(time.Second*2),
			)
			if err != nil {
				log.Printf("执行步骤失败: %v", err)
				msg = "选择库存模板失败"
				iso = false
				return msg, iso
			}
			// 获取所有 option 的 value 值

			var optionValues []string
			err = chromedp.Run(ctx,
				chromedp.Evaluate(`(function() {
            var options = document.evaluate("//select[contains(@id, 'ld_select_')][1]/option", document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
            var values = [];
            for (var i = 0; i < options.snapshotLength; i++) {
                values.push(options.snapshotItem(i).value);
            }
            return values;
        })()`, &optionValues),
			)
			if err != nil {
				log.Fatalf("获取 option 的 value 失败: %v", err)
			}

			// 打印所有的 option value 值
			for i, val := range optionValues {
				fmt.Printf("Option %d: %s\n", i+1, val)
			}

			// 查找目标值并选择
			for i, val := range optionValues {
				if val == "MP_INVENTORY" {
					log.Println(i)
					// 模拟键盘操作选择到指定的 option
					for j := 0; j < i; j++ {
						err = chromedp.Run(ctx,
							// 模拟按下箭头键 (ArrowDown) 选择
							chromedp.SendKeys(`//select[contains(@id, 'ld_select_')][1]`, kb.ArrowDown),
							chromedp.Sleep(500*time.Millisecond),
						)
						if err != nil {
							log.Printf("模拟箭头键失败: %v", err)
						}
					}

					// 最后按下回车键来选中
					err = chromedp.Run(ctx,
						chromedp.Sleep(500*time.Millisecond),
						chromedp.SendKeys(`//select[contains(@id, 'ld_select_')][1]`, kb.Enter),
					)
					if err != nil {
						log.Fatalf("选择 option 失败: %v", err)
					}

					fmt.Println("成功选择 MP_INVENTORY")
					break
				}
			}

			// 取消文件输入框的 hidden 属性，设置 display 样式，并移除可能的隐藏类
			err = chromedp.Run(timeout,
				chromedp.Evaluate(`(function() {
        var elem = document.getElementById("bulkUpdateFileUploaderInputNode");
        if (elem) {
            elem.removeAttribute("hidden");
            elem.style.display = "block";
            elem.style.visibility = "visible"; // 确保元素可见
            elem.classList.remove("hidden"); // 移除可能的隐藏类
        }
    })()`, nil),
			)
			if err != nil {
				log.Printf("显示文件输入框失败: %v", err)
				msg = "显示文件输入框失败"
				iso = false
				return msg, iso
			}

			// 通过 SendKeys 上传文件
			err = chromedp.Run(timeout,
				chromedp.SendKeys(`#bulkUpdateFileUploaderInputNode`, fmt.Sprintf(`C:\Users\Administrator\Desktop\file\%s`, tu.File)),
			)
			if err != nil {
				log.Printf("发送文件路径失败: %v", err)
				msg = "发送文件路径失败"
				iso = false
				return msg, iso
			}

			// 轮询页面查找成功提示
			var success bool
			startTime := time.Now()
			for {
				// 每次轮询检查页面文本
				err = chromedp.Run(ctx,
					chromedp.Evaluate(`document.body.innerText.includes("File submitted successfully.")`, &success),
				)
				if err != nil {
					log.Printf("检查上传结果失败: %v", err)
				}

				// 如果找到文本，表示上传成功
				if success {
					fmt.Println("文件上传成功")
					msg = "文件上传成功"
					iso = true
					return msg, iso

				}

				// 如果超过最大等待时间 10 秒，停止轮询
				if time.Since(startTime) > 10*time.Second {
					fmt.Println("文件上传成功的提示未检测到，可能上传失败或超时。")
					break
				}

				// 短暂的延迟，避免 CPU 过度使用
				time.Sleep(500 * time.Millisecond)
			}

		case "类别":
			// 开始 chromedp 操作
			err := chromedp.Run(timeout,
				chromedp.Sleep(time.Second*2),
				// 点击"Manage Items"按钮
				chromedp.Click(`button[data-automation-id="itemListPageHeaderManageItemsButtonNode"]`),
				chromedp.Sleep(time.Second*2),
				// 点击第一个 td 元素
				chromedp.Click(`(//td[@class="Options-module_cell__nel-r flex-grow"])[1]`),
				chromedp.Sleep(time.Second*2),
				// 等待第一个下拉框出现并点击
				chromedp.WaitVisible(`//select[contains(@id, 'ld_select_')][1]`),
				chromedp.Click(`//select[contains(@id, 'ld_select_')][1]`),
				chromedp.Sleep(time.Second*2),
			)
			if err != nil {
				log.Printf("执行步骤失败: %v", err)
				msg = "选择库存模板失败"
				iso = false
				return msg, iso
			}
			// 获取所有 option 的 value 值

			var optionValues []string
			err = chromedp.Run(ctx,
				chromedp.Evaluate(`(function() {
            var options = document.evaluate("//select[contains(@id, 'ld_select_')][1]/option", document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
            var values = [];
            for (var i = 0; i < options.snapshotLength; i++) {
                values.push(options.snapshotItem(i).value);
            }
            return values;
        })()`, &optionValues),
			)
			if err != nil {
				log.Fatalf("获取 option 的 value 失败: %v", err)
			}

			// 打印所有的 option value 值
			for i, val := range optionValues {
				fmt.Printf("Option %d: %s\n", i+1, val)
			}

			// 查找目标值并选择
			for i, val := range optionValues {
				if val == "UPDATE_ALL_ATTRIBUTES" {
					log.Println(i)
					// 模拟键盘操作选择到指定的 option
					for j := 0; j < i; j++ {
						err = chromedp.Run(ctx,
							// 模拟按下箭头键 (ArrowDown) 选择
							chromedp.SendKeys(`//select[contains(@id, 'ld_select_')][1]`, kb.ArrowDown),
							chromedp.Sleep(500*time.Millisecond),
						)
						if err != nil {
							log.Printf("模拟箭头键失败: %v", err)
						}
					}

					// 最后按下回车键来选中
					err = chromedp.Run(ctx,
						//chromedp.Sleep(500*time.Millisecond),
						chromedp.SendKeys(`//select[contains(@id, 'ld_select_')][1]`, kb.Enter),
						chromedp.Sleep(500*time.Millisecond),
					)
					if err != nil {
						log.Fatalf("选择 option 失败: %v", err)
					}

					fmt.Println("成功选择 UPDATE_ALL_ATTRIBUTES")
					break
				}
			}
			fmt.Println("第一次获取的选项值:", optionValues)
			// 获取第二个select的所有 option 的 value 值

			var optionValues2 []string
			err = chromedp.Run(ctx,
				chromedp.WaitVisible(`(//select[contains(@id, 'ld_select_')])[2]`),
				chromedp.Click(`(//select[contains(@id, 'ld_select_')])[2]`),
				chromedp.Evaluate(`(function() {
            var options = document.evaluate("(//select[contains(@id, 'ld_select_')])[2]/option", document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
            var values = [];
            for (var i = 0; i < options.snapshotLength; i++) {
                values.push(options.snapshotItem(i).value);
            }
            return values;
        })()`, &optionValues2),
			)
			if err != nil {
				log.Fatalf("获取 option 的 value 失败: %v", err)
			}

			// 打印所有的 option value 值
			for i, val := range optionValues2 {
				fmt.Printf("Option %d: %s\n", i+1, val)
			}

			// 查找目标值并选择
			for i, val := range optionValues2 {
				if val == "sfs" {
					log.Println(i)
					// 模拟键盘操作选择到指定的 option
					for j := 0; j < i; j++ {
						err = chromedp.Run(ctx,
							// 模拟按下箭头键 (ArrowDown) 选择
							chromedp.SendKeys(`(//select[contains(@id, 'ld_select_')])[2]`, kb.ArrowDown),
							chromedp.Sleep(500*time.Millisecond),
						)
						if err != nil {
							log.Printf("模拟箭头键失败: %v", err)
						}
					}

					// 最后按下回车键来选中
					err = chromedp.Run(ctx,
						//chromedp.Sleep(500*time.Millisecond),
						chromedp.SendKeys(`(//select[contains(@id, 'ld_select_')])[2]`, kb.Enter),
						chromedp.Sleep(500*time.Millisecond),
					)
					if err != nil {
						log.Fatalf("选择 option 失败: %v", err)
					}

					fmt.Println("成功选择 sfs(Seller fulfilled)")
					break
				}
			}

			// 取消文件输入框的 hidden 属性，设置 display 样式，并移除可能的隐藏类
			err = chromedp.Run(timeout,
				chromedp.Evaluate(`(function() {
        var elem = document.getElementById("bulkUpdateFileUploaderInputNode");
        if (elem) {
            elem.removeAttribute("hidden");
            elem.style.display = "block";
            elem.style.visibility = "visible"; // 确保元素可见
            elem.classList.remove("hidden"); // 移除可能的隐藏类
        }
    })()`, nil),
			)
			if err != nil {
				log.Printf("显示文件输入框失败: %v", err)
				msg = "显示文件输入框失败"
				iso = false
				return msg, iso
			}

			// 通过 SendKeys 上传文件
			err = chromedp.Run(timeout,
				chromedp.SendKeys(`#bulkUpdateFileUploaderInputNode`, fmt.Sprintf(`C:\Users\Administrator\Desktop\file\%s`, tu.File)),
			)
			if err != nil {
				log.Printf("发送文件路径失败: %v", err)
				msg = "发送文件路径失败"
				iso = false
				return msg, iso
			}

			// 轮询页面查找成功提示
			var success bool
			startTime := time.Now()
			for {
				// 每次轮询检查页面文本
				err = chromedp.Run(ctx,
					chromedp.Evaluate(`document.body.innerText.includes("File submitted successfully.")`, &success),
				)
				if err != nil {
					log.Printf("检查上传结果失败: %v", err)
				}

				// 如果找到文本，表示上传成功
				if success {
					fmt.Println("文件上传成功")
					msg = "文件上传成功"
					iso = true
					return msg, iso

				}

				// 如果超过最大等待时间 10 秒，停止轮询
				if time.Since(startTime) > 10*time.Second {
					fmt.Println("文件上传成功的提示未检测到，可能上传失败或超时。")
					msg = "文件上传提示失败。"
					iso = false
					return msg, iso
				}

				// 短暂的延迟，避免 CPU 过度使用
				time.Sleep(500 * time.Millisecond)
			}
		}
		chromedp.Sleep(time.Second * 2).Do(ctx)
		log.Println("提交过程完成")
		// 关闭当前页面
		//log.Println("关闭当前页面")
		//chromedp.Sleep(time.Second * 2).Do(ctx)
		//cancel()
		return "完成", true

	}
	return msg, iso
}
