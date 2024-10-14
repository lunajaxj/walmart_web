package mail

import (
	"fmt"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddMail(ma mode.Mail) {
	_, err := config.Db.NamedExec(`INSERT INTO mail (count,seller,msg,update_time)
       									VALUES (:count,:seller,:msg:update_time)`, ma)
	tools.ErrPr(err, "")
}

func GetMails(page, limit int) ([]mode.Mail, int) {
	var mas []mode.Mail
	var counts []int
	page = (page - 1) * limit
	err := config.Db.Select(&mas, fmt.Sprintf("SELECT * FROM mail LIMIT %d,%d", page, limit))
	tools.ErrPr(err, "")
	err = config.Db.Select(&counts, fmt.Sprintf("SELECT count(*) FROM mail"))
	tools.ErrPr(err, "")
	return mas, counts[0]
}
func GetMail(id string) mode.Mail {
	var ma []mode.Mail
	err := config.Db.Select(&ma, fmt.Sprintf("SELECT * FROM mail where m_id ='%s'", id))
	tools.ErrPr(err, "")
	if len(ma) > 0 {
		return ma[0]
	}
	return mode.Mail{}

}

func UploadMail(tu mode.Mail) int {
	stmt := "UPDATE mail set count=:count,msg=:msg,seller=:seller  WHERE m_id=:m_id"
	affected, err := config.Db.NamedExec(stmt, map[string]interface{}{
		"m_id":   tu.MId,
		"count":  tu.Count,
		"msg":    tu.Msg,
		"seller": tu.Seller,
	})

	tools.ErrPr(err, "")
	if err == nil {
		rowsAffected, _ := affected.RowsAffected()
		return int(rowsAffected)
	}
	return 0

}

func Remove(ids string) int {
	split := strings.Split(ids, ",")
	var wheres string
	if len(split) == 0 {
		return 0
	}
	for i := range split {
		wheres = tools.WhereOrInt(wheres, "=", "m_id", split[i])
	}
	updateResult := config.Db.MustExec(fmt.Sprintf("DELETE FROM mail %s", wheres))
	affected, err := updateResult.RowsAffected()
	tools.ErrPr(err, "")
	return int(affected)
}

func DelMsg() {
	config.Db.MustExec("UPDATE mail SET msg = '',count = 0")

}
