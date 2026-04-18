package main

import (
	"fmt"
	_ "public_disclosure/models"
	_ "public_disclosure/routers"

	"strings"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	beecontext "github.com/beego/beego/v2/server/web/context"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbUser, _ := beego.AppConfig.String("db_user")
	dbPass, _ := beego.AppConfig.String("db_pass")
	dbHost, _ := beego.AppConfig.String("db_host")
	dbPort, _ := beego.AppConfig.Int("db_port")
	dbName, _ := beego.AppConfig.String("db_name")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&tls=true", dbUser, dbPass, dbHost, dbPort, dbName)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", conn)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *beecontext.Context) {
		ctx.Output.Header("X-Frame-Options", "SAMEORIGIN")
		ctx.Output.Header("X-Content-Type-Options", "nosniff")

		ctx.Output.Header("X-XSS-Protection", "1; mode=block")

	})

	beego.InsertFilter("*", beego.FinishRouter, func(ctx *beecontext.Context) {
		cookies := ctx.ResponseWriter.Header()["Set-Cookie"]
		for i, c := range cookies {
			isXsrf := strings.Contains(c, "_xsrf=")
			isSession := strings.Contains(c, "sessionid=")

			if isXsrf || isSession {
				newCookie := c
				if !strings.Contains(c, "SameSite=") {
					if isXsrf {
						newCookie += "; SameSite=Strict"
					} else if isSession {
						newCookie += "; SameSite=None"
					}
				}
				if !strings.Contains(c, "Secure") {
					newCookie += "; Secure"
				}
				cookies[i] = newCookie
			}
		}
	})

	beego.Run()
}
