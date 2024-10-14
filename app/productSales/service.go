package productSales

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func UpdateSalesDay(sales mode.ProductSalesDay) {
	// 获取当前时间
	now := time.Now()

	// 将时间戳截断为天
	truncatedTime := now.Truncate(24 * time.Hour)

	// 输出精确到天的时间
	day := truncatedTime.Format("2006-01-02")
	sales.CreateDate = day

	_, err := config.Db.NamedExec(`
		INSERT INTO product_sales_day (
			item_id,
			catalog_seller_id,
			day0,
			create_date
		) VALUES (
			:item_id,
			:catalog_seller_id,
			:day0,
			:create_date
		)
	ON DUPLICATE KEY UPDATE day0=:day0
	`, sales)
	tools.ErrPr(err, "")
	if err != nil {
		log.Println("库存1插入数据失败:", err)
	}

}
func UpdateSales(sales mode.ProductSales) {
	_, err := config.Db.NamedExec(`
		UPDATE product_sales SET img=:img,day_01=:day_01,day_02=:day_02,day_03=:day_03,day_04=:day_04,day_05=:day_05,day_06=:day_06,day_07=:day_07,day_15=:day_15,day_30=:day_30,day_60=:day_60,day_90=:day_90 WHERE item_id=:item_id
	`, sales)
	tools.ErrPr(err, "")
	if err != nil {
		log.Println("库存2插入数据失败:", err)
	}
}
func SelectSales(page, limit int, id string) ([]mode.ProductSales, int) {
	var prs []mode.ProductSales
	var counts []int
	var wheres string
	if len(id) != 0 {
		split := strings.Split(id, "\n")
		for i := range split {
			wheres = tools.WhereOrInt(wheres, "=", "item_id", split[i])
		}
	}
	page = (page - 1) * limit
	sql := fmt.Sprintf("SELECT * FROM product_sales %s LIMIT %d,%d", wheres, page, limit)
	err := config.Db.Select(&prs, sql)
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(item_id) FROM product_sales %s", wheres))
	tools.ErrPr(err, "")

	return prs, counts[0]
}

func RemoveSales(id string) []mode.ProductSales {
	split := strings.Split(id, ",")
	var wheres string
	var prs []mode.ProductSales
	if len(split) == 0 {
		return nil
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "item_id", split[i])
	}
	sql := fmt.Sprintf("delete from product_sales %s", wheres)
	err := config.Db.Select(&prs, sql)
	tools.ErrPr(err, "")
	return prs
}

func AddSales(sales []mode.ProductSales) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Fatal(err)
	}

	for i := range sales {
		_, err = tx.NamedExec(`
		INSERT INTO product_sales (
			item_id,
			catalog_seller_id,
		    img, 
		    day_01,
		    day_02,
		    day_03,
		    day_04,
		    day_05,
		    day_06,
		    day_07,
		    day_15,
		    day_30,
		    day_60,
		    day_90,
			create_date
		) VALUES (
			:item_id,
		    :catalog_seller_id,
		    :img, 
		    :day_01,
		    :day_02,
		    :day_03,
		    :day_04,
		    :day_05,
		    :day_06,
		    :day_07,
		    :day_15,
		    :day_30,
		    :day_60,
		    :day_90,
			:create_date
		)
	ON DUPLICATE KEY UPDATE catalog_seller_id=:catalog_seller_id
	`, sales[i])
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true

}

func SelectSalesDayByItemId(id string) []mode.ProductSalesDay {
	var prs []mode.ProductSalesDay
	var wheres string
	if len(id) != 0 {
		split := strings.Split(id, "\n")
		for i := range split {
			wheres = tools.WhereOrInt(wheres, "=", "item_id", split[i])
		}
	}
	sql := fmt.Sprintf("SELECT * FROM product_sales_day %s ORDER BY create_date DESC LIMIT 200;", wheres)
	err := config.Db.Select(&prs, sql)
	tools.ErrPr(err, "")
	return prs
}
func Update(sales mode.ProductSalesDay, img string) {
	UpdateSalesDay(sales)
	prs := SelectSalesDayByItemId(strconv.Itoa(sales.ITEMID))
	var xl int
	var seae mode.ProductSales
	seae.ITEMID = sales.ITEMID
	seae.Img = img
j:
	for i := 0; i < len(prs); i++ {
		var res string
		if i == len(prs)-1 && (prs[i].Day0 != "有跟卖" || prs[i].Day0 != "自发货") {
			//第一天没有以往数据，因此将销售额设为缺失数据
			res = "缺少计算数据"
		} else {
			prevQuantity := prs[i+1].Day0
			currQuantity := prs[i].Day0

			if currQuantity == "有跟卖" || currQuantity == "自发货" {
				// 库存数量无效，将销售标记为错误
				res = currQuantity
			} else if prevQuantity == "有跟卖" || prevQuantity == "自发货" {
				// 如果前一天的数量无效，则搜索最后有效的数量
				j := i + 1
				for j < len(prs) && (prs[j].Day0 == "有跟卖" || prs[j].Day0 == "自发货") {
					j++
				}
				if j == len(prs) {
					// 未找到有效的先前数据，将销售额标记为缺失数据
					res = "缺少计算数据"
				} else {
					prevQuantity = prs[j].Day0
					atoi1, err := strconv.Atoi(prevQuantity)
					if err != nil {
						log.Println(err)
					}
					atoi2, err := strconv.Atoi(currQuantity)
					if err != nil {
						log.Println(err)
					}
					z := atoi1 - atoi2
					res = strconv.Itoa(z)
					if z < 0 {
						z = 0
						res = "0"
					}
					xl += z

				}
			} else {
				atoi1, err := strconv.Atoi(prevQuantity)
				if err != nil {
					log.Println(err)
				}
				atoi2, err := strconv.Atoi(currQuantity)
				if err != nil {
					log.Println(err)
				}
				z := atoi1 - atoi2
				res = strconv.Itoa(z)
				if z < 0 {
					z = 0
					res = "0"
				}
				xl += z
			}
		}

		switch i + 1 {
		case 1:
			seae.Day01 = res
		case 2:
			seae.Day02 = res
		case 3:
			seae.Day03 = res
		case 4:
			seae.Day04 = res
		case 5:
			seae.Day05 = res
		case 6:
			seae.Day06 = res
		case 7:
			seae.Day07 = strconv.Itoa(xl)
		case 15:
			seae.Day15 = strconv.Itoa(xl)
		case 30:
			seae.Day30 = strconv.Itoa(xl)
		case 60:
			seae.Day60 = strconv.Itoa(xl)
		case 90:
			seae.Day90 = strconv.Itoa(xl)
			break j
		}

	}
	UpdateSales(seae)
}
