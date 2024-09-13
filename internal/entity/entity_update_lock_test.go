package entity

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// Function to hold a lock on the photos table for 15 seconds.  Hopefully long enough.
func lock_photos_for_test(t *testing.T, m interface{}, keyNames ...string) {

	t.Logf("lock_photos_for_test Starting")
	// Extract interface slice with all values including zero.
	values, keys, err := ModelValues(m, keyNames...)

	// Has keys and values?
	if err != nil {
		t.Logf("lock_photos_for_test ModelValues Error = %s", err.Error())
		return
	} else if len(keys) != len(keyNames) {
		t.Logf("lock_photos_for_test keys issue")
		return
	}

	db_test := UnscopedDb()
	db_test.Transaction(func(tx *gorm.DB) error {
		if tx.Error != nil {
			t.Logf("lock_photos_for_test Begin = %s", tx.Error.Error())
			return tx.Error
		}

		t.Logf("lock_photos_for_test Before updating the photo")
		result := tx.Model(m).Updates(values)
		t.Logf("lock_photos_for_test After updating the photo")

		if result.Error != nil {
			t.Logf("lock_photos_for_test Model Updates Error = %s", result.Error.Error())
			return result.Error
		}
		
		t.Logf("lock_photos_for_test Sleeping")
		time.Sleep(4 * time.Second)
		t.Logf("lock_photos_for_test Awake, rolling back now")
		tx.Rollback()
		return nil
	})

	t.Logf("lock_photos_for_test Rollback Done")
}

func TestEntityLock(t *testing.T) {
	t.Run("LockedDB", func(t *testing.T) {
	if 1 == 0 {
		m := PhotoFixtures.Pointer("Photo01")
		updatedAt := m.UpdatedAt

		// Should be updated without any issues.
		if err := Update(m, "ID", "PhotoUID"); err != nil {
			assert.Greater(t, m.UpdatedAt.UTC(), updatedAt.UTC())
			t.Fatal(err)
			return
		} else {
			assert.Greater(t, m.UpdatedAt.UTC(), updatedAt.UTC())
			t.Logf("(1) UpdatedAt: %s -> %s", updatedAt.UTC(), m.UpdatedAt.UTC())
			t.Logf("(1) Successfully updated values")
		}
		updatedAt = m.UpdatedAt
		go lock_photos_for_test(t, m, "ID", "PhotoUID")
		// Wait a bit for the other thread to start waiting.
		time.Sleep(time.Second * 1)

		// Should return an error with lock timeout.
		if err := Update(m, "ID", "PhotoUID"); err != nil {
			assert.Greater(t, m.UpdatedAt.UTC(), updatedAt.UTC())
			assert.ErrorContains(t, err, "no such table: photos") // Sql Lite doesn't have wait locking it just throws a missing table message.
			t.Logf("(2) UpdatedAt: %s -> %s", updatedAt.UTC(), m.UpdatedAt.UTC())
			t.Logf("(2) Error was %s", err.Error())
			// Wait for the rollback to be done.
			time.Sleep(time.Second * 8)
			// Find the record that was there before the lock
			n := m.Find()
			t.Logf("(3) UpdatedAt: %s -> %s", updatedAt.UTC(), n.UpdatedAt.UTC())
			assert.Equal(t, n.UpdatedAt.UTC(), updatedAt.UTC())
			return
		} else {
			t.Logf("(2) Error not found")
			t.Fail()
		}
	} else {
		t.Logf("Test Skipped as a way to run this test in a mutually exclusive way hasn't been found.")
		t.SkipNow()
		return
	}
	})

}