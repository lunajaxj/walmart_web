package category

import (
	"fmt"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

//func GetCategory() []*mode.Category {
//	var cats []*mode.Category
//	err := config.Db.Select(&cats, fmt.Sprintf("SELECT * FROM category"))
//	tools.ErrPr(err,"")
//	for i := range cats {
//		for ii := range cats {
//			if cats[i].CategoryName == cats[ii].CategoryUpName {
//				cats[i].Children = append(cats[i].Children, cats[ii])
//			}
//		}
//	}
//	return cats
//}

func GetCategory(name string) []map[string]interface{} {
	var cats1 []*mode.Category
	var where1 string
	if name == "" {
		where1 = tools.WhereAnd(where1, "=", "category_up_name", "")
		err := config.Db.Select(&cats1, fmt.Sprintf("SELECT * FROM category %s ", where1))
		tools.ErrPr(err, "")
	} else {
		where1 = tools.WhereAnd(where1, "=", "category_up_name", name)
		err := config.Db.Select(&cats1, fmt.Sprintf("SELECT * FROM category %s ", where1))
		tools.ErrPr(err, "")
	}
	var cats2 []*mode.Category
	var where2 string
	for i := range cats1 {
		where2 = tools.WhereOr(where2, "=", "category_up_name", cats1[i].CategoryName)
	}
	err := config.Db.Select(&cats2, fmt.Sprintf("SELECT * FROM category %s ", where2))
	tools.ErrPr(err, "")
	var m []map[string]interface{}
	for i := range cats1 {
		node := map[string]interface{}{"value": cats1[i].CategoryName, "label": cats1[i].CategoryName}
		node["leaf"] = true
		for i2 := range cats2 {
			if cats1[i].CategoryName == cats2[i2].CategoryUpName {
				node["leaf"] = false
			}
		}
		m = append(m, node)
	}

	return m
}

func AddCategory(prs []mode.ProductDetails) []mode.ProductDetails {
	for i := range prs {
		var categoryNames []string
		var categoryUp string
		var category string
		for ii := 7; ii >= 1; ii-- {
			switch ii {
			case 1:
				categoryUp = ""
				category = prs[i].Category1
			case 2:
				categoryUp = prs[i].Category1
				category = prs[i].Category2
			case 3:
				categoryUp = prs[i].Category2
				category = prs[i].Category3
			case 4:
				categoryUp = prs[i].Category3
				category = prs[i].Category4
			case 5:
				categoryUp = prs[i].Category4
				category = prs[i].Category5
			case 6:
				categoryUp = prs[i].Category5
				category = prs[i].Category6
			case 7:
				categoryUp = prs[i].Category6
				category = prs[i].Category7
			}
			if category != "" {
				categoryc := strings.Replace(category, "'", `\'`, -1)
				categoryc = strings.Replace(categoryc, `"`, `\"`, -1)
				categoryUpc := strings.Replace(categoryUp, "'", `\'`, -1)
				categoryUpc = strings.Replace(categoryUpc, `"`, `\"`, -1)
				err := config.Db.Select(&categoryNames, fmt.Sprintf("SELECT category_up_name FROM category WHERE category_name = '%s' AND category_up_name = '%s'", categoryc, categoryUpc))
				tools.ErrPr(err, "")
				if len(categoryNames) == 0 {
					config.Db.NamedExec("INSERT INTO category (category_up_name,category_name) VALUES (:categoryUpName,:categoryName)",
						map[string]interface{}{
							"categoryUpName": categoryUp,
							"categoryName":   category,
						})
				}
			}
		}
	}
	return prs
}
