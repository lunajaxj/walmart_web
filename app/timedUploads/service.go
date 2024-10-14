package timedUploads

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddTimedUploads(tu mode.TimedUploads) {
	_, err := config.Db.NamedExec(`INSERT INTO timed_uploads (name,genre,seller,cron,file,msg)
       									VALUES (:name,:genre,:seller,:cron,:file,:msg)`, tu)
	tools.ErrPr(err, "")
}

func GetTimedUploadsWithShopName(page, limit int, genre, name, msg, file, seller string) ([]mode.TimedUploads, int) {
	var tus []mode.TimedUploads
	var count int
	var wheres string

	if genre != "" {
		wheres = tools.WhereAnd(wheres, "=", "tu.genre", genre)
	}
	if file != "" {
		wheres = tools.WhereAnd(wheres, "=", "tu.file", file)
	}
	if name != "" {
		wheres = tools.WhereAnd(wheres, "LIKE", "tu.name", "%"+name+"%")
	}
	if seller != "" {
		wheres = tools.WhereAnd(wheres, "=", "tu.seller", seller)
	}
	if msg != "" {
		if msg == " " {
			msg = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "tu.msg", msg)
	}

	page = (page - 1) * limit

	// SQL 查询，联接 shop_mapping 表获取 shop_name
	query := fmt.Sprintf(`
    SELECT tu.*, DATE_FORMAT(tu.update_time, '%%Y-%%m-%%d %%H:%%i:%%s') AS update_time, sm.shop_name 
    FROM timed_uploads tu
    LEFT JOIN shop_mapping sm ON tu.seller = sm.pid
    %s
    LIMIT %d, %d
`, wheres, page, limit)

	// 打印查询结果
	err := config.Db.Select(&tus, query)
	//log.Printf("查询结果: %+v", tus)
	if err != nil {
		// 记录 SQL 查询错误及相关参数
		log.Printf("SQL 查询错误: %v | 查询: %s | 参数: page=%d, limit=%d, genre=%s, name=%s, file=%s, seller=%s", err, query, page, limit, genre, name, file, seller)
		return nil, 0
	}
	// 查询总数并打印
	countQuery := fmt.Sprintf("SELECT COUNT(tu.tu_id) FROM timed_uploads tu LEFT JOIN shop_mapping sm ON tu.seller = sm.pid %s", wheres)
	//log.Printf("执行计数查询: %s", countQuery)

	err = config.Db.Get(&count, countQuery)
	if err != nil {
		log.Printf("计数查询错误: %v", err)
	}

	// 打印总数
	//log.Printf("查询到的总记录数: %d", count)
	return tus, count
}

func GetTimedUploadsCron(genre, pid string) []string {
	var list []string
	err := config.Db.Select(&list, fmt.Sprintf("SELECT distinct cron FROM timed_uploads where genre ='%s' and seller = '%s'", genre, pid))
	if err != nil {
		log.Printf("查询 cron 表达式失败: %v", err)
	}
	//log.Println("list", list)
	return list
}
func parseUpdateTime(task mode.TimedUploads) (*time.Time, error) {
	// 手动解析时间格式
	layout := "2006-01-02 15:04:05" // MySQL datetime 的格式
	parsedTime, err := time.Parse(layout, task.UpdateTime)
	if err != nil {
		return nil, fmt.Errorf("无法解析时间: %v", err)
	}
	return &parsedTime, nil
}
func GetTimedUploadsByCron(genre, cron, pid string) []mode.TimedUploads {
	var tasks []mode.TimedUploads
	// 使用 pid 和 cron 表达式获取任务
	//query := "SELECT * FROM timed_uploads WHERE genre = ? AND seller = ? AND cron = ?"
	query := "SELECT t.tu_id,t.name,t.genre,t.file,t.seller,s.shop_name,t.msg,t.cron,DATE_FORMAT(t.update_time,'%Y-%m-%d %h:%m:%s') AS 'update_time' FROM timed_uploads t,shop_mapping s WHERE t.seller = s.pid AND t.genre = ? AND t.seller = ? AND t.cron = ?"
	err := config.Db.Select(&tasks, query, genre, pid, cron)
	if err != nil {
		log.Printf("查询任务失败: %v | genre: %s | pid: %s | cron: %s", err, genre, pid, cron)
		return tasks
	}

	for i := range tasks {
		parsedTime, err := parseUpdateTime(tasks[i])
		if err != nil {
			log.Printf("时间解析错误: %v", err)
		} else {
			// 如果解析成功，更新 UpdateTime 字段
			tasks[i].UpdateTime = parsedTime.Format("2006-01-02 15:04:05")
		}
	}

	log.Printf("查询任务成功: %v", tasks)
	return tasks
}

func UploadTimedUploads(tu mode.TimedUploads) int {
	stmt := "UPDATE timed_uploads set name=:name, genre=:genre, cron=:cron, file=:file, msg=:msg, seller=:seller WHERE tu_id=:tu_id"
	affected, err := config.Db.NamedExec(stmt, map[string]interface{}{
		"name":   tu.Name,
		"genre":  tu.Genre,
		"cron":   tu.Cron,
		"file":   tu.File,
		"msg":    tu.Msg,
		"seller": tu.Seller, // 此处的 seller 是 pid
		"tu_id":  tu.TuId,
	})

	if err != nil {
		log.Printf("任务 %s 状态更新失败: %v", tu.Name, err)
		return 0 // 返回 0 表示更新失败
	}

	rowsAffected, err := affected.RowsAffected()
	if err != nil {
		log.Printf("获取受影响行数失败: %v", err)
		return 0
	}

	return int(rowsAffected)
}

func Remove(ids string) int {
	split := strings.Split(ids, ",")
	var wheres string
	if len(split) == 0 {
		return 0
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "tu_id", split[i])
	}
	updateResult := config.Db.MustExec(fmt.Sprintf("DELETE FROM timed_uploads %s", wheres))
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func DelMsg() {
	config.Db.MustExec("UPDATE timed_uploads SET msg = ''")
}
func GetMsg() []string {
	var list []string
	err := config.Db.Select(&list, "SELECT distinct msg FROM timed_uploads")
	tools.ErrPr(err, "")
	return list
}

type Seller struct {
	PID      string `db:"pid"`
	ShopName string `db:"shop_name"`
}

func GetSeller() []Seller {
	var sellers []Seller
	query := `
		SELECT DISTINCT sm.pid, sm.shop_name
		FROM timed_uploads tu
		JOIN shop_mapping sm ON tu.seller = sm.pid
	`
	err := config.Db.Select(&sellers, query)
	if err != nil {
		log.Printf("SQL 查询错误: %v", err)
	}
	return sellers
}

func GetFile2() (list []map[string]string) {
	directory := "file" // 指定目录的路径

	// 获取目录下的所有文件和目录
	files, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	// 遍历所有文件和目录
	for _, file := range files {
		// 判断是否为文件
		if !file.IsDir() {
			info, _ := file.Info()
			modTimeString := info.ModTime().Format("2006-01-02 15:04:05") // 格式化为 "年-月-日 时:分:秒"

			// 获取文件名并去除后缀
			list = append(list, map[string]string{"name": file.Name(), "time": modTimeString})
		}
	}
	return list
}

func GetFile() (list []string) {
	directory := "file" // 指定目录的路径

	// 获取目录下的所有文件和目录
	files, err := os.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	// 遍历所有文件和目录
	for _, file := range files {
		// 判断是否为文件
		if !file.IsDir() {
			// 获取文件名
			list = append(list, file.Name())
		}
	}
	return list
}

// ShopMapping 表示 PID 和店铺名称的结构体
type ShopMapping struct {
	PID      string `db:"pid"`
	ShopName string `db:"shop_name"`
}

// GetShopMapping 获取 PID 和店铺名称的映射
func GetShopMapping() []ShopMapping {
	var mappings []ShopMapping
	err := config.Db.Select(&mappings, "SELECT pid, shop_name FROM shop_mapping")
	if err != nil {
		log.Println("Failed to retrieve shop mappings:", err)
		return nil
	}
	return mappings
}
