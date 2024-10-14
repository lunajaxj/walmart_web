package product

import (
	"fmt"
	"log"
	"strings"
	"walmart_web/app/category"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddProducts(prs []mode.ProductDetails) error {
	prs = category.AddCategory(prs)
	//_, err := config.Db.NamedExec(`INSERT INTO product_details (id,img,code_type,code,brands,tags,title,rating,comments,price,sellers,distribution,variants1,variants2,variants_id,arrival_time,category1,category2,category3,category4,category5,category6,category7,category_name)
	//    VALUES (:id,:img,:code_type,:code,:brands,:tags,:title,:rating,:comments,:price,:sellers,:distribution,:variants1,:variants2,:variants_id,:arrival_time,:category1,:category2,:category3,:category4,:category5,:category6,:category7,:category_name)
	//    ON DUPLICATE KEY UPDATE tags = :tags, price = IF(price = 0, :price, price)`, prs)
	_, err := config.Db.NamedExec(`INSERT INTO product_details (id,img,code_type,code,brands,tags,title,rating,comments,price,sellers,distribution,variants1,variants2,variants_id,arrival_time,category1,category2,category3,category4,category5,category6,category7,category_name,star_from)
        VALUES (:id,:img,:code_type,:code,:brands,:tags,:title,:rating,:comments,:price,:sellers,:distribution,:variants1,:variants2,:variants_id,:arrival_time,:category1,:category2,:category3,:category4,:category5,:category6,:category7,:category_name,:star_from)
        ON DUPLICATE KEY UPDATE img=:img,code_type=:code_type,code=:code,brands=:brands,tags=:tags,title=:title,rating=:rating,comments=:comments,price=:price,sellers=:sellers,distribution=:distribution,variants1=:variants1,variants2=:variants2,variants_id=:variants_id,arrival_time=:arrival_time,category1=:category1,category2=:category2,category3=:category3,category4=:category4,category5=:category5,category6=:category6,category7=:category7,category_name=:category_name,star_from=:star_from`, prs)
	tools.ErrPr(err, "")
	return err
}

func AddProduct(pr mode.ProductDetails) {
	AddProducts([]mode.ProductDetails{pr})
}

func EditProductMark(id, remark, mark string) int {
	var count int
	err := config.Db.Select(&count, fmt.Sprintf("UPDATE product_details SET remark='%s', mark='%s'  where id=%s", remark, mark, id))
	tools.ErrPr(err, "")
	return count
}

func UploadProductGtin(id, gtin string) []int {
	defer func() { <-ChGtin }()
	var count []int
	err := config.Db.Select(&count, fmt.Sprintf("UPDATE product_details SET code='%s'  where id=%s", gtin, id))
	tools.ErrPr(err, "")
	return count
}

func UploadProduct(pro mode.ProductDetails) error {
	_, err := config.Db.NamedExec("UPDATE product_details SET img=:img,code_type=:code_type,code=:code,brands=:brands,tags=:tags,title=:title,rating=:rating,comments=:comments,price=:price,sellers=:sellers,distribution=:distribution,variants1=:variants1,variants2=:variants2,variants_id=:variants_id,arrival_time=:arrival_time,category1=:category1,category2=:category2,category3=:category3,category4=:category4,category5=:category5,category6=:category6,category7=:category7,category_name=:category_name  where id=:id", pro)
	tools.ErrPr(err, "")
	return err
}

func UploadProductNoIs(pro mode.ProductDetails) error {
	_, err := config.Db.NamedExec("UPDATE product_details SET remark=:remark  where id=:id", pro)
	tools.ErrPr(err, "")
	return err
}

func GETProductMark() []string {
	var marks []string
	err := config.Db.Select(&marks, "select distinct mark from product_details")
	tools.ErrPr(err, "")
	return marks
}

func GetProduct(page, limit int, id, title, sellers, sellersType, tags, brands, remark, mark, categoryTree, rating, comments, price string, actids, keyids []string) ([]mode.ProductDetails, int) {
	var prs []mode.ProductDetails
	var counts []int
	var wheres string
	if len(actids) != 0 {
		wheres = tools.WhereAndInsAndInt(wheres, "id", actids)
	}
	if len(brands) != 0 {
		split := strings.Split(brands, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("brands", len(split)),
			split)
	}
	if len(keyids) != 0 {
		wheres = tools.WhereAndInsAndInt(wheres, "id", keyids)
	}
	if id != "" {
		split := strings.Split(id, "\n")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndInsAndInt(wheres, "id", split)
	}

	if sellers != "" {
		split := strings.Split(sellers, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("sellers", len(split)),
			split)
	}
	//if sellersType != "" {
	//	if sellersType == "1" {
	//		wheres = tools.WhereAnd(wheres, "!=", "sellers", "Walmart.com")
	//	} else if sellersType == "2" {
	//		wheres = tools.WhereAnd(wheres, "=", "sellers", "Walmart.com")
	//	}
	//}
	if sellersType != "" {
		split := strings.Split(sellersType, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAnds(wheres,
			tools.WhereREPEAT("!=", len(split)),
			tools.WhereREPEAT("sellers", len(split)),
			split)
	}

	if tags != "" {
		split := strings.Split(tags, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		var spl []string
		for i := range split {
			spl = append(spl, "%"+split[i]+"%")
		}
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("LIKE", len(split)),
			tools.WhereREPEAT("tags", len(split)),
			spl)
	}
	if title != "" {
		split := strings.Split(title, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		var spl []string
		for i := range split {
			spl = append(spl, "%"+split[i]+"%")
		}
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("LIKE", len(split)),
			tools.WhereREPEAT("title", len(split)),
			spl)
	}
	if categoryTree != "" {
		wheres = tools.WhereAndOrs(wheres,
			[]string{"LIKE", "LIKE", "LIKE", "LIKE", "LIKE", "LIKE", "LIKE"},
			[]string{"category1", "category2", "category3", "category4", "category5", "category6", "category7"},
			[]string{"%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%"})
	}
	if remark != "" {
		split := strings.Split(remark, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		var spl []string
		for i := range split {
			spl = append(spl, "%"+split[i]+"%")
		}
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("LIKE", len(split)),
			tools.WhereREPEAT("remark", len(split)),
			spl)
	}
	if mark != "" {
		split := strings.Split(mark, "|")
		split = tools.RemoveEmptyStringsFromArray(split)
		wheres = tools.WhereAndOrs(wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("mark", len(split)),
			split)
	}
	if rating != "" {
		split := strings.Split(rating, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "rating", rating)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"rating", "rating"}, split)
		}
	}
	if comments != "" {
		split := strings.Split(comments, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "comments", comments)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"comments", "comments"}, split)
		}
	}
	if price != "" {
		split := strings.Split(price, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "price", price)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"price", "price"}, split)
		}
	}

	page = (page - 1) * limit
	sql := fmt.Sprintf("SELECT * FROM product_details %s LIMIT %d,%d", wheres, page, limit)
	err := config.Db.Select(&prs, sql)
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(id) FROM product_details %s", wheres))
	tools.ErrPr(err, "")
	return prs, counts[0]
}

func GetProductById(sellersType, categoryTree, rating, comments, price string) ([]string, int) {
	var prs []string
	var counts []int
	var wheres string

	if sellersType != "" {
		if sellersType == "1" {
			wheres = tools.WhereAnd(wheres, "!=", "sellers", "Walmart.com")
		} else if sellersType == "2" {
			wheres = tools.WhereAnd(wheres, "=", "sellers", "Walmart.com")
		}
	}
	if categoryTree != "" {
		wheres = tools.WhereAndOrs(wheres,
			[]string{"LIKE", "LIKE", "LIKE", "LIKE", "LIKE", "LIKE", "LIKE"},
			[]string{"category1", "category2", "category3", "category4", "category5", "category6", "category7"},
			[]string{"%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%", "%" + categoryTree + "%"})
	}
	if rating != "" {
		split := strings.Split(rating, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "rating", rating)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"rating", "rating"}, split)
		}
	}
	if comments != "" {
		split := strings.Split(comments, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "comments", comments)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"comments", "comments"}, split)
		}
	}
	if price != "" {
		split := strings.Split(price, "-")
		split = tools.RemoveEmptyStringsFromArray(split)
		if len(split) == 1 {
			wheres = tools.WhereAnd(wheres, "=", "price", price)
		} else {
			wheres = tools.WhereAnds(wheres, []string{">=", "<="}, []string{"price", "price"}, split)
		}
	}

	sql := fmt.Sprintf("SELECT id FROM product_details %s ", wheres)
	err := config.Db.Select(&prs, sql)
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(id) FROM product_details %s", wheres))
	tools.ErrPr(err, "")
	return prs, counts[0]
}

func GetProductId(ids []string) []mode.ProductDetails {
	var prs []mode.ProductDetails
	var wheres string
	if len(ids) != 0 {
		wheres = tools.WhereAndOrsInt(wheres,
			tools.WhereREPEAT("=", len(ids)),
			tools.WhereREPEAT("id", len(ids)),
			ids)
	}
	err := config.Db.Select(&prs, fmt.Sprintf("SELECT * FROM product_details %s", wheres))
	tools.ErrPr(err, fmt.Sprintf("SELECT * FROM product_details %s", wheres))
	return prs
}

func GetChart(brands []string, m, bfb1, bfb2 int) ([]string, map[string]mode.ProductBrandss) {
	for i := range brands {
		brands[i] = strings.ToLower(brands[i])
	}
	var sellers1 []mode.ProductBrands
	sellers2 := make(map[string][]mode.ProductBrands)
	uniqueMap := make(map[string]int)
	for i := range brands {
		brands[i] = strings.Replace(brands[i], "'", `\'`, -1)
	}
	wheres := "'" + strings.Join(brands, "', '") + "'"
	err := config.Db.Select(&sellers1, fmt.Sprintf(`select  LOWER(sellers) as sellers ,LOWER(brands) as brands,COUNT(*) AS count from product_details where brands IN (%s) and sellers !="" and sellers !="Walmart.com" GROUP BY sellers HAVING count > 1`, wheres))
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	wheres = ""
	//改
	//var sellerss []string

	for i := range sellers1 {
		//改
		//sellerss = append(sellerss, sellers1[i].Sellers)

		if wheres != "" {
			wheres += ",'" + strings.Replace(sellers1[i].Sellers, "'", `\'`, -1) + "'"
		} else {
			wheres += "'" + strings.Replace(sellers1[i].Sellers, "'", `\'`, -1) + "'"
		}

	}

	//改
	//return GetChartSellers(sellerss)

	type SellerCount struct {
		Sellers string `db:"sellers"`
		Count   int    `db:"count"`
	}
	var count []SellerCount
	err = config.Db.Select(&count, fmt.Sprintf(`select LOWER(sellers) as sellers,COUNT(*) AS count from product_details where sellers in (%s) GROUP BY sellers`, wheres))
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	wheres = ""
	for i := range sellers1 {
		for _, v := range count {
			if sellers1[i].Sellers == v.Sellers && sellers1[i].Count*100 > (v.Count)*bfb1 {
				if wheres != "" {
					wheres += ",'" + strings.Replace(sellers1[i].Sellers, "'", `\'`, -1) + "'"
				} else {
					wheres += "'" + strings.Replace(sellers1[i].Sellers, "'", `\'`, -1) + "'"
				}
				break
			}
		}
	}

	var sellers0 []mode.ProductBrands
	iss := make(map[string]bool)
	err = config.Db.Select(&sellers0, fmt.Sprintf(`SELECT LOWER(sellers) as sellers ,LOWER(brands) as brands, COUNT(*) as count FROM product_details where  sellers IN (%s)  and brands != '",' GROUP BY sellers,brands`, wheres))
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	for i2 := range sellers0 {
		if _, ok := iss[sellers0[i2].Sellers+sellers0[i2].Brands]; !ok {
			uniqueMap[sellers0[i2].Brands] += 1
			iss[sellers0[i2].Sellers+sellers0[i2].Brands] = true
		}
	}

	for s := range uniqueMap {
		if uniqueMap[s]*100 > len(sellers1)*bfb2 {
			brands = append(brands, strings.TrimSpace(s))
		}
	}
	brands = tools.UniqueArr(brands)
	if m == 1 {
		return GetChart(brands, 2, bfb1, bfb2)
	}

fo:
	for i := range sellers0 {
		for i2 := range brands {
			if sellers0[i].Brands == brands[i2] {
				sellers2[sellers0[i].Sellers] = append(sellers2[sellers0[i].Sellers], sellers0[i])
				continue fo
			}
		}

	}
	ps := make(map[string]mode.ProductBrandss)
	for s := range sellers2 {
		var bs []int
	fo1:
		for i := range brands {
			for ii := range sellers2[s] {
				if brands[i] == sellers2[s][ii].Brands {
					bs = append(bs, sellers2[s][ii].Count)
					continue fo1
				}
			}
			bs = append(bs, 0)
		}
		ps[s] = mode.ProductBrandss{
			Sellers: s,
			Brandss: bs,
		}

	}

	return brands, ps
}
func GetChartSellers(sellers []string) ([]string, map[string]mode.ProductBrandss) {
	sellers2 := make(map[string][]mode.ProductBrands)
	var wheres string
	for i := range sellers {
		if wheres != "" {
			wheres += ",'" + strings.Replace(sellers[i], "'", `\'`, -1) + "'"
		} else {
			wheres += "'" + strings.Replace(sellers[i], "'", `\'`, -1) + "'"
		}

	}
	type SellerCount struct {
		Sellers string `db:"sellers"`
		Count   int    `db:"count"`
	}
	var count []SellerCount
	err := config.Db.Select(&count, fmt.Sprintf(`select LOWER(sellers) as sellers,COUNT(*) AS count from product_details where sellers in (%s) GROUP BY sellers`, wheres))
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	var sellers0 []mode.ProductBrands
	err = config.Db.Select(&sellers0, fmt.Sprintf(`SELECT LOWER(sellers) as sellers ,LOWER(brands) as brands, COUNT(*) as count FROM product_details where  sellers IN (%s)  and brands != '",' GROUP BY sellers,brands`, wheres))
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	var brands []string
	for i := range sellers0 {
		brands = append(brands, sellers0[i].Brands)
		sellers2[sellers0[i].Sellers] = append(sellers2[sellers0[i].Sellers], sellers0[i])
	}
	brands = tools.UniqueArr(brands)
	ps := make(map[string]mode.ProductBrandss)
	for s := range sellers2 {
		var bs []int
	fo1:
		for i := range brands {
			for ii := range sellers2[s] {
				if brands[i] == sellers2[s][ii].Brands {
					bs = append(bs, sellers2[s][ii].Count)
					continue fo1
				}
			}
			bs = append(bs, 0)
		}
		ps[s] = mode.ProductBrandss{
			Sellers: s,
			Brandss: bs,
		}

	}

	return brands, ps
}
