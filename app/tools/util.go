package tools

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"walmart_web/app/mode"
	"walmart_web/app/walLog"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器
}

func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ErrPr(err error, sql string) bool {
	if err != nil {
		if strings.Contains(err.Error(), "' for key '") {
			return true
		}
		log.Println("sql错误：", err)
		lo := mode.Log{Classify: "sql", Msg: "sql错误：" + err.Error(), Val: sql}
		walLog.AddLog(lo)
		return true
	}
	return false
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return s
	}
	return d
}

// 文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func WhereOr(str, whe, key, value string) string {
	value = strings.Replace(value, "'", `\'`, -1)
	value = strings.Replace(value, `"`, `\"`, -1)
	if str == "" {
		str = "where " + key + " " + whe + " '" + value + "' "
	} else {
		str += "OR " + key + " " + whe + " '" + value + "' "
	}
	return str
}
func WhereOrInt(str, whe, key, value string) string {
	value = strings.Replace(value, "'", `\'`, -1)
	value = strings.Replace(value, `"`, `\"`, -1)
	if str == "" {
		str = "where " + key + " " + whe + " " + value + " "
	} else {
		str += "OR " + key + " " + whe + " " + value + " "
	}
	return str
}

func SafeDeleteFile(filePath string) error {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 指定一个备份目录
	backupDir := "./backup"

	// 确保备份目录存在
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// 生成备份文件名，使用时间戳作为文件名
	backupFileName := fmt.Sprintf("%d%s", timestamp, filepath.Base(filePath))
	backupFilePath := filepath.Join(backupDir, backupFileName)

	// 移动文件到备份目录
	err = os.Rename(filePath, backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to move file to backup directory: %v", err)
	}

	return nil
}
func WhereAndOrs(str string, whe, key, value []string) string {
	for i := range value {
		value[i] = strings.Replace(value[i], "'", `\'`, -1)
		value[i] = strings.Replace(value[i], `"`, `\"`, -1)
	}
	for i := range value {
		if i == 0 && str == "" {
			str = "where (" + key[i] + " " + whe[i] + " '" + value[i] + "' "
		} else if i == 0 && str != "" {
			str += "AND ( " + key[i] + " " + whe[i] + " '" + value[i] + "' "
		} else {
			str += "OR " + key[i] + " " + whe[i] + " '" + value[i] + "' "
		}
	}
	if len(value) > 0 {
		str += ") "
	}
	return str
}
func WhereAndOrsInt(str string, whe, key, value []string) string {
	for i := range value {
		value[i] = strings.Replace(value[i], "'", `\'`, -1)
		value[i] = strings.Replace(value[i], `"`, `\"`, -1)
	}
	for i := range value {
		if i == 0 && str == "" {
			str = "where (" + key[i] + " " + whe[i] + " " + value[i] + " "
		} else if i == 0 && str != "" {
			str += "AND ( " + key[i] + " " + whe[i] + " " + value[i] + " "
		} else {
			str += "OR " + key[i] + " " + whe[i] + " " + value[i] + " "
		}
	}
	if len(value) > 0 {
		str += ") "
	}
	return str
}

func WhereAndInsAndInt(str, key string, value []string) string {
	for i := range value {
		value[i] = strings.Replace(value[i], "'", `\'`, -1)
		value[i] = strings.Replace(value[i], `"`, `\"`, -1)
	}
	for i := range value {
		if i == 0 && str == "" {
			str = "where " + key + " IN (" + value[i]
		} else if i == 0 && str != "" {
			str += " AND  " + key + " IN (" + value[i]
		} else {
			str += "," + value[i]
		}
	}
	if len(value) > 0 {
		str += ") "
	}
	return str
}
func WhereAnds(str string, whe, key, value []string) string {
	for i := range value {
		value[i] = strings.Replace(value[i], "'", `\'`, -1)
		value[i] = strings.Replace(value[i], `"`, `\"`, -1)
	}
	for i := range value {
		if i == 0 && str == "" {
			str = "where (" + key[i] + " " + whe[i] + " '" + value[i] + "' "
		} else if i == 0 && str != "" {
			str += "AND ( " + key[i] + " " + whe[i] + " '" + value[i] + "' "
		} else {
			str += "AND " + key[i] + " " + whe[i] + " '" + value[i] + "' "
		}
	}
	if len(value) > 0 {
		str += ") "
	}
	return str
}
func WhereAnd(str, whe, key, value string) string {
	value = strings.Replace(value, "'", `\'`, -1)
	value = strings.Replace(value, `"`, `\"`, -1)
	if str == "" {
		str = "where " + key + " " + whe + " '" + value + "' "
	} else {
		str += "AND " + key + " " + whe + " '" + value + "' "
	}
	return str
}

func Remove(s []string, i int) ([]string, string) {
	str := s[i]
	s[i] = s[len(s)-1]
	return s[:len(s)-1], str
}

func ToTree(cats []*mode.Category, is bool, name []string) []mode.Tree {
	var trees []mode.Tree
noe:
	for i := range cats {
		tree := mode.Tree{
			Id:       cats[i].CategoryName,
			Title:    cats[i].CategoryName,
			ParentId: cats[i].CategoryUpName,
		}
		if len(cats[i].Children) == 0 {
			tree.Last = true
		} else {
			for i2 := range name {
				if cats[i].CategoryName == name[i2] || cats[i].CategoryName == cats[i].CategoryUpName {
					continue noe
				}
			}
			toTrees := ToTree(cats[i].Children, false, append(name, tree.ParentId))
			if len(toTrees) == 0 {
				tree.Last = true
				tree.Children = nil
			} else {
				tree.Children = toTrees
			}

		}
		trees = append(trees, tree)
	}
	if is {
		var treess []mode.Tree
		for i := range trees {
			if len(trees[i].ParentId) == 0 {
				treess = append(treess, trees[i])
			}
		}
		return treess
	}
	return trees

}

//func ToTree(cats []*mode.Category) []mode.Tree {
//	var trees []*mode.Tree
//	for i := range cats {
//		tree := &mode.Tree{
//			PrId:       cats[i].CategoryName,
//			Title:    cats[i].CategoryName,
//			ParentId: cats[i].CategoryUpName,
//			Last:     true,
//		}
//		trees = append(trees, tree)
//	}
//	for i := range trees {
//		for i2 := range trees {
//			if trees[i2].ParentId == trees[i].PrId {
//				for i3 := range trees {
//					if trees[i3].ParentId == trees[i2].PrId {
//						for i4 := range trees {
//							if trees[i4].ParentId == trees[i3].PrId {
//								for i5 := range trees {
//									if trees[i5].ParentId == trees[i4].PrId {
//										for i6 := range trees {
//											if trees[i6].ParentId == trees[i5].PrId {
//												for i7 := range trees {
//													if trees[i7].ParentId == trees[i6].PrId {
//														trees[i6].Last = false
//														trees[i6].Children = append(trees[i6].Children, *trees[i7])
//														fmt.Println(trees[i6].Title, trees[i7].Title)
//														time.Sleep(1000)
//													}
//												}
//												trees[i5].Last = false
//												trees[i5].Children = append(trees[i5].Children, *trees[i6])
//											}
//										}
//										trees[i4].Last = false
//										trees[i4].Children = append(trees[i4].Children, *trees[i5])
//									}
//								}
//								trees[i3].Last = false
//								trees[i3].Children = append(trees[i3].Children, *trees[i4])
//							}
//						}
//						trees[i2].Last = false
//						trees[i2].Children = append(trees[i2].Children, *trees[i3])
//					}
//				}
//				trees[i].Last = false
//				trees[i].Children = append(trees[i].Children, *trees[i2])
//			}
//		}
//	}
//	var treess []mode.Tree
//	for i := range trees {
//		if len(trees[i].ParentId) == 0 {
//			treess = append(treess, *trees[i])
//		}
//	}
//	return treess
//
//}

// 数组去重
func UniqueArr(arr []string) []string {
	newArr := make([]string, 0)
	tempArr := make(map[string]bool, len(newArr))
	for _, v := range arr {
		if len(v) == 0 {
			continue
		}
		if tempArr[v] == false {
			tempArr[v] = true
			newArr = append(newArr, v)
		}
	}
	tempArr = nil
	return newArr
}

// MergeArray 合并数组
func MergeArray(dest []string, src []string) (result []string) {
	result = make([]string, len(dest)+len(src))
	//将第一个数组传入result
	copy(result, dest)
	//将第二个数组接在尾部，也就是 len(dest):
	copy(result[len(dest):], src)
	return
}

// 数组删除一样的
func UniqueArrT(arr []string, arr2 []string) []string {
	newArr := make([]string, 0)
	for i := range arr {
		is := true
		for i2 := range arr2 {
			if arr[i] == arr2[i2] {
				is = false
				break
			}
		}
		if is {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func WhereREPEAT(str string, num int) []string {
	var strs []string
	for i := 0; i < num; i++ {
		strs = append(strs, str)
	}
	return strs
}

func RemoveEmptyStringsFromArray(arr []string) []string {
	newSlice := make([]string, len(arr))
	j := 0
	for _, v := range arr {
		if v != "" {
			newSlice[j] = v
			j++
		}
	}
	newSlice = newSlice[:j]
	return newSlice
}

func DeleteAtIndex(arr []mode.ProductBrands, index int) []mode.ProductBrands {
	if index < 0 || index >= len(arr) {
		return arr
	}
	return append(arr[:index], arr[index+1:]...)
}
