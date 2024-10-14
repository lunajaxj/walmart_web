package user

import (
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func GetUser(username string) []mode.User {
	var u []mode.User
	strings.Replace(username, "'", "", -1)
	err := config.Db.Select(&u, "select * from user where username = ?", username)
	tools.ErrPr(err, "")

	return u
}
