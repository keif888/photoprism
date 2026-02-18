package testextras

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"

	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/dsn"
	pfs "github.com/photoprism/photoprism/pkg/fs"
)

// Stores the number of test databases that are supported.
// The MariaDB scripts/sql/mariadb/reset-testdb.sql and PostgreSQL scripts/sql/postgres/reset-testdb.sql scripts need to create this number of databases.
const dbCount = 8

// dbID holds the database identifier number for this instance
var dbID int

// TestDBChoice structure to assist finding the available databases
type TestDBChoice struct {
	ID uint `gorm:"primaryKey;"`
}

// TestDBMutex structure to store the currently active mutex
type TestDBMutex struct {
	ID          uint      `gorm:"primaryKey;autoIncrement:false"`
	RequestType string    `gorm:"primaryKey;autoIncrement:false;size:50"`
	CreateAt    time.Time `sql:"index:idx_testdbmutex_create_at"`
	ProcessID   int
	Caller      string `gorm:"size:255"`
}

// lockDBMutex Attempts to acquire a database controlled mutex.  Using the table primary key to prevent more than 1 insert succeeding.
// Will retry 60 times with 10s interval, before returning false on failure to get mutex.
// The mutex uses the process id and request type to ensure uniqueness between processes.
func lockDBMutex(db *gorm.DB, requestType, caller string) (ok bool, dbNum int) {
	type Result struct {
		ID uint
	}
	var result Result
	var results []TestDBMutex

	pid := os.Getpid()
	err := errors.New("so i am not nil")
	counter := 0
	ok = false
	for err != nil {
		if len(caller) > 255 {
			caller = caller[:255]
		}

		if err = db.Model(&TestDBChoice{}).Select("test_db_choices.id").Joins("left join test_db_mutexes on test_db_choices.id = test_db_mutexes.id and test_db_mutexes.request_type = ?", requestType).Where("test_db_mutexes.id is null").Order("test_db_choices.id ASC").First(&result).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				LogMessage(db, fmt.Sprintf("%v LockDBMutex No Database Available %v", caller, counter))
				counter++
				time.Sleep(10 * time.Second)

				// Check if any of the stored process id's are no longer active...
				if dberr := db.Model(TestDBMutex{}).Where("id is not null").Find(&results).Error; dberr != nil {
					LogMessage(db, fmt.Sprintf("%v LockDBMutex Find Failed Attempt %v with %s", caller, counter, dberr.Error()))
					return ok, dbNum
				}
				for _, existing := range results {
					if proc, oserr := os.FindProcess(existing.ProcessID); oserr == nil {
						running := false
						if runtime.GOOS == "windows" {
							running = true
							LogMessage(db, fmt.Sprintf("Process %d is running on windows", existing.ProcessID))
						} else {
							if procerr := proc.Signal(syscall.Signal(0)); procerr == nil {
								running = true
								LogMessage(db, fmt.Sprintf("Process %d is running on *nix", existing.ProcessID))
							} else if procerr == os.ErrProcessDone {
								running = false
							} else {
								LogMessage(db, fmt.Sprintf("Unable to Signal %d due to %s", existing.ProcessID, procerr.Error()))
							}
						}
						if !running {
							if dberr := db.Where("process_id = ?", existing.ProcessID).Delete(existing); dberr.Error != nil {
								LogMessage(db, fmt.Sprintf("Unable to delete not running %d due to %s", existing.ProcessID, dberr.Error))
							} else {
								LogMessage(db, fmt.Sprintf("Cleaned up not running process id %d from DBMutex", existing.ProcessID))
							}
						}
					} else {
						LogMessage(db, fmt.Sprintf("Unable to FindProcess %d due to %s", existing.ProcessID, oserr.Error()))
					}
				}
			} else {
				LogMessage(db, fmt.Sprintf("%v LockDBMutex Failed Attempt %v with %s", caller, counter, err.Error()))
				return ok, dbNum
			}
		} else {
			record := TestDBMutex{ID: result.ID, RequestType: requestType, CreateAt: time.Now().UTC(), ProcessID: pid, Caller: caller}
			if err = db.Create(&record).Error; err != nil {
				// Assumption is that this will be a unique index error, because someone else got it before us...
				LogMessage(db, fmt.Sprintf("%v LockDBMutex Failed Attempt %v with %s", caller, counter, err.Error()))
			} else {
				ok = true
				dbNum = int(result.ID)
			}
		}
	}
	return ok, dbNum
}

// UnlockDBMutex deletes the mutex using the processes id.  This should be called with a defer to try and ensure that it always get cleared.
func UnlockDBMutex(db *gorm.DB) {
	pid := os.Getpid()
	record := TestDBMutex{ProcessID: pid}
	db.Where("process_id = ?", pid).Delete(&record)
}

// ReleaseDBMutex clears out a mutex lock and logs messages about it
func ReleaseDBMutex(db *gorm.DB, log event.Logger, caller string, code int) {
	LogMessage(db, fmt.Sprintf("%v UnlockDBMutex", caller))
	UnlockDBMutex(db)
	log.Info("database mutex released")
	LogMessage(db, fmt.Sprintf("%v ending with %v", caller, code))
}

// AcquireDBMutex opens a database connection, and then attempts to acquire a mutex for this process.
func AcquireDBMutex(log event.Logger, caller string) (dbc *DbConn, dbn int, err error) {

	err = nil

	driver, dsn := dsn.PhotoPrismTestToDriverDSN(0)

	// Set default test database driver.
	if driver == "test" || driver == "sqlite" || driver == "" || dsn == "" {
		driver = SQLite3
	}

	dbc, dbn, err = acquireDBMutexCore(log, dsn, dbc, caller, driver, dbn, err)
	dbID = dbn
	return dbc, dbn, err
}

// AcquireMigrationDBMutex opens a database connection, and then attempts to acquire a mutex for this process
// for the purpose of Migration command testing.
func AcquireMigrationDBMutex(log event.Logger, caller string) (dbc *DbConn, dbn int, err error) {

	err = nil
	var dsn string
	return acquireDBMutexCore(log, dsn, dbc, caller, "migration", dbn, err)
}

// acquireDBMutexCore is the core logic to acquiring a database mutex.
// opens a database connection, and then attempts to acquire a mutex for this process
func acquireDBMutexCore(log event.Logger, dsn string, dbc *DbConn, caller string, requestType string, dbn int, err error) (*DbConn, int, error) {
	var dbPath string
	if cwd, err := os.Getwd(); err == nil {
		tmpPaths := strings.SplitAfter(cwd, "photoprism/photoprism")
		if len(tmpPaths) == 2 {
			dbPath = filepath.Join(tmpPaths[0], pfs.StorageDir)
		}
	} else {
		log.Warningf("testextras: Getwd error %s", err.Error())
	}

	dbPath = filepath.Join(dbPath, "testdata")
	dbFile := filepath.Join(dbPath, "unit.mutex.db")
	dsn = fmt.Sprintf("%s?_busy_timeout=5000&_foreign_keys=on", dbFile)
	// Try to create the path, ignoring errors
	_ = os.MkdirAll(dbPath, fs.ModePerm)

	// Create gorm.DB connection provider to SQLite3 as all tests share the same mutex database
	dbc = &DbConn{
		Driver: SQLite3,
		Dsn:    dsn,
	}

	SetDbProvider(dbc)
	log.Info("migrating test extras")
	MigrateTestExtras(dbc.Db().Debug())
	LogMessage(dbc.Db(), fmt.Sprintf("%v starting", caller))
	if ok, n := lockDBMutex(dbc.Db(), requestType, caller); ok {
		LogMessage(dbc.Db(), fmt.Sprintf("%v LockDBMutex database %d acquired", caller, n))
		log.Info("database mutex acquired")
		dbn = n
		dbID = n
	} else {
		log.Error("Unable to get DBMutex")
		err = errors.New("unable to acquire DBMutex")
	}

	return dbc, dbn, err
}

// GetDBMutexID returns the cached database id.
func GetDBMutexID() int {
	return dbID
}
