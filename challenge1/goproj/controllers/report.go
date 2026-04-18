package controllers

import (
	"public_disclosure/services"

	beego "github.com/beego/beego/v2/server/web"
)

type ReportController struct {
	beego.Controller
}

func (c *ReportController) Prepare() {
	c.Layout = "layout/staff.tpl"
}

func (c *ReportController) Get() {
	c.TplName = "staff/report.tpl"
}

func (c *ReportController) Post() {
	uri := c.GetString("uri")
	content := c.GetString("content")

	if uri == "" || content == "" {
		c.Data["Error"] = "All fields are required."
		c.TplName = "staff/report.tpl"
		return
	}
	err := services.VisitLocalUri(uri)
	if err != nil {
		c.Data["Error"] = "Failed to submit report: " + err.Error()
		c.TplName = "staff/report.tpl"
		return
	}

	c.Data["Success"] = "Report submitted successfully!"

	// Optional: Clear fields so the form is empty on reload
	c.Data["uri"] = ""
	c.Data["content"] = ""

	c.TplName = "staff/report.tpl"
}
