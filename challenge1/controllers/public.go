package controllers

import (
	"html/template"
	"public_disclosure/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type PublicController struct {
	beego.Controller
}

func (c *PublicController) Prepare() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/public.tpl"
}

// Get handles the homepage listing and search
func (c *PublicController) Get() {
	search := c.GetString("search")
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.DisclosureFile))

	// 1. Always create a base condition for global rules (like IsPublished)
	mainCond := orm.NewCondition().And("IsPublished", true)

	if search != "" {
		// 2. Create a separate condition for the search OR logic
		searchCond := orm.NewCondition().
			Or("Title__icontains", search).
			Or("Description__icontains", search).
			Or("Uuid__icontains", search)

		// 3. Combine them: (IsPublished == true) AND (Title OR Description OR Uuid)
		mainCond = mainCond.AndCond(searchCond)
	}

	// 4. Apply the final combined condition to the QuerySeter
	qs = qs.SetCond(mainCond)

	var files []*models.DisclosureFile
	qs.OrderBy("-CreatedAt").All(&files)

	c.Data["Files"] = files
	c.Data["Search"] = search
	c.TplName = "public/index.tpl"
}

// GetFile handles the file detail view
func (c *PublicController) GetFile() {
	uuid := c.Ctx.Input.Param(":uuid")

	o := orm.NewOrm()
	file := models.DisclosureFile{Uuid: uuid}

	err := o.Read(&file, "Uuid")
	if err == orm.ErrNoRows {
		c.Abort("404")
		return
	}

	if !file.IsPublished {
		c.Abort("404")
		return
	}

	o.LoadRelated(&file, "Attachments")

	c.Data["File"] = file
	c.TplName = "public/detail.tpl"
}
