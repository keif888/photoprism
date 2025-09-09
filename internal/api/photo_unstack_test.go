package api

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/fs"
)

func TestPhotoUnstack(t *testing.T) {
	t.Run("unstack xmp sidecar file", func(t *testing.T) {
		app, router, _ := NewApiTest()
		PhotoUnstack(router)
		r := PerformRequest(app, "POST", "/api/v1/photos/ps6sg6be2lvl0yh7/files/fs6sg6bw45bnlqdw/unstack")
		// Sidecar files can not be unstacked.
		assert.Equal(t, http.StatusBadRequest, r.Code)
		// t.Logf("RESP: %s", r.Body.String())
	})

	t.Run("unstack bridge3.jpg", func(t *testing.T) {
		app, router, _ := NewApiTest()

		// Create test files so unstack can function
		c := get.Config()
		os.RemoveAll(filepath.Join(c.OriginalsPath(), "1990"))
		os.RemoveAll(filepath.Join(c.OriginalsPath(), "London"))
		fs.Copy("testdata/london_160x160.jpg", filepath.Join(c.OriginalsPath(), "1990", "04", "bridge2.jpg"))
		fs.Copy("testdata/london_160x160.jpg", filepath.Join(c.OriginalsPath(), "London", "bridge3.jpg"))

		PhotoUnstack(router)
		r := PerformRequest(app, "POST", "/api/v1/photos/ps6sg6be2lvl0yh7/files/fs6sg6bwhhbnlqdn/unstack")
		assert.Equal(t, http.StatusOK, r.Code)
		// Cleanup
		os.RemoveAll(filepath.Join(c.OriginalsPath(), "1990"))
		os.RemoveAll(filepath.Join(c.OriginalsPath(), "London"))
	})

	t.Run("not existing file", func(t *testing.T) {
		app, router, _ := NewApiTest()
		PhotoUnstack(router)
		r := PerformRequest(app, "POST", "/api/v1/photos/ps6sg6be2lvl0yh7/files/xxx/unstack")
		assert.Equal(t, http.StatusNotFound, r.Code)
		// t.Logf("RESP: %s", r.Body.String())
	})
}
