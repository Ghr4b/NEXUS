package controllers

import (
	"fmt"
	"html/template"
	"public_disclosure/models"
	"public_disclosure/services"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type AuditController struct {
	beego.Controller
}

func (c *AuditController) Prepare() {
	isStaff := c.GetSession("is_staff")
	isActive := c.GetSession("is_active")
	if isStaff == nil || !isStaff.(bool) || isActive == nil || !isActive.(bool) {
		c.Redirect("/staff/login", 302)
		c.StopRun()
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/staff.tpl"
}

func (c *AuditController) GetOne() {
	id, _ := c.GetInt(":id")
	o := orm.NewOrm()
	audit := models.AuditLog{Id: id}

	if err := o.Read(&audit); err == nil {
		o.LoadRelated(&audit, "Staff")
		if audit.Staff != nil {
			o.LoadRelated(audit.Staff, "User")
			o.LoadRelated(audit.Staff, "Department")
		}
		c.Data["Audit"] = audit
		c.TplName = "staff/audit_detail.tpl"
	} else {
		c.Data["Error"] = "Audit log not found"
		c.TplName = "staff/audit_detail.tpl"
	}
}

func (c *AuditController) Get() {
	queryParams := c.Ctx.Input.Context.Request.URL.Query()
	filterMap := make(map[string]interface{})
	hasFilter := false

	for key, values := range queryParams {
		if len(values) > 0 && values[0] != "" {
			filterMap[key] = values[0]
			hasFilter = true
		}
	}

	logs, err := services.SearchAuditLogs(filterMap)
	if err == nil {
		if len(logs) == 1 && hasFilter {
			c.Redirect(fmt.Sprintf("/staff/auditlog/%d", logs[0].Id), 302)
			return
		}

		o := orm.NewOrm()
		for _, log := range logs {
			o.LoadRelated(log, "Staff")
			if log.Staff != nil {
				o.LoadRelated(log.Staff, "User")
			}
		}
		c.Data["Logs"] = logs
	}

	c.TplName = "staff/audit_list.tpl"
}
