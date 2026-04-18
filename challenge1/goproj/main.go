package main

import (
	"fmt"
	_ "public_disclosure/models"
	_ "public_disclosure/routers"

	"net/http"
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

	beego.RunWithMiddleWares(":8111", cookieHardeningMiddleware)
}

type cookieResponseWriter struct {
	http.ResponseWriter
}

func (w *cookieResponseWriter) modifyCookies() {
	cookies := w.Header()["Set-Cookie"]
	for i, c := range cookies {
		isXsrf := strings.Contains(c, "_xsrf=")
		isSession := strings.Contains(c, "sessionid=")

		if isXsrf || isSession {
			// Extract all parts except potential existing SameSite/Secure
			parts := strings.Split(c, ";")
			newParts := []string{}
			hasSecure := false
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if strings.HasPrefix(strings.ToLower(trimmed), "samesite=") {
					continue // Skip existing SameSite
				}
				if strings.ToLower(trimmed) == "secure" {
					hasSecure = true
				}
				newParts = append(newParts, p)
			}

			// Add desired SameSite
			if isXsrf {
				newParts = append(newParts, " SameSite=Strict")
			} else if isSession {
				newParts = append(newParts, " SameSite=None")
			}

			// Add Secure if missing
			if !hasSecure {
				newParts = append(newParts, " Secure")
			}
			cookies[i] = strings.Join(newParts, ";")
		}
	}
}

func (w *cookieResponseWriter) Write(b []byte) (int, error) {
	w.modifyCookies()
	return w.ResponseWriter.Write(b)
}

func (w *cookieResponseWriter) WriteHeader(statusCode int) {
	w.modifyCookies()
	w.ResponseWriter.WriteHeader(statusCode)
}

func cookieHardeningMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &cookieResponseWriter{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
	})
}
