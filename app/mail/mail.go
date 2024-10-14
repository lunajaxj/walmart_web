package mail

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"golang.org/x/net/context"
	"log"
	"strings"
	"sync"
	"time"
	"walmart_web/app/mode"
)

//var users = map[string][]string{"demo": {"xcitnpkj@126.com", "*IcI4$fz152NFjI"}}
//var se = map[string]string{"demo": "8444"}

var users = map[string][]string{"Online mini-mart": {"yumumandzsw@163.com", "Dyumuman123!"}, "GUBIN": {"guibinwalmart@163.com", "DGuibin123!"}, "Money Saving Center": {"ketianwalmart@163.com", "DKetian123!"}, "Money Saving World": {"ziyuewlkj@163.com", "Dziyue123!"}, "Kiaote Center": {"qiaotewalmart@163.com", "Dqiaote123!"}, "Moremuma": {"moremuma@163.com", "$Moremuma778"}, "Juno Kael": {"keyshinewalmart@163.com", "&Keyshine788"}}
var se = map[string]string{"GUBIN": "8333", "Money Saving Center": "9333", "Money Saving World": "7333", "Kiaote Center": "10222", "Moremuma": "11222", "Juno Kael": "12222"}

var IsRun = map[string]bool{"Online mini-mart": false, "GUBIN": false, "Money Saving Center": false, "Money Saving World": false, "Kiaote Center": false, "Moremuma": false, "Juno Kael": false}
var Mutexs = map[string]*sync.Mutex{"Online mini-mart": new(sync.Mutex), "GUBIN": new(sync.Mutex), "Money Saving Center": new(sync.Mutex), "Money Saving World": new(sync.Mutex), "Kiaote Center": new(sync.Mutex), "Moremuma": new(sync.Mutex), "Juno Kael": new(sync.Mutex)}

// 操作浏览器
func Run(ma mode.Mail) {
	//return
	Mutexs[ma.Seller].Lock()
	IsRun[ma.Seller] = true
	defer func() {
		IsRun[ma.Seller] = false
		Mutexs[ma.Seller].Unlock()
	}()
	// 创建一个分配器来连接到已经打开的浏览器
	allocator, _ := chromedp.NewRemoteAllocator(context.Background(), fmt.Sprintf("ws://127.0.0.1:%s/devtools/browser", se[ma.Seller]))

	// 创建一个新的浏览器上下文
	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()
	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	chromedp.Run(ctx,
		startup(ma),
	)

}

// 控制器
func startup(ma mode.Mail) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		ma.Count = 1
		ma.Msg = "开始执行中"
		UploadMail(ma)

		if loginup(ctx, ma) {
			ma.Msg, _ = showMail(ctx, ma)
		} else {
			ma.Msg, _ = "登录失败", false
		}
		ma.Count = 0
		UploadMail(ma)

		return err
	}
}

// 检测是否登录，未登录就登录
func loginup(ctx context.Context, tu mode.Mail) bool {
	timeout, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout)
	for i := 0; i < 2; i++ {
		timeout0, cancel0 := context.WithTimeout(ctx, 30*time.Second)
		defer cancel0()

		err := chromedp.WaitVisible(`input[id="search"]`).Do(timeout0)
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
				chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
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
	chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout)
	chromedp.Sleep(time.Second * 10).Do(timeout)
	timeout03, cancel03 := context.WithTimeout(ctx, 30*time.Second)
	defer cancel03()
	err := chromedp.WaitVisible(`input[id="search"]`).Do(timeout03)
	if err != nil {
		log.Println("登录失败")
		chromedp.Stop().Do(ctx)
		return false
	}
	log.Println("登录成功")
	return true

}

// 查看邮件
func showMail(ctx context.Context, ma mode.Mail) (string, bool) {
	var msg string
	var iso bool
	//for i := 0; i < 3; i++ {
	//	if i != 0 {
	//		log.Println("开始重试：", i)
	//		timeout00, cance00 := context.WithTimeout(ctx, 10*time.Second)
	//		defer cance00()
	//		chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout00)
	//	}

	timeout, cancel := context.WithTimeout(ctx, 6*time.Minute)
	defer cancel()
	for i := 0; i < 2; i++ {
		timeout0, cancel0 := context.WithTimeout(ctx, 30*time.Second)
		defer cancel0()
		timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
		defer cancel01()
		err := chromedp.WaitVisible(`input[id="search"]`).Do(timeout0)
		if err != nil {
			log.Println("重新加载页面")
			chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
		} else {
			break
		}
	}
	log.Println("开始点击未读邮件")
	//橫幅
	var i string
	chromedp.EvaluateAsDevTools(`document.getElementsByClassName("h-48")[0].style.display = "none"`, &i).Do(timeout)
	//右下角弹窗
	chromedp.EvaluateAsDevTools(`
			var elements = document.querySelectorAll('[data-vertical-alignment="Bottom Right Aligned"]');
			for (var i = 0; i < elements.length; i++) {
				elements[i].style.display = 'none';
		}`, &i).Do(timeout)
	for ii := 0; ii < 6; {
		chromedp.Click(`div.dM4RB  div.css-1hwfws3`).Do(timeout)
		chromedp.Sleep(time.Second * 2).Do(timeout)
		chromedp.Click(`[id^="react-select-"][id$="-option-1"]`).Do(timeout)
		chromedp.Sleep(time.Second * 2).Do(timeout)
		var source string
		if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools("document.documentElement.outerHTML", &source)); err != nil {
			log.Println("获取源码失败")
			msg = "获取源码失败"
			iso = false
			ii++
			timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel01()
			chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
			continue
		}
		if !strings.Contains(source, "Unread</div>") {
			msg = "选择出现错误"
			iso = false
			ii++
			timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel01()
			chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
			continue
		}
		timeout00, cancel00 := context.WithTimeout(ctx, 3*time.Second)
		defer cancel00()
		err := chromedp.Click(`div._12Qdh:first-child`).Do(timeout00)
		if err != nil {
			log.Println("已无未读邮件,任务结束")
			return "无未读邮件", true
		}
		chromedp.Sleep(time.Second * 2).Do(timeout)
		if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools("document.documentElement.outerHTML", &source)); err != nil {
			log.Println("获取源码失败")
			msg = "获取源码失败"
			iso = false
			ii++
			timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel01()
			chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
			continue
		}
		if strings.Contains(source, "Select and view messages from customers") {
			log.Println("点击邮件失败")
			msg = "点击邮件失败"
			iso = false
			ii++
			timeout01, cancel01 := context.WithTimeout(ctx, 10*time.Second)
			defer cancel01()
			chromedp.Navigate("https://seller.walmart.com/seller-communication/Customer").Do(timeout01)
			continue
		}

		ma.Msg = fmt.Sprint("已点击", ma.Count, "次")
		UploadMail(ma)
		ma.Count += 1
	}
	return msg, iso
	//}
	//return msg, iso
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
