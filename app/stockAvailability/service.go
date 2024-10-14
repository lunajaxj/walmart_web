package stockAvailability

import (
	"fmt"
	"log"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func GetStockAvailability(page, limit int, id string) ([]mode.StockAvailability, int) {
	var sas []mode.StockAvailability
	var counts []int
	var wheres string
	if id != "" {
		split := strings.Split(id, "\n")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("item_id", len(split)),
			split)
	}
	page = (page - 1) * limit
	err := config.Db.Select(&sas, fmt.Sprintf("SELECT * FROM stock_availability %s LIMIT %d,%d", wheres, page, limit))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(item_id) FROM stock_availability %s", wheres))
	tools.ErrPr(err, "")
	return sas, counts[0]
}
func Remove(ids string) int {
	split := strings.Split(ids, ",")
	var wheres string
	if len(split) == 0 {
		return 0
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "item_id", split[i])
	}
	updateResult := config.Db.MustExec(fmt.Sprintf("DELETE FROM stock_availability %s", wheres))
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func Install(sas []mode.StockAvailability) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Fatal(err)
	}

	// 准备插入语句
	stmt, err := tx.Preparex("INSERT INTO stock_availability (item_id,sales_user,img,cy_sku ,cy_name ,gtin,pt_sku ,declaration,num,warehouse ,lead_time , counts,transit_num,library_num,weighted,remarks1,remarks2) VALUES (?, ?, ?, ?, ?, ?,? , ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// 执行批量插入
	for _, sa := range sas {
		_, err := stmt.Exec(sa.ItemId, sa.SalesUser, sa.Img, sa.CySku, sa.CyName, sa.Gtin, sa.PtSku, sa.Declaration, sa.Num, sa.Warehouse, sa.LeadTime, sa.Counts, sa.TransitNum, sa.LibraryNum, sa.Weighted, sa.Remarks1, sa.Remarks2)
		if err != nil {
			tx.Rollback()
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
func GetStockAvailabilityAll() ([]mode.StockAvailability, int) {
	var sas []mode.StockAvailability
	var counts []int

	err := config.Db.Select(&sas, fmt.Sprint("SELECT * FROM stock_availability"))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprint("SELECT count(item_id) FROM stock_availability"))
	tools.ErrPr(err, "")
	return sas, counts[0]
}

func UpdateStockAvailability(sas []mode.StockAvailability) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	for k := range sas {
		_, err := tx.NamedExec("UPDATE stock_availability SET  sales_user=:sales_user,img=:img,cy_sku=:cy_sku ,cy_name=:cy_name ,gtin=:gtin ,pt_sku=:pt_sku ,declaration=:declaration ,num=:num ,warehouse=:warehouse ,lead_time=:lead_time ,counts=:counts,transit_num=:transit_num,library_num=:library_num,weighted=:weighted,remarks1=:remarks1,remarks2=:remarks2 WHERE item_id=:item_id",
			map[string]interface{}{
				"sales_user":  sas[k].SalesUser,
				"img":         sas[k].Img,
				"cy_sku":      sas[k].CySku,
				"cy_name":     sas[k].CyName,
				"gtin":        sas[k].Gtin,
				"pt_sku":      sas[k].PtSku,
				"declaration": sas[k].Declaration,
				"num":         sas[k].Num,
				"warehouse":   sas[k].Warehouse,
				"lead_time":   sas[k].LeadTime,
				"counts":      sas[k].Counts,
				"transit_num": sas[k].TransitNum,
				"library_num": sas[k].LibraryNum,
				"weighted":    sas[k].Weighted,
				"remarks1":    sas[k].Remarks1,
				"remarks2":    sas[k].Remarks2,
				"item_id":     sas[k].ItemId,
			})
		if err != nil {
			tx.Rollback()
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

func UpdateStockAvailabilityImg(sas []mode.StockAvailability) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	for k := range sas {
		_, err := tx.NamedExec("UPDATE stock_availability SET  img=:img WHERE item_id=:item_id",
			map[string]interface{}{
				"img":     sas[k].Img,
				"item_id": sas[k].ItemId,
			})
		if err != nil {
			tx.Rollback()
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
func DelTockAvailabilityMsg() {

	//config.Db.Exec("UPDATE stock_availability SET  sales_user=:sales_user")

}
