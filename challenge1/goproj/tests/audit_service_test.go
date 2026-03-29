package test

import (
	"public_disclosure/models"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Register SQLite for this test file.
	// Since filenames are alphabetical, this might run before user_search_test.go's init.
	// user_search_test.go's init has a recover block, so it should be fine if it fails to register "default" on top of this.
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "file:memory_audit_test?mode=memory&cache=shared")
}

func TestSearchAuditLogs(t *testing.T) {
	// Initialize DB schemas
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		t.Fatalf("Failed to sync db: %v", err)
	}

	o := orm.NewOrm()

	// Seed Data
	// Create Department
	dept := models.Department{Name: "AuditDept"}
	_, err = o.Insert(&dept)
	assert.NoError(t, err)

	// Create User & Staff
	user := models.User{Username: "audit_user", Email: "audit@example.com", IsActive: true}
	_, err = o.Insert(&user)
	assert.NoError(t, err)

	staff := models.Staff{User: &user, Department: &dept}
	_, err = o.Insert(&staff)
	assert.NoError(t, err)

	user2 := models.User{Username: "other_user", Email: "other@example.com", IsActive: true}
	_, err = o.Insert(&user2)
	assert.NoError(t, err)
	staff2 := models.Staff{User: &user2, Department: &dept}
	_, err = o.Insert(&staff2)
	assert.NoError(t, err)

	// Create Logs
	logsData := []models.AuditLog{
		{Staff: &staff, Action: "LOGIN", TargetType: "User", TargetId: user.Id, Message: "User logged in"},
		{Staff: &staff, Action: "VIEW_FILE", TargetType: "File", TargetId: 101, Message: "Viewed file"},
		{Staff: &staff2, Action: "LOGOUT", TargetType: "User", TargetId: user2.Id, Message: "User logged out"},
	}

	for _, l := range logsData {
		_, err := o.Insert(&l)
		assert.NoError(t, err)
	}

	t.Run("Filter by exact Action (Implicit)", func(t *testing.T) {
		filters := map[string]interface{}{
			"Action": "LOGIN",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "LOGIN", res[0].Action)
	})

	t.Run("Filter by exact Action (Explicit Operator) - BUG REPRO", func(t *testing.T) {
		// This test is expected to fail with the current bug.
		// logic: "Action__exact" -> parsed operator "exact" -> reset to "" -> field "Action__exact" -> Not Allowed -> Ignored -> Returns all logs (3)
		filters := map[string]interface{}{
			"Action__exact": "VIEW_FILE",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		// Expectation if working: 1 result.
		// Expectation if bug exists: 3 results (or however many total).
		assert.Len(t, res, 1, "Should filter strictly by Action__exact")
		if len(res) == 1 {
			assert.Equal(t, "VIEW_FILE", res[0].Action)
		}
	})

	t.Run("Filter by Staff Username", func(t *testing.T) {
		filters := map[string]interface{}{
			"Staff__User__Username": "audit_user",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		// Should match LOGIN and VIEW_FILE
		assert.Len(t, res, 2)
		for _, l := range res {
			// We need to load related to check, but ID matches
			assert.Equal(t, staff.Id, l.Staff.Id)
		}
	})

	t.Run("Filter by Department Name", func(t *testing.T) {
		filters := map[string]interface{}{
			"Staff__Department__Name": "AuditDept",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("Filter invalid field", func(t *testing.T) {
		filters := map[string]interface{}{
			"InvalidField": "Something",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		// Should ignore filter and return all
		assert.Len(t, res, 3)
	})

	t.Run("Unsupported operator", func(t *testing.T) {
		// "isnull" is in operators map but NOT in allowedOperators map in compiled code logic (if I recall correctly from view_file)
		// Let's check source code seen earlier.
		// allowedOperators has: icontains, in, gte, lte, gt, lt, exact.
		// operators has: ... startswith, endswith, isnull ...

		filters := map[string]interface{}{
			"Action__startswith": "LOG",
		}
		res, err := SearchAuditLogs(filters)
		assert.NoError(t, err)
		// Should ignore filter and fail to filter? Or just ignore?
		// Code says: "operator ... is not supported, the query will be skipped."
		// So it returns all 3.
		assert.Len(t, res, 3)
	})

}
