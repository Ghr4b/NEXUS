package controllers

import (
	"encoding/json"
	"html/template"
	"public_disclosure/models"
	"public_disclosure/services"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type StaffManagementController struct {
	beego.Controller
}

func (c *StaffManagementController) Prepare() {
	isStaff := c.GetSession("is_staff")
	isActive := c.GetSession("is_active")
	if isStaff == nil || !isStaff.(bool) || isActive == nil || !isActive.(bool) {
		c.Redirect("/staff/login", 302)
		c.StopRun()
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/staff.tpl"
}

func (c *StaffManagementController) Get() {
	o := orm.NewOrm()
	var departments []*models.Department
	o.QueryTable(new(models.Department)).All(&departments)
	c.Data["Departments"] = departments

	// SECURED: Passed as a standard string.
	// Go's html/template will safely escape this for the JavaScript context in the view.
	filter := c.GetString("filter")
	c.Data["SavedFilter"] = filter

	c.TplName = "staff/management.tpl"
}

func (c *StaffManagementController) Search() {
	var filterMap map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &filterMap)
	users, err := services.SearchUsers(filterMap)
	if err != nil {
		panic(err)
	}
	c.Data["json"] = users
	c.ServeJSON()
}

// NEW: Dedicated method to view and process a single user
func (c *StaffManagementController) ViewUser() {
	userId, err := c.GetInt("id")
	if err != nil {
		c.Redirect("/staff/management", 302)
		return
	}

	o := orm.NewOrm()
	user := models.User{Id: userId}

	if err := o.Read(&user); err != nil {
		c.Redirect("/staff/management", 302)
		return
	}

	var departments []*models.Department
	o.QueryTable(new(models.Department)).All(&departments)

	c.Data["User"] = user
	c.Data["Departments"] = departments
	c.TplName = "staff/user_process.tpl"
}

func (c *StaffManagementController) Approve() {
	userId, _ := c.GetInt("user_id")
	departmentId, _ := c.GetInt("department_id")

	o := orm.NewOrm()
	user := models.User{Id: userId}
	if err := o.Read(&user); err == nil {
		user.IsActive = true
		user.IsStaff = true
		if _, err := o.Update(&user); err == nil {
			staff := models.Staff{
				User:       &models.User{Id: userId},
				Department: &models.Department{Id: departmentId},
			}
			o.Insert(&staff)
		}
	}
	c.Redirect("/staff/management", 302)
}

func (c *StaffManagementController) Reject() {
	userId, _ := c.GetInt("user_id")

	o := orm.NewOrm()
	o.Delete(&models.User{Id: userId})

	c.Redirect("/staff/management", 302)
}
