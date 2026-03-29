package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id          int    `orm:"auto"`
	Username    string `orm:"unique;size(100)"`
	FirstName   string `orm:"size(100)"`
	LastName    string `orm:"size(100)"`
	Email       string `orm:"unique;size(100)"`
	Password    string `orm:"size(100)" filter:"false"`
	IsStaff     bool   `orm:"default(false)"`
	IsActive    bool   `orm:"default(true)"`
	IsSuperuser bool   `orm:"default(false)"`
	Staff       *Staff `orm:"reverse(one)"`
}

type Department struct {
	Id    int      `orm:"auto"`
	Name  string   `orm:"unique;size(100)"`
	Staff []*Staff `orm:"reverse(many)"`
}

type Staff struct {
	Id           int               `orm:"auto"`
	User         *User             `orm:"rel(one)"`
	Department   *Department       `orm:"rel(fk)"`
	CreatedFiles []*DisclosureFile `orm:"reverse(many)"`
}

type DisclosureFile struct {
	Id          int           `orm:"auto"`
	Uuid        string        `orm:"unique;index;size(36)"`
	Title       string        `orm:"size(255)"`
	Description string        `orm:"type(text)"`
	CreatedBy   *Staff        `orm:"rel(fk)"`
	CreatedAt   time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time     `orm:"auto_now;type(datetime)"`
	IsPublished bool          `orm:"default(false)"`
	Attachments []*Attachment `orm:"rel(m2m)"`
}

type Attachment struct {
	Id              int    `orm:"auto"`
	FileName        string `orm:"size(255)"`
	FilePath        string `orm:"size(512)"`
	FileSize        int64
	Sha256Hash      string            `orm:"size(64)"`
	UploadedAt      time.Time         `orm:"auto_now_add;type(datetime)"`
	DisclosureFiles []*DisclosureFile `orm:"reverse(many)"`
}

type AuditLog struct {
	Id         int    `orm:"auto"`
	Staff      *Staff `orm:"rel(fk)"`
	Action     string `orm:"size(255)"`
	TargetType string `orm:"size(50)"`
	TargetId   int
	Message    string    `orm:"type(text)"`
	Timestamp  time.Time `orm:"auto_now_add;type(datetime)"`
}

type Report struct {
	Id        int       `orm:"auto"`
	Url       string    `orm:"size(512)"`
	Content   string    `orm:"type(text)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User), new(Department), new(Staff), new(DisclosureFile), new(Attachment), new(AuditLog), new(Report))
}
