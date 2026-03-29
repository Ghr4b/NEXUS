package routers

import (
	"public_disclosure/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Public routes
	beego.Router("/", &controllers.PublicController{}, "get:Get")
	beego.Router("/files/:uuid", &controllers.PublicController{}, "get:GetFile")

	// Auth routes
	beego.Router("/staff/login", &controllers.AuthController{}, "get,post:Login")
	beego.Router("/staff/logout", &controllers.AuthController{}, "get:Logout")
	beego.Router("/staff/register", &controllers.AuthController{}, "get,post:Register")

	// Staff routes
	beego.Router("/staff/dashboard", &controllers.StaffController{}, "get:Dashboard")
	beego.Router("/staff/create", &controllers.StaffController{}, "get,post:CreateFile")
	beego.Router("/staff/files/:uuid", &controllers.StaffController{}, "get:ViewFile")
	beego.Router("/staff/files/:uuid/update", &controllers.StaffController{}, "get,post:UpdateFile")
	beego.Router("/staff/files/:uuid/delete", &controllers.StaffController{}, "post:DeleteFile")
	beego.Router("/staff/upload", &controllers.StaffController{}, "post:UploadAttachment")
	beego.Router("/staff/attachments/:id/delete", &controllers.StaffController{}, "post:DeleteAttachment")

	// Staff Management
	beego.Router("/staff/management", &controllers.StaffManagementController{}, "get:Get")
	beego.Router("/staff/management/search", &controllers.StaffManagementController{}, "post:Search")
	beego.Router("/staff/management/approve", &controllers.StaffManagementController{}, "post:Approve")
	beego.Router("/staff/management/reject", &controllers.StaffManagementController{}, "post:Reject")
	beego.Router("/staff/management/user", &controllers.StaffManagementController{}, "get:ViewUser")

	// Profiles
	beego.Router("/staff/profile", &controllers.StaffController{}, "get,post:Profile")
	beego.Router("/staff/profile/:id", &controllers.StaffController{}, "get:PublicProfile")

	// Audit Logs
	beego.Router("/staff/auditlog", &controllers.AuditController{}, "get:Get")
	beego.Router("/staff/auditlog/:id", &controllers.AuditController{}, "get:GetOne")

	// Reports
	beego.Router("/staff/report", &controllers.ReportController{}, "get:Get;post:Post")

	// Add upload processing route if needed, but the form will validly post to /upload
}
