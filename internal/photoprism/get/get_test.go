package get

import (
	"os"
	"testing"
	"time"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/testextras"
	"github.com/photoprism/photoprism/pkg/dsn"
	"github.com/photoprism/photoprism/pkg/fs"
)

// TestMain executes testMain returning it's results.  It is done this way so that defer can be used to cleanup.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) (code int) {
	// Remove temporary SQLite files before running the tests.
	fs.PurgeTestDbFiles(".", false)

	caller := "internal/photoprism/get/get_test.go/TestMain"
	dbc, dbn, err := testextras.AcquireDBMutex(log, caller)
	if err != nil {
		log.Errorf("testMain: AcquireDBMutex error %+v", err)
		return 1
	}
	defer testextras.UnlockDBMutex(dbc.Db())

	_, dsname := dsn.PhotoPrismTestToDriverDSN(dbn)
	dsn.SetDSNToEnv(dsname)

	tempDir, err := os.MkdirTemp("", "internal-photoprism-get")
	if err != nil {
		panic(err)
	}

	c := config.NewMinimalTestConfigWithDb("test", tempDir)
	code = 999

	defer func() {
		c.CleanupTestFolder()
		if err := c.CloseDb(); err != nil {
			log.Errorf("close db: %v", err)
		}
		// Remove temporary SQLite files after running the tests.
		fs.PurgeTestDbFiles(".", false)
		testextras.ReleaseDBMutex(dbc.Db(), log, caller, code)
		dbc.Close()
		os.RemoveAll(tempDir)
	}()

	SetConfig(c)

	// Run unit tests.
	beforeTimestamp := time.Now().UTC()
	code = m.Run()
	code = testextras.ValidateDBErrors(c.Db(), log, beforeTimestamp, code)

	return code
}
