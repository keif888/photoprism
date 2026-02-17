package frame

import (
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/http/header"

	"github.com/stretchr/testify/assert"
)

func TestPolaroid(t *testing.T) {
	dir := t.TempDir()
	t.Run("RandomAngle", func(t *testing.T) {
		img, err := imaging.Open("testdata/500x500.jpg")
		assert.NoError(t, err)

		saveName := filepath.Join(dir, "test-polaroid.png")

		out, err := polaroid(img, RandomAngle(30))

		assert.NoError(t, err)

		err = imaging.Save(out, saveName)

		assert.NoError(t, err)
		mimeType, _ := fs.DetectMimeType(saveName)
		assert.Equal(t, header.ContentTypePng, mimeType)
	})
}
