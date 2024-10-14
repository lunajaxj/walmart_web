package storeInformation

import (
	"fmt"
	"log"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func GetStoreInformation(page, limit int, name, sid, affiliatedCompany, hasTargetLink, onSaleProductCount string) ([]mode.StoreInformation, int) {
	var sas []mode.StoreInformation
	var counts []int
	var wheres string
	if name != "" {
		split := strings.Split(name, "\n")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("account_shop_name", len(split)),
			split)
	}
	if sid != "" {
		wheres = tools.WhereAnd(wheres, "=", "sid", sid)
	}
	if affiliatedCompany != "" {
		wheres = tools.WhereAnd(wheres, "=", "affiliated_company", affiliatedCompany)
	}
	if hasTargetLink != "" {
		split := strings.Split(hasTargetLink, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "has_target_link", hasTargetLink)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"has_target_link", "has_target_link"}, split)
		}
	}
	if onSaleProductCount != "" {
		split := strings.Split(onSaleProductCount, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "on_sale_product_count", onSaleProductCount)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"on_sale_product_count", "on_sale_product_count"}, split)
		}
	}

	page = (page - 1) * limit
	err := config.Db.Select(&sas, fmt.Sprintf("SELECT * FROM store_information %s LIMIT %d,%d", wheres, page, limit))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(shop_id) FROM store_information %s", wheres))
	tools.ErrPr(err, "")
	return sas, counts[0]
}
func Remove(sids string) int {
	split := strings.Split(sids, ",")
	var wheres string
	if len(split) == 0 {
		return 0
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "sid", split[i])
	}
	updateResult := config.Db.MustExec(fmt.Sprintf("DELETE FROM store_information %s", wheres))
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func Install(sas []mode.StoreInformation) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Println(err)
	}
	// 执行批量插入
	for _, sa := range sas {
		_, err = tx.NamedExec(`
		INSERT INTO store_information (
			account_shop_name,
			shop_id,
			pid,
			on_sale_product_count,
			has_target_link,
			bsr_link,
			pp_link,
			wfs_delivery_count,
			wfs_delivery_percentage,
			price_0_10,
			price_10_15,
			price_15_20,
			price_20_40,
			price_40_60,
			price_60_above,
			reviews_0_5,
			reviews_5_10,
			reviews_10_20,
			reviews_20_100,
			reviews_100_above,
			brand,
			affiliated_company,
		    seller_name,
			address,
			country,
			readx,
			seller_reviews_num,
			note_1,
			note_2,
			note_3,
			note_4,
			note_5
		) VALUES (
			:account_shop_name,
			:shop_id,
			:pid,
			:on_sale_product_count,
			:has_target_link,
			:bsr_link,
			:pp_link,
			:wfs_delivery_count,
			:wfs_delivery_percentage,
			:price_0_10,
			:price_10_15,
			:price_15_20,
			:price_20_40,
			:price_40_60,
			:price_60_above,
			:reviews_0_5,
			:reviews_5_10,
			:reviews_10_20,
			:reviews_20_100,
			:reviews_100_above,
			:brand,
			:affiliated_company,
		    :seller_name,
			:address,
			:country,
			:readx,
			:seller_reviews_num,
			:note_1,
			:note_2,
			:note_3,
			:note_4,
			:note_5
		)
	ON DUPLICATE KEY UPDATE shop_id=:shop_id,pid=:pid
	`, sa)
		if err != nil {
			log.Println("插入数据失败:", err)
		}
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return false
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
func GetStoreInformationAll() ([]mode.StoreInformation, int) {
	var sas []mode.StoreInformation
	var counts []int

	err := config.Db.Select(&sas, fmt.Sprint("SELECT * FROM store_information"))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprint("SELECT count(shop_id) FROM store_information"))
	tools.ErrPr(err, "")
	return sas, counts[0]
}

func UpdateStoreInformation(sas []mode.StoreInformation) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Println(err)
	}
	for k := range sas {
		_, err = tx.NamedExec(`
	UPDATE store_information SET
		pid = :pid,
		shop_id = :shop_id,
		on_sale_product_count = :on_sale_product_count,
		has_target_link = :has_target_link,
		bsr_link = :bsr_link,
		pp_link = :pp_link,
		wfs_delivery_count = :wfs_delivery_count,
		wfs_delivery_percentage = :wfs_delivery_percentage,
		price_0_10 = :price_0_10,
		price_10_15 = :price_10_15,
		price_15_20 = :price_15_20,
		price_20_40 = :price_20_40,
		price_40_60 = :price_40_60,
		price_60_above = :price_60_above,
		reviews_0_5 = :reviews_0_5,
		reviews_5_10 = :reviews_5_10,
		reviews_10_20 = :reviews_10_20,
		reviews_20_100 = :reviews_20_100,
		reviews_100_above = :reviews_100_above,
		brand = :brand,
		affiliated_company = :affiliated_company,
		note_1 = :note_1,
		note_2 = :note_2,
		note_3 = :note_3,
		note_4 = :note_4,
		note_5 = :note_5
	WHERE account_shop_name = :account_shop_name
`, sas[k])
		if err != nil {
			log.Println("3更新数据失败:", sas[k].AccountShopName, err)
		}
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return false
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func UpdateStoreInformationNoteAndAffiliatedCompany(sas []mode.StoreInformation) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Println(err)
	}
	for k := range sas {
		_, err = tx.NamedExec(`
	UPDATE store_information SET
		affiliated_company = :affiliated_company,
		note_1 = :note_1,
		note_2 = :note_2,
		note_3 = :note_3,
		note_4 = :note_4,
		note_5 = :note_5
	WHERE shop_id = :shop_id
`, sas[k])
		if err != nil {
			log.Println("2更新数据失败:", sas[k].AccountShopName, err)
		}
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return false
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func UpdateStoreInformationNote(sas []mode.StoreInformation) bool {
	// 开始事务
	tx, err := config.Db.Beginx()
	if err != nil {
		log.Println(err)
	}
	for k := range sas {
		_, err = tx.NamedExec(`
	UPDATE store_information SET
		affiliated_company= :affiliated_company,                        
		note_1 = :note_1,
		note_2 = :note_2,
		note_3 = :note_3,
		note_4 = :note_4,
		note_5 = :note_5
	WHERE sid = :sid
`, sas[k])
		if err != nil {
			log.Println("1更新数据失败:", sas[k].AccountShopName, err)
		}
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return false
		}
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func UpdateInfo(sa mode.StoreInformation) {
	_, err := config.Db.NamedExec(`
	UPDATE store_information SET
		seller_name = :seller_name,
		address = :address,
		country = :country,
		readx = :readx,
		seller_reviews_num =:seller_reviews_num
	WHERE shop_id = :shop_id
`, sa)
	if err != nil {
		log.Println("4更新数据失败:", sa.AccountShopName, err)
	}
	if err != nil {
		log.Println(err)

	}

}
