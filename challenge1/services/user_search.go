package services

import (
	"public_disclosure/models"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func SearchUsers(filterMap map[string]interface{}) ([]*models.User, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.User)).Filter("IsActive", false)

	for key, value := range filterMap {
		// Split key into field and operator
		keyPieces := strings.Split(key, "__")
		if len(keyPieces) > 2 {
			logs.Warning("the key '%s' is not allowed, the query will be skipped.", key)
			continue
		}

		field := keyPieces[0]

		// Check if field is allowed
		allowedFields := map[string]struct{}{
			"Username":  {},
			"FirstName": {},
			"LastName":  {},
			"Email":     {},
		}

		if _, ok := allowedFields[field]; !ok {
			logs.Warning("the field '%s' is not allowed, the query will be skipped.", field)
			continue
		}

		// Check operator if provided
		if len(keyPieces) == 2 {
			operator := orm.ExprSep + keyPieces[1]
			allowedOperators := map[string]struct{}{
				"__icontains": {},
				"__in":        {},
				"__gte":       {},
				"__lte":       {},
				"__gt":        {},
				"__lt":        {},
				"__exact":     {},
			}
			if _, ok := allowedOperators[operator]; !ok {
				logs.Warning("the operator '%s' in query parameter '%s' is not supported, the query will be skipped.", operator, key)
				continue
			}
			// Reconstruct key with proper operator prefix
			key = field + operator
		}

		qs = qs.Filter(key, toArgs(value)...)
	}

	var users []*models.User
	_, err := qs.All(&users)
	return users, err
}

func toArgs(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	// Check if it's a slice/array
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		n := rv.Len()
		if n == 0 {
			return nil
		}
		args := make([]interface{}, n)
		for i := 0; i < n; i++ {
			args[i] = rv.Index(i).Interface()
		}
		return args
	}

	// Single value - wrap in slice
	return []interface{}{v}
}
