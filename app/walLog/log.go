package walLog

import (
	"log"
	"walmart_web/app/config"
	"walmart_web/app/mode"
)

func AddLog(lo mode.Log) {
	_, err := config.Db.NamedExec(`INSERT INTO log (classify,msg,val)
        VALUES (:classify,:msg,:val)`, lo)
	if err != nil {
		log.Println("sql错误：", err)
	}
}

func GetLogs() []mode.Log {
	var los []mode.Log
	err := config.Db.Select(&los, "SELECT * FROM log")
	if err != nil {
		log.Println("sql错误：", err)
	}
	return los
}
