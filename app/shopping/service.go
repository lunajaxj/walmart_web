package shopping

import (
	"fmt"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddShopping(sho mode.Shopping) {
	var up = ""
	if sho.Sku != "" {
		up += "sku = :sku,"
	}
	if sho.FloorPrice != "" {
		up += "floor_price = :floor_price,"
	}
	if sho.XPrice != "" {
		up += "xprice = :xprice,"
	}
	if sho.Inventory != "" {
		up += "inventory = :inventory,"
	}
	if sho.CenterId != "" {
		up += "center_id = :center_id,"
	}
	if sho.Seller != "" {
		up += "seller = :seller,"
	}
	if sho.Note != "" {
		up += "note = :note,"
	}
	if sho.Name != "" {
		up += "name = :name,"
	}
	if sho.Sales != 159057 {
		up += "sales = :sales,"
	}
	if sho.ShoppingCron != "" {
		up += "shopping_cron = :shopping_cron,"
	}
	if sho.TheShelvesCron != "" {
		up += "the_shelves_cron = :the_shelves_cron,"
	}
	if sho.InventoryCron != "" {
		up += "inventory_cron = :inventory_cron,"
	}
	if sho.XInventoryCron != "" {
		up += "xinventory_cron = :xinventory_cron,"
	}
	if sho.StatusCron1 != "" {
		up += "status_cron1 = :status_cron1,"
	}
	if sho.StatusCron2 != "" {
		up += "status_cron2 = :status_cron2,"
	}
	if sho.StatusCron3 != "" {
		up += "status_cron3 = :status_cron3,"
	}
	if sho.StatusCron4 != "" {
		up += "status_cron4 = :status_cron4,"
	}
	if sho.StatusCron5 != "" {
		up += "status_cron5 = :status_cron5,"
	}
	if sho.StatusCron5 != "" {
		up += "status_cron5 = :status_cron5,"
	}

	if up != "" {
		up = up[:len(up)-1]
	}
	_, err := config.Db.NamedExec(fmt.Sprintf(`INSERT INTO shopping (pr_id,img,sku,price,xprice,inventory,center_id,seller,floor_price,is_active,promotions_status,promo_price,promo_start_date,promo_end_date,is_up,msg,status1,status2,status3,status4,status5,note,name,sales,shopping_cron,the_shelves_cron,inventory_cron,xinventory_cron,status_cron1,status_cron2,status_cron3,status_cron4,status_cron5,j99991,j99992,j99993,j99994,j99995)
       									VALUES (:pr_id,:img,:sku,:price,:xprice,:inventory,:center_id,:seller,:floor_price,:is_active,:promotions_status,:promo_price,:promo_start_date,:promo_end_date,:is_up,:msg,:status1,:status2,:status3,:status4,:status5,:note,:name,:sales,:shopping_cron,:the_shelves_cron,:inventory_cron,:xinventory_cron,:status_cron1,:status_cron2,:status_cron3,:status_cron4,:status_cron5,:j99991,:j99992,:j99993,:j99994,:j99995)
       									ON DUPLICATE KEY UPDATE %s`, up), sho)
	tools.ErrPr(err, "")
}

func GetShopping(page, limit int, id, note, name, msg, seller, sales, status1, status2, status3, status4, status5, j99991, j99992, j99993, j99994, j99995 string) ([]mode.Shopping, int) {
	var shoss []mode.Shopping
	var counts []int
	var wheres string
	if id != "" {
		split := strings.Split(id, "\n")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("pr_id", len(split)),
			split)
	}
	if j99991 == "有" {
		wheres = tools.WhereAnd(wheres, "!=", "j99991", "无")
	} else if j99991 == "无" {
		wheres = tools.WhereAnd(wheres, "=", "j99991", "无")
	}
	if j99992 == "有" {
		wheres = tools.WhereAnd(wheres, "!=", "j99992", "无")
	} else if j99992 == "无" {
		wheres = tools.WhereAnd(wheres, "=", "j99992", "无")
	}
	if j99993 == "有" {
		wheres = tools.WhereAnd(wheres, "!=", "j99993", "无")
	} else if j99993 == "无" {
		wheres = tools.WhereAnd(wheres, "=", "j99993", "无")
	}
	if j99994 == "有" {
		wheres = tools.WhereAnd(wheres, "!=", "j99994", "无")
	} else if j99994 == "无" {
		wheres = tools.WhereAnd(wheres, "=", "j99994", "无")
	}
	if j99995 == "有" {
		wheres = tools.WhereAnd(wheres, "!=", "j99995", "无")
	} else if j99995 == "无" {
		wheres = tools.WhereAnd(wheres, "=", "j99995", "无")
	}
	if sales != "" {
		split := strings.Split(sales, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "sales", sales)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"sales", "sales"}, split)
		}
	}
	if note != "" {
		wheres = tools.WhereAnd(wheres, "LIKE", "note", "%"+note+"%")
	}
	if status1 != "" {
		if status1 == " " {
			status1 = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "status1", status1)
	}
	if status2 != "" {
		if status2 == " " {
			status2 = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "status2", status2)
	}
	if status3 != "" {
		if status3 == " " {
			status3 = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "status3", status3)
	}
	if status4 != "" {
		if status4 == " " {
			status4 = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "status4", status4)
	}
	if status5 != "" {
		if status5 == " " {
			status5 = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "status5", status5)
	}
	if name != "" {
		wheres = tools.WhereAnd(wheres, "LIKE", "name", "%"+name+"%")
	}
	if seller != "" {
		wheres = tools.WhereAnd(wheres, "=", "seller", seller)
	}
	if msg != "" {
		if msg == " " {
			msg = ""
		}
		wheres = tools.WhereAnd(wheres, "=", "msg", msg)
	}
	page = (page - 1) * limit
	err := config.Db.Select(&shoss, fmt.Sprintf("SELECT * FROM shopping %s LIMIT %d,%d", wheres, page, limit))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(pr_id) FROM shopping %s", wheres))
	tools.ErrPr(err, "")
	return shoss, counts[0]
}

func GetShoppingByIds(id []string) []mode.Shopping {
	var shoss []mode.Shopping
	var wheres string
	if len(id) > 0 {
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(id)),
			tools.WhereREPEAT("pr_id", len(id)),
			id)
	} else {
		return shoss
	}
	err := config.Db.Select(&shoss, fmt.Sprintf("SELECT * FROM shopping %s", wheres))
	tools.ErrPr(err, "")
	return shoss
}

func GetShoppingByCron(name, cron, seller string) []mode.Shopping {
	var shoss []mode.Shopping
	err := config.Db.Select(&shoss, fmt.Sprintf("SELECT * FROM shopping where %s = '%s' and seller = '%s'", name, cron, seller))
	tools.ErrPr(err, "")
	return shoss
}
func GetShoppingCron(name, seller string) []string {
	var list []string
	err := config.Db.Select(&list, fmt.Sprintf("SELECT distinct %s FROM shopping where %s != '' and seller = '%s'", name, name, seller))
	tools.ErrPr(err, "")
	return list
}

func GetShoppingMsg() []string {
	var list []string
	err := config.Db.Select(&list, fmt.Sprintf("SELECT distinct msg FROM shopping"))
	tools.ErrPr(err, "")
	return list
}

func GetShoppingSeller() []string {
	var list []string
	err := config.Db.Select(&list, fmt.Sprintf("SELECT distinct seller FROM shopping"))
	tools.ErrPr(err, "")
	return list
}

func UploadShopping(sho mode.Shopping) int {
	stmt := "UPDATE shopping SET sku=:sku, price=:price, xprice=:xprice, inventory=:inventory, center_Id=:center_Id, seller=:seller, floor_price=:floor_price, is_active=:is_active, promotions_status=:promotions_status, promo_price=:promo_price, promo_start_date=:promo_start_date, promo_end_date=:promo_end_date, is_up=:is_up, msg=:msg,img=:img,name=:name,sales=:sales WHERE pr_id=:pr_id"
	updateResult, err := config.Db.NamedExec(stmt, map[string]interface{}{
		"sku":               sho.Sku,
		"price":             sho.Price,
		"xprice":            sho.XPrice,
		"inventory":         sho.Inventory,
		"center_Id":         sho.CenterId,
		"seller":            sho.Seller,
		"floor_price":       sho.FloorPrice,
		"is_active":         sho.IsActive,
		"promotions_status": sho.PromotionsStatus,
		"promo_price":       sho.PromoPrice,
		"promo_start_date":  sho.PromoStartDate,
		"promo_end_date":    sho.PromoEndDate,
		"is_up":             sho.IsUp,
		"msg":               sho.Msg,
		"name":              sho.Name,
		"sales":             sho.Sales,
		"img":               sho.Img,
		"pr_id":             sho.PrId,
	})
	tools.ErrPr(err, "")
	//updateResult := config.Db.MustExec("UPDATE shopping SET sku=?,price=?,xprice=?,inventory=?,center_Id=?,seller=?,floor_price=?,is_active=?,promotions_status=?,promo_price=?,promo_start_date=?,promo_end_date=?,is_up=?,msg=?,img=? where pr_id=?", sho.Sku, sho.Price, sho.XPrice, sho.Inventory, sho.CenterId, sho.Seller, sho.FloorPrice, sho.IsActive, sho.PromotionsStatus, sho.PromoPrice, sho.PromoStartDate, sho.PromoEndDate, sho.IsUp, sho.Msg, sho.Img, sho.PrId)
	if err != nil && updateResult != nil {
		affected, err := updateResult.RowsAffected()
		tools.ErrPr(err, "")
		if err != nil {
			return int(affected)
		}
	}
	return 0
}

func UploadShoppingAll(sho mode.Shopping) int {
	stmt := "UPDATE shopping SET sku=:sku, price=:price, xprice=:xprice, inventory=:inventory, center_Id=:center_Id, seller=:seller, floor_price=:floor_price, is_active=:is_active, promotions_status=:promotions_status, promo_price=:promo_price, promo_start_date=:promo_start_date, promo_end_date=:promo_end_date, is_up=:is_up, msg=:msg,img=:img,name=:name,sales=:sales,status1=:status1,status2=:status2,status3=:status3,status4=:status4,status5=:status5 WHERE pr_id=:pr_id"
	updateResult, err := config.Db.NamedExec(stmt, map[string]interface{}{
		"sku":               sho.Sku,
		"price":             sho.Price,
		"xprice":            sho.XPrice,
		"inventory":         sho.Inventory,
		"center_Id":         sho.CenterId,
		"seller":            sho.Seller,
		"floor_price":       sho.FloorPrice,
		"is_active":         sho.IsActive,
		"promotions_status": sho.PromotionsStatus,
		"promo_price":       sho.PromoPrice,
		"promo_start_date":  sho.PromoStartDate,
		"promo_end_date":    sho.PromoEndDate,
		"is_up":             sho.IsUp,
		"msg":               sho.Msg,
		"name":              sho.Name,
		"sales":             sho.Sales,
		"img":               sho.Img,
		"status1":           sho.Status1,
		"status2":           sho.Status2,
		"status3":           sho.Status3,
		"status4":           sho.Status4,
		"status5":           sho.Status5,
		"pr_id":             sho.PrId,
	})
	tools.ErrPr(err, "")
	//updateResult := config.Db.MustExec("UPDATE shopping SET sku=?,price=?,xprice=?,inventory=?,center_Id=?,seller=?,floor_price=?,is_active=?,promotions_status=?,promo_price=?,promo_start_date=?,promo_end_date=?,is_up=?,msg=?,img=? where pr_id=?", sho.Sku, sho.Price, sho.XPrice, sho.Inventory, sho.CenterId, sho.Seller, sho.FloorPrice, sho.IsActive, sho.PromotionsStatus, sho.PromoPrice, sho.PromoStartDate, sho.PromoEndDate, sho.IsUp, sho.Msg, sho.Img, sho.PrId)
	if err != nil && updateResult != nil {
		affected, err := updateResult.RowsAffected()
		tools.ErrPr(err, "")
		if err != nil {
			return int(affected)
		}
	}
	return 0
}

func UploadShoppingStatus(sho mode.Shopping) int {
	stmt := "UPDATE shopping SET status1=:status1,status2=:status2,status3=:status3,status4=:status4,status5=:status5,img=:img,price=:price,j99991=:j99991,j99992=:j99992,j99993=:j99993,j99994=:j99994,j99995=:j99995 WHERE pr_id=:pr_id"
	updateResult, err := config.Db.NamedExec(stmt, map[string]interface{}{
		"price":   sho.Price,
		"img":     sho.Img,
		"status1": sho.Status1,
		"status2": sho.Status2,
		"status3": sho.Status3,
		"status4": sho.Status4,
		"status5": sho.Status5,
		"j99991":  sho.J99991,
		"j99992":  sho.J99992,
		"j99993":  sho.J99993,
		"j99994":  sho.J99994,
		"j99995":  sho.J99995,
		"pr_id":   sho.PrId,
	})
	tools.ErrPr(err, "")
	//updateResult := config.Db.MustExec("UPDATE shopping SET sku=?,price=?,xprice=?,inventory=?,center_Id=?,seller=?,floor_price=?,is_active=?,promotions_status=?,promo_price=?,promo_start_date=?,promo_end_date=?,is_up=?,msg=?,img=? where pr_id=?", sho.Sku, sho.Price, sho.XPrice, sho.Inventory, sho.CenterId, sho.Seller, sho.FloorPrice, sho.IsActive, sho.PromotionsStatus, sho.PromoPrice, sho.PromoStartDate, sho.PromoEndDate, sho.IsUp, sho.Msg, sho.Img, sho.PrId)
	if err != nil && updateResult != nil {
		affected, err := updateResult.RowsAffected()
		tools.ErrPr(err, "")
		if err != nil {
			return int(affected)
		}
	}
	return 0
}

func UploadSho(sho mode.Shopping) int {
	updateResult := config.Db.MustExec("UPDATE shopping SET seller=?,floor_price=?,xprice=?,inventory=?,note=?,name=?,sales=?,shopping_cron=?,the_shelves_cron=?,inventory_cron=?,xinventory_cron=?,status_cron1=?,status_cron2=?,status_cron3=?,status_cron4=?,status_cron5=? where pr_id=?", sho.Seller, sho.FloorPrice, sho.XPrice, sho.Inventory, sho.Note, sho.Name, sho.Sales, sho.ShoppingCron, sho.TheShelvesCron, sho.InventoryCron, sho.XInventoryCron, sho.StatusCron1, sho.StatusCron2, sho.StatusCron3, sho.StatusCron4, sho.StatusCron5, sho.PrId)
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func Remove(ids string) int {
	split := strings.Split(ids, ",")
	var wheres string
	if len(split) == 0 {
		return 0
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "pr_id", split[i])
	}
	updateResult := config.Db.MustExec(fmt.Sprintf("DELETE FROM shopping %s", wheres))
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func DelMsg() {
	config.Db.MustExec("UPDATE shopping SET msg = '',status1='',status2='',status3='',status4='',status5=''")
}
