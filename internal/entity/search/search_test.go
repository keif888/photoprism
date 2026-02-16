package search

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/testextras"
	"github.com/photoprism/photoprism/pkg/dsn"
	"github.com/photoprism/photoprism/pkg/fs"
)

// TestMain executes testMain returning it's results.  It is done this way so that defer can be used to cleanup.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) (code int) {
	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)

	// Remove temporary SQLite files before running the tests.
	fs.PurgeTestDbFiles(".", false)

	caller := "internal/entity/search/search_test.go/TestMain"
	dbc, dbn, err := testextras.AcquireDBMutex(log, caller)
	if err != nil {
		log.Errorf("testMain: AcquireDBMutex error %+v", err)
		return 1
	}
	defer testextras.UnlockDBMutex(dbc.Db())

	if err := testextras.SetupStorage(); err != nil {
		log.Errorf("testMain: SetupStorage error %+v", err)
		return 1
	}
	defer testextras.CleanupStorage()

	driver, dsn := dsn.PhotoPrismTestToDriverDSN(dbn)
	db := entity.InitTestDb(
		driver,
		dsn)

	code = 999

	defer func() {
		// Remove temporary SQLite files after running the tests.
		fs.PurgeTestDbFiles(".", false)
		testextras.ReleaseDBMutex(dbc.Db(), log, caller, code)
		dbc.Close()
	}()

	// Run unit tests.
	beforeTimestamp := time.Now().UTC()
	code = m.Run()
	code = testextras.ValidateDBErrors(db.Db(), log, beforeTimestamp, code)

	return code
}
