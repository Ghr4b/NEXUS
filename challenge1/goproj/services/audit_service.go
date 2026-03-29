package services

import (
	"public_disclosure/models"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type FieldCategory int

var (
	exactMatchRegex      = regexp.MustCompile(`^(?:TargetId|Action)$`)
	staffDepartmentRegex = regexp.MustCompile(`^Staff__Department__(?:Name|Id)`)
	staffUserRegex       = regexp.MustCompile(`^Staff__User__(?:Id|Username|Email)$`)
)
var allowedOperators = map[string]struct{}{
	"icontains":  {},
	"in":         {},
	"gte":        {},
	"lte":        {},
	"gt":         {},
	"lt":         {},
	"exact":      {},
	"startswith": {},
}
var operators = map[string]struct{}{
	"exact":       {},
	"iexact":      {},
	"contains":    {},
	"icontains":   {},
	"gt":          {},
	"gte":         {},
	"lt":          {},
	"lte":         {},
	"in":          {},
	"startswith":  {},
	"istartswith": {},
	"endswith":    {},
	"iendswith":   {},
	"isnull":      {},
}

func SearchAuditLogs(filterMap map[string]interface{}) ([]*models.AuditLog, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.AuditLog))

	for key, value := range filterMap {
		operator := ""
		parts := strings.Split(key, "__")
		if len(parts) >= 2 {
			operator = parts[len(parts)-1]

			if _, ok := operators[operator]; ok {

				if _, ok := allowedOperators[operator]; !ok {
					logs.Warning("the operator '%s' in query parameter '%s' is not supported, the query will be skipped.", operator, key)
					continue
				}
			}
			operator = ""
		}
		field := strings.TrimSuffix(key, operator)
		if !AllowedFields(field) {
			logs.Warning("the field '%s' is not allowed, the query will be skipped.", key)
			continue
		}
		qs = qs.Filter(key, toArgs(value)...)
	}

	qs = qs.OrderBy("-Timestamp")

	var logs []*models.AuditLog
	_, err := qs.All(&logs)
	return logs, err
}

// AllowedFields returns true if the field is allowed
func AllowedFields(field string) bool {
	if staffDepartmentRegex.MatchString(field) {
		return true
	}
	if staffUserRegex.MatchString(field) {
		return true
	}
	if exactMatchRegex.MatchString(field) {
		return true
	}
	return false
}
