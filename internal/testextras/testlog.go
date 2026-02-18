package testextras

import (
	"os"
	"time"

	"gorm.io/gorm"
)

// TestLog Test Logging structure
type TestLog struct {
	ID        uint      `gorm:"primaryKey;"`
	LogTime   time.Time `sql:"index:idx_testlog_log_time"`
	ProcessID int
	Message   string `gorm:"size:200;default:''"`
}

// LogMessage stores a message into the testextras mutex database
func LogMessage(db *gorm.DB, message string) {
	pid := os.Getpid()
	record := TestLog{LogTime: time.Now().UTC(), ProcessID: pid, Message: message}
	db.Create(&record)
}
