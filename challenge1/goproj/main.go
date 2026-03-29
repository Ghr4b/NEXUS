package main

import (
	"fmt"
	"net/http" // Added for http.SameSiteNoneMode
	_ "public_disclosure/models"
	_ "public_disclosure/routers"

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

	// Note: ensure your MySQL server supports TLS if using tls=true,
	// otherwise use tls=false or remove it for local dev.
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&tls=true", dbUser, dbPass, dbHost, dbPort, dbName)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", conn)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// ---------------------------------------------------------
	// 1. FORCE SESSION AND COOKIE POLICIES PROGRAMMATICALLY
	// ---------------------------------------------------------
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionCookieHttpOnly = true
	beego.BConfig.WebConfig.Session.SessionCookieSecure = true
	beego.BConfig.WebConfig.Session.SessionSameSite = http.SameSiteNoneMode

	// ---------------------------------------------------------
	// 2. FORCE XSRF POLICIES
	// ---------------------------------------------------------
	beego.BConfig.WebConfig.EnableXSRF = true
	// Note: Beego uses the global request scheme (HTTPS) to determine if
	// the XSRF cookie itself should be marked Secure. Because you are using
	// SameSite=None, you MUST serve this app over HTTPS, which will
	// automatically secure the _xsrf cookie as well.

	// ---------------------------------------------------------
	// 3. ROUTER AND FILTERS
	// ---------------------------------------------------------
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *beecontext.Context) {
		// Prevent site from being embedded in iframes (Anti-Clickjacking)
		ctx.Output.Header("X-Frame-Options", "SAMEORIGIN")

		// Prevent browsers from sniffing MIME types
		ctx.Output.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS filtering
		ctx.Output.Header("X-XSS-Protection", "1; mode=block")

		// HSTS (Strict Transport Security) - Uncomment when on HTTPS!
		// ctx.Output.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	})

	beego.Run()
}
