package provisioner

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/config"
)

func TestCreateDB(t *testing.T) {
	t.Run("MariaDB", func(t *testing.T) {
		c := config.NewTestConfig("cluster")
		// Ensure we're on MariaDB in tests.
		if c.DatabaseDriver() != config.MySQL {
			t.Skip("test requires MariaDB driver in test config")
		}
		creds, created, err := EnsureNodeDatabase(nil, c, "pp-node-01", false)
		assert.Empty(t, err)
		assert.True(t, created)
		assert.Contains(t, creds.DSN, "tcp")

		// Check if connection possible
		db, err := gorm.Open("mysql", creds.DSN)
		defer db.Close()
		if err != nil {
			t.Fatal(err)
		}

		if created {
			if err := c.Db().Unscoped().Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", quoteIdent(creds.Name))).Error; err != nil {
				assert.Empty(t, err)
				t.Logf("Unable to drop database %s", quoteIdent(creds.Name))
			}
			if err := c.Db().Unscoped().Exec(fmt.Sprintf("DROP USER IF EXISTS %s", quoteIdent(creds.User))).Error; err != nil {
				assert.Empty(t, err)
				t.Logf("Unable to drop user %s", quoteIdent(creds.User))
			}
		}
	})
}
