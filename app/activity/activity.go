package activity

import (
	"fmt"
	"strings"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddActivity(act mode.Activity) {
	_, err := config.Db.NamedExec(`INSERT INTO activity (name,date,ids)
        VALUES (:name,:date,:ids)`, act)
	tools.ErrPr(err, "")
}

func GetActivityWhere(name, date string) []mode.Activity {
	var wheres string

	if name != "" {
		split := strings.Split(name, ",")
		wheres = tools.WhereAndOrs(
			wheres,
			tools.WhereREPEAT("=", len(split)),
			tools.WhereREPEAT("name", len(split)),
			split)
	}
	if date != "" {
		split := strings.Split(date, " ~ ")
		wheres = tools.WhereAnd(wheres, ">=", "date", split[0])
		wheres = tools.WhereAnd(wheres, "<=", "date", split[1])
	}
	var acts []mode.Activity
	err := config.Db.Select(&acts, fmt.Sprintf("SELECT * FROM activity %s", wheres))
	tools.ErrPr(err, "")
	return acts
}

func GetActivity() []mode.Activity {
	var acts []mode.Activity
	err := config.Db.Select(&acts, "SELECT * FROM activity")
	tools.ErrPr(err, "")
	return acts
}
