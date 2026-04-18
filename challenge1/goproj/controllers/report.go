package controllers

import (
	"html/template"
	"public_disclosure/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type ReportController struct {
	beego.Controller
}

func (c *ReportController) Prepare() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/staff.tpl"
}

func (c *ReportController) Get() {
	c.TplName = "staff/report.tpl"
}

func (c *ReportController) Post() {
	url := c.GetString("url")
	content := c.GetString("content")

	if url == "" || content == "" {
		c.Data["Error"] = "All fields are required."
		c.TplName = "staff/report.tpl"
		return
	}

	o := orm.NewOrm()
	report := models.Report{
		Url:     url,
		Content: content,
	}

	if _, err := o.Insert(&report); err == nil {
		c.Data["Success"] = "Report submitted successfully."
	} else {
		c.Data["Error"] = "Failed to submit report: " + err.Error()
	}

	c.TplName = "staff/report.tpl"
}
