package commands

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
)

func TestConfigShutdown(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		ctx := config.CliTestContext()
		conf, err := InitConfig(ctx)

		if err != nil {
			t.Errorf("unexpected error %s", err)
		}

		conf.InitDb()

		var beforeCount int64
		conf.Db().Model(&entity.Error{}).Count(&beforeCount)

		event.Publish(
			"log.Warn",
			event.Data{
				"time":    time.Now().UTC().Truncate(time.Second),
				"level":   "Warn",
				"message": "testing: Push a warning record before Shutdown",
			},
		)

		var count int64

		// Wait up to 15 seconds so, hopefully, the event logger will have written to the database.
		count = beforeCount
		for i := 0; i < 15; i++ {
			time.Sleep(time.Second * 1)
			conf.Db().Model(&entity.Error{}).Count(&count)
			if count > beforeCount {
				break
			}
		}

		assert.Greater(t, count, beforeCount, "The event hasn't been written to the database")

		conf.Shutdown()

		// If Shutdown has not stopped the event logger, this will cause a Fatal error and stop the test.
		event.Publish(
			"log.Warn",
			event.Data{
				"time":    time.Now().UTC().Truncate(time.Second),
				"level":   "Warn",
				"message": "testing: Push a warning record after Shutdown",
			},
		)

		// Allow 5 seconds for the published event to propogate
		time.Sleep(time.Second * 5)
		t.Log("no fatal error has caused TestConfigShutdown to fail, there should be at least 1 logevents message(s) before this")
	})

	t.Run("SecondAttempt", func(t *testing.T) {
		ctx := config.CliTestContext()
		conf, err := InitConfig(ctx)

		if err != nil {
			t.Errorf("unexpected error %s", err)
		}

		conf.InitDb()

		var beforeCount int64
		conf.Db().Model(&entity.Error{}).Count(&beforeCount)

		event.Publish(
			"log.Warn",
			event.Data{
				"time":    time.Now().UTC().Truncate(time.Second),
				"level":   "Warn",
				"message": "testing: Push a 2nd warning record before Shutdown",
			},
		)

		var count int64

		// Wait up to 15 seconds so, hopefully, the event logger will have written to the database.
		count = beforeCount
		for i := 0; i < 15; i++ {
			time.Sleep(time.Second * 1)
			conf.Db().Model(&entity.Error{}).Count(&count)
			if count > beforeCount {
				break
			}
		}

		assert.Greater(t, count, beforeCount, "The event hasn't been written to the database")

		conf.Shutdown()

		// If Shutdown has not stopped the event logger, this will cause a Fatal error and stop the test.
		event.Publish(
			"log.Warn",
			event.Data{
				"time":    time.Now().UTC().Truncate(time.Second),
				"level":   "Warn",
				"message": "testing: Push a 2nd warning record after Shutdown",
			},
		)

		// Allow 5 seconds for the published event to propogate
		time.Sleep(time.Second * 5)
		t.Log("no fatal error has caused TestConfigShutdown to fail, there should be at least 1 logevents message(s) before this")
	})

}
