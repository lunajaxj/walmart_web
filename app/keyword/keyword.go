package keyword

import (
	"fmt"
	"walmart_web/app/config"
	"walmart_web/app/mode"
	"walmart_web/app/tools"
)

func AddKeyword(key mode.Keyword) {
	config.Db.Exec(fmt.Sprintf(`DELETE FROM keyword WHERE name=%s`, key))
	_, err := config.Db.NamedExec(`INSERT INTO keyword (name,ids)
        VALUES (:name,:ids)`, key)
	tools.ErrPr(err, "")
}

func GetKeywordName(name []string) []mode.Keyword {
	var wheres string
	for i := range name {
		wheres = tools.WhereOr(wheres, "=", "name", name[i])
	}

	var keys []mode.Keyword
	err := config.Db.Select(&keys, fmt.Sprintf("SELECT * FROM keyword %s", wheres))
	tools.ErrPr(err, "")
	return keys
}

func GetKeyword() []mode.Keyword {
	var keys []mode.Keyword
	err := config.Db.Select(&keys, "SELECT * FROM keyword")
	tools.ErrPr(err, "")
	return keys
}
