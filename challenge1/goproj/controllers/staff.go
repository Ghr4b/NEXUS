package controllers

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"public_disclosure/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StaffController struct {
	beego.Controller
}

func (c *StaffController) Prepare() {
	isStaff := c.GetSession("is_staff")
	if isStaff == nil || !isStaff.(bool) {
		c.Redirect("/staff/login", 302)
		c.StopRun()
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Layout = "layout/staff.tpl"
}

func (c *StaffController) Dashboard() {
	o := orm.NewOrm()
	var files []*models.DisclosureFile
	o.QueryTable(new(models.DisclosureFile)).OrderBy("-CreatedAt").All(&files)
	c.Data["Files"] = files
	c.TplName = "staff/dashboard.tpl"
}

func (c *StaffController) CreateFile() {
	if c.Ctx.Input.IsPost() {
		title := c.GetString("title")
		description := c.GetString("description")
		isPublished, _ := c.GetBool("is_published")

		fileUUID := uuid.New().String()

		// Get current user (Staff)
		userId := c.GetSession("user_id").(int)
		var staff models.Staff
		o := orm.NewOrm()
		o.QueryTable(new(models.Staff)).Filter("User__Id", userId).One(&staff)

		disclosureFile := models.DisclosureFile{
			Uuid:        fileUUID,
			Title:       title,
			Description: description,
			IsPublished: isPublished,
			CreatedBy:   &staff,
		}

		if _, err := o.Insert(&disclosureFile); err == nil {
			// Log audit
			audit := models.AuditLog{
				Staff:      &staff,
				Action:     "Created Case",
				TargetType: "DisclosureFile",
				TargetId:   disclosureFile.Id,
				Message:    fmt.Sprintf("Created file %s", title),
			}
			o.Insert(&audit)

			c.Redirect("/staff/dashboard", 302)
		} else {
			c.Data["Error"] = "Failed to create file: " + err.Error()
		}
	}
	c.TplName = "staff/create_file.tpl"
}

func (c *StaffController) UploadAttachment() {
	fileUUID := c.GetString("uuid")
	// Since uuid is likely passed as hidden field
	if fileUUID == "" {
		fileUUID = c.Ctx.Input.Param(":uuid")
	}

	f, h, err := c.GetFile("attachment")
	if err == nil {
		fileName := h.Filename
		// Save file
		uploadDir := "static/uploads/"
		filePath := filepath.Join(uploadDir, fileName)
		c.SaveToFile("attachment", filePath)

		// Calculate SHA256
		f.Seek(0, 0)
		hasher := sha256.New()
		if _, err := io.Copy(hasher, f); err == nil {
			hash := fmt.Sprintf("%x", hasher.Sum(nil))

			// Create Attachment record
			attachment := models.Attachment{
				FileName:   fileName,
				FilePath:   filePath,
				FileSize:   h.Size,
				Sha256Hash: hash,
			}

			o := orm.NewOrm()
			if id, err := o.Insert(&attachment); err == nil {
				// Link to File
				disclosureFile := models.DisclosureFile{Uuid: fileUUID}
				o.Read(&disclosureFile, "Uuid")

				m2m := o.QueryM2M(&disclosureFile, "Attachments")
				m2m.Add(&attachment)

				// Audit log
				userId := c.GetSession("user_id").(int)
				var staff models.Staff
				o.QueryTable(new(models.Staff)).Filter("User__Id", userId).One(&staff)

				audit := models.AuditLog{
					Staff:      &staff,
					Action:     "Uploaded Attachment",
					TargetType: "Attachment",
					TargetId:   int(id),
					Message:    fmt.Sprintf("Uploaded %s to case %s", fileName, fileUUID),
				}
				o.Insert(&audit)
			}
		}
		f.Close()
	}
	c.Redirect("/staff/files/"+fileUUID, 302)
}

func (c *StaffController) ViewFile() {
	uuid := c.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	file := models.DisclosureFile{Uuid: uuid}
	if err := o.Read(&file, "Uuid"); err == nil {
		o.LoadRelated(&file, "Attachments")
		c.Data["File"] = file
		c.TplName = "staff/view_file.tpl"
	} else {
		c.Redirect("/staff/dashboard", 302)
	}
}

func (c *StaffController) UpdateFile() {
	uuid := c.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	file := models.DisclosureFile{Uuid: uuid}

	if err := o.Read(&file, "Uuid"); err != nil {
		c.Redirect("/staff/dashboard", 302)
		return
	}

	if c.Ctx.Input.IsPost() {
		title := c.GetString("title")
		description := c.GetString("description")
		isPublished, _ := c.GetBool("is_published")

		file.Title = title
		file.Description = description
		file.IsPublished = isPublished

		if _, err := o.Update(&file); err == nil {
			// Audit log
			userId := c.GetSession("user_id").(int)
			var staff models.Staff
			o.QueryTable(new(models.Staff)).Filter("User__Id", userId).One(&staff)

			audit := models.AuditLog{
				Staff:      &staff,
				Action:     "Updated Case",
				TargetType: "DisclosureFile",
				TargetId:   file.Id,
				Message:    fmt.Sprintf("Updated file %s details", title),
			}
			o.Insert(&audit)

			c.Redirect("/staff/files/"+uuid, 302)
		} else {
			c.Data["Error"] = "Failed to update file"
		}
	}

	c.Data["File"] = file
	c.TplName = "staff/edit_file.tpl"
}

func (c *StaffController) DeleteFile() {
	uuid := c.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	file := models.DisclosureFile{Uuid: uuid}

	if err := o.Read(&file, "Uuid"); err == nil {
		// Load attachments to delete them first?
		// Or let DB handle cascade if configured, but ORM might need manual handling for M2M
		// For now, let's keep attachments orphan or delete M2M relation.

		// Audit log before deletion (so we have record of what was deleted)
		userId := c.GetSession("user_id").(int)
		var staff models.Staff
		o.QueryTable(new(models.Staff)).Filter("User__Id", userId).One(&staff)

		audit := models.AuditLog{
			Staff:      &staff,
			Action:     "Deleted Case",
			TargetType: "DisclosureFile",
			TargetId:   file.Id,
			Message:    fmt.Sprintf("Deleted file %s (%s)", file.Title, file.Uuid),
		}
		o.Insert(&audit)

		// Delete M2M relations is automatic in some ORMs or needs explicit clear
		m2m := o.QueryM2M(&file, "Attachments")
		m2m.Clear()

		o.Delete(&file)
	}
	c.Redirect("/staff/dashboard", 302)
}

func (c *StaffController) DeleteAttachment() {
	id, _ := c.GetInt(":id")
	o := orm.NewOrm()
	attachment := models.Attachment{Id: id}

	if err := o.Read(&attachment); err == nil {
		// Audit log
		userId := c.GetSession("user_id").(int)
		var staff models.Staff
		o.QueryTable(new(models.Staff)).Filter("User__Id", userId).One(&staff)

		audit := models.AuditLog{
			Staff:      &staff,
			Action:     "Deleted Attachment",
			TargetType: "Attachment",
			TargetId:   attachment.Id,
			Message:    fmt.Sprintf("Deleted attachment %s", attachment.FileName),
		}
		o.Insert(&audit)

		// Remove file from disk?
		// os.Remove(attachment.FilePath)

		o.Delete(&attachment)

		// Redirect back to file view. We need to know which file it belonged to.
		// Since it's M2M, it could belong to multiple.
		// Check Referer header or passed query param?
		// Ideally we should pass ?file_uuid=...
		fileUuid := c.GetString("file_uuid")
		if fileUuid != "" {
			c.Redirect("/staff/files/"+fileUuid, 302)
			return
		}
	}
	c.Redirect("/staff/dashboard", 302)
}

func (c *StaffController) Profile() {
	userId := c.GetSession("user_id").(int)
	o := orm.NewOrm()
	user := models.User{Id: userId}

	if c.Ctx.Input.IsPost() {
		// Update profile
		if err := o.Read(&user); err == nil {
			user.FirstName = c.GetString("first_name")
			user.LastName = c.GetString("last_name")
			// user.Email = c.GetString("email") // Unique constraint might fail if changed to existing
			// For simplicity let's allow trying to update it
			email := c.GetString("email")
			if email != "" {
				user.Email = email
			}

			// Password change logic (optional for now)
			newPass := c.GetString("password")
			if newPass != "" {
				// hash and set
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
				user.Password = string(hashedPassword)
			}

			if _, err := o.Update(&user); err == nil {
				c.Data["Success"] = "Profile updated successfully"
			} else {
				c.Data["Error"] = "Failed to update profile: " + err.Error()
			}
		}
	} else {
		o.Read(&user)
	}

	// Load staff info
	var staff models.Staff
	o.QueryTable(new(models.Staff)).Filter("User__Id", userId).RelatedSel().One(&staff)
	c.Data["User"] = user
	c.Data["Staff"] = staff

	c.TplName = "staff/profile.tpl"
}

func (c *StaffController) PublicProfile() {
	id, _ := c.GetInt(":id")
	o := orm.NewOrm()

	// Find staff by user id
	var staff models.Staff
	if err := o.QueryTable(new(models.Staff)).Filter("User__Id", id).RelatedSel().One(&staff); err == nil {
		c.Data["Staff"] = staff
		c.TplName = "staff/public_profile.tpl"
	} else {
		c.Redirect("/staff/dashboard", 302)
	}
}
