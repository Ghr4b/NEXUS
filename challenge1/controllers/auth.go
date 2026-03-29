package controllers

import (
	"html/template"
	"public_disclosure/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Prepare() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/public.tpl"
}

// Login
// Login
func (c *AuthController) Login() {
	// 1. Grab the 'next' parameter (this is the vulnerable part)
	next := c.GetString("next")
	if next == "" {
		next = "/staff/dashboard" // Default fallback
	}

	// 2. If the user is ALREADY logged in, bounce them immediately
	if c.GetSession("user_id") != nil {
		c.Redirect(next, 302)
		return
	}

	// 3. Handle actual login attempts
	if c.Ctx.Input.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")

		o := orm.NewOrm()
		user := models.User{Username: username}
		err := o.Read(&user, "Username")

		if err == nil {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err == nil {
				c.SetSession("user_id", user.Id)
				c.SetSession("is_staff", user.IsStaff)
				c.SetSession("is_active", user.IsActive)

				// 4. Redirect to 'next' upon successful login
				c.Redirect(next, 302)
				return
			}
		}
		c.Data["Error"] = "Invalid credentials"
	}
	c.TplName = "auth/login.tpl"
}

func (c *AuthController) Logout() {
	c.DelSession("user_id")
	c.DelSession("is_staff")
	c.DelSession("is_active")
	c.Redirect("/", 302)
}

func (c *AuthController) Register() {
	o := orm.NewOrm()
	if c.Ctx.Input.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")
		firstName := c.GetString("first_name")
		lastName := c.GetString("last_name")
		email := c.GetString("email")

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := models.User{
			Username:  username,
			Password:  string(hashedPassword),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			IsStaff:   false,
			IsActive:  false,
		}

		_, err := o.Insert(&user)
		if err == nil {
			// Redirect to login with a message?
			// "Registration successful. Please wait for approval."
			// For now redirecting to login as per request
			c.Redirect("/staff/login", 302)
			return
		}
		c.Data["Error"] = "Registration failed"
	}

	// No need to fetch departments for registration anymore
	c.TplName = "auth/register.tpl"
}
