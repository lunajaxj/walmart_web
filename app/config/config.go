package config

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
	"time"
)

var Rule *casbin.Enforcer
var KeyWg = sync.WaitGroup{}
var IdWg = sync.WaitGroup{}
var UrlWg = sync.WaitGroup{}
var ActWg = sync.WaitGroup{}
var MuxId sync.Mutex
var MuxKey sync.Mutex
var MuxUrl sync.Mutex
var MuxAct sync.Mutex
var MuxDowIds sync.Mutex
var Ch = make(chan int, 10)
var DowIds []string
var DowwIds []string

var ProxyUrl = "l752.kdltps.com:15818"
var ProxyUser = "t19932187800946"
var ProxyPass = "wsad123456"

var Db *sqlx.DB
var KeysFile []string
var KeyssFile []string
var UrlsFile []string
var IdsFile []string

//0:80 本地测试
//1:80 跑产品环境
//2:8080 跑购物车环境
var IsC = false
var IsC2 = true

// 切换环境
var Mode = 0

func init() {
	var db *sqlx.DB
	var err error
	log.Println("连接数据库...")
	if Mode == 0 {
		db, err = sqlx.Open("mysql", "walmart:jHPhbZpM4G7e4NH8@tcp(192.168.2.8:3306)/walmart?charset=utf8mb4")
	} else {
		db, err = sqlx.Open("mysql", "root:disen88888888@tcp(192.168.2.8:3316)/walmart?charset=utf8mb4")

	}

	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	db.SetConnMaxLifetime(120 * time.Minute)
	db.SetMaxIdleConns(50)  // 最大空闲连接数
	db.SetMaxOpenConns(100) // 最大连接数
	Db = db
	if Mode == 1 {
		GetFile()
	}
	//rule()
}
