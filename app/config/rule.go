package config

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func rule() {
	log.Println("加载权限配置...")
	var err error
	Rule, err = casbin.NewEnforcer("./config/rbac_models.conf", "./config/rbac.csv")
	if err != nil {
		fmt.Println("rule err:", err)
		return
	}
	//Rule.AddPolicy("1000", "/login/index", "GET")
	//Rule.AddPolicy("1", "/admin/shopping/*", "GET")

}

func RuleConfig(c *gin.Context) {
	session := sessions.Default(c)
	re := session.Get("rule")

	obj := c.Request.URL.Path

	act := c.Request.Method
	if re == nil && !strings.Contains(obj, "login") && !strings.Contains(obj, "public") {
		c.Redirect(http.StatusFound, "/login/index")
		c.Abort()
	}
	if re == nil {
		re = "1000"
	}
	if obj == "/" {
		c.Redirect(http.StatusFound, "/admin/index")
		c.Abort()
	}
	enforce, _ := Rule.Enforce(re, obj, act)
	if enforce {
		c.Next()
	} else {
		c.HTML(403, "403.html", nil)
		c.Abort()
	}
}
