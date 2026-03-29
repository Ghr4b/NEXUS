package test

import (
	"fmt"
	"public_disclosure/models"
	"testing"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Database configuration for testing
	// Using the same credentials as dev environment
	dbUser := "28g7uo7XM1BNj1c.root"
	dbPass := "FxCvYxQAZoanwx5s"
	dbHost := "gateway01.eu-central-1.prod.aws.tidbcloud.com"
	dbPort := 4000
	dbName := "test"

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&tls=true", dbUser, dbPass, dbHost, dbPort, dbName)

	// Register driver and database if not already registered (init is called once per package)
	// To be safe, we check or recover from panic if re-registered?
	// Orm registry acts globally.
	// If running `go test ./...`, other packages might have registered it.
	// But `services` package test is isolated process usually.

	// Check if default is registered? No easy way. Just register and let it panic if duplicate (which means it's already there)
	// Or recover.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from DB registration:", r)
		}
	}()

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", conn)
}

func TestSearchUsers(t *testing.T) {
	o := orm.NewOrm()

	// Create unique prefix for this test run to avoid collision
	timestamp := time.Now().UnixNano()
	prefix := fmt.Sprintf("test_%d_", timestamp)

	// Seed data
	users := []models.User{
		{Username: prefix + "inactive_1", IsActive: false, IsStaff: false, Email: prefix + "john@gov.tn", FirstName: "John", LastName: "Doe"},
		{Username: prefix + "inactive_2", IsActive: false, IsStaff: false, Email: prefix + "jane@gmail.com", FirstName: "Jane", LastName: "Smith"},
		{Username: prefix + "active_staff", IsActive: true, IsStaff: true, Email: prefix + "staff@gov.tn"},
	}

	createdUsers := []*models.User{}
	for _, u := range users {
		// Create copy to insert
		newUser := u
		_, err := o.Insert(&newUser)
		if err != nil {
			t.Fatalf("Failed to insert test user: %v", err)
		}
		createdUsers = append(createdUsers, &newUser)
	}

	// Cleanup function
	defer func() {
		for _, u := range createdUsers {
			o.Delete(u)
		}
	}()

	t.Run("No filters - should find inactive users", func(t *testing.T) {
		res, err := SearchUsers(nil)
		assert.NoError(t, err)

		found := 0
		for _, u := range res {
			if u.Username == prefix+"inactive_1" || u.Username == prefix+"inactive_2" {
				found++
			}
			assert.False(t, u.IsActive)
		}
		assert.Equal(t, 2, found)
	})

	t.Run("Filter by exact username", func(t *testing.T) {
		filters := map[string]interface{}{
			"Username": prefix + "inactive_1",
		}
		res, err := SearchUsers(filters)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, prefix+"inactive_1", res[0].Username)
	})

	t.Run("Filter by icontains email", func(t *testing.T) {
		filters := map[string]interface{}{
			"Email__icontains": "gov.tn",
		}
		res, err := SearchUsers(filters)
		assert.NoError(t, err)

		// This might match other users in DB, so check if our test user is present
		found := false
		for _, u := range res {
			if u.Username == prefix+"inactive_1" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find inactive_1 with email filter")
	})
}
