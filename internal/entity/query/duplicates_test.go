package query

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestDuplicates(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		if files, err := Duplicates(10, 0, ""); err != nil {
			t.Fatal(err)
		} else if files == nil {
			t.Fatal("files must not be nil")
		}
	})
	t.Run("pathname not empty", func(t *testing.T) {
		files, err := Duplicates(10, 0, "/holiday/sea.jpg")

		if err != nil {
			t.Fatal(err)
		}

		assert.Empty(t, files)
	})
	t.Run("University", func(t *testing.T) {
		entity.AddDuplicate("education/university/BSc-Thesis.pdf", "/", "abcdefg", int64(555), int64(1234))
		entity.AddDuplicate("education/university/BSc-Thesis.pdf", "backup", "abcdefg", int64(555), int64(1234))
		if files, err := Duplicates(10, 0, "education/university"); err != nil {
			t.Fatal(err)
		} else if files == nil {
			t.Fatal("files must not be nil")
		} else {
			assert.Equal(t, 2, len(files))
		}
		if err := UnscopedDb().Where("1=1").Delete(&entity.Duplicate{}).Error; err != nil {
			t.Fatal(err)
		}
	})
}
