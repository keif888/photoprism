package query

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestMain(m *testing.M) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	db := entity.InitTestDb(
		os.Getenv("PHOTOPRISM_TEST_DRIVER"),
		os.Getenv("PHOTOPRISM_TEST_DSN"))

	defer db.Close()

	code := m.Run()

	os.Exit(code)
}

func TestDbDialect(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		if DbDialect() != SQLite3 {
			t.SkipNow()
		}
		assert.Equal(t, "sqlite", DbDialect())
	})

	t.Run("MariaDB", func(t *testing.T) {
		if DbDialect() != MySQL {
			t.SkipNow()
		}
		assert.Equal(t, MySQL, DbDialect())
	})
}

func TestBatchSize(t *testing.T) {
	t.Run("SQLite", func(t *testing.T) {
		if DbDialect() != SQLite3 {
			t.SkipNow()
		}
		assert.Equal(t, 333, BatchSize())
	})
}
