package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestErrors(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		errors, err := Errors(1000, 0, "notexistingErrorString")
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, errors)
	})
	t.Run("Error", func(t *testing.T) {
		errors, err := Errors(1000, 0, "errors")
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, errors)
	})
	t.Run("warning", func(t *testing.T) {
		errors, err := Errors(1000, 0, "warnings")
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, errors)
	})
	t.Run("OneError", func(t *testing.T) {
		expected := "OneError Testing Message"
		if err := Db().Create(&entity.Error{ID: 999999, ErrorTime: time.Now(), ErrorLevel: "debug", ErrorMessage: expected}).Error; err != nil {
			t.Fatal(err)
		}

		errors, err := Errors(1000, 0, "OneError")
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, errors)
		if len(errors) > 0 {
			assert.Equal(t, expected, errors[0].ErrorMessage)
		}

		if err := UnscopedDb().Delete(&entity.Error{ID: 999999}).Error; err != nil {
			t.Fatal(err)
		}
	})

}
