package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShareToken(t *testing.T) {
	t.Run("InvalidToken", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareToken(router)
		r := PerformRequest(app, "GET", "/api/v1/xxx")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})
	// TODO Why does it panic?
	/*t.Run("ValidToken", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareToken(router)
		r := PerformRequest(app, "GET", "/api/v1/4jxf3jfn2k")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})*/
}

func TestShareBootstrapConfig(t *testing.T) {
	t.Run("ClearsUnusedValues", func(t *testing.T) {
		_, _, conf := NewApiTest()

		cfg := conf.ClientShare()
		cfg.PreviewToken = "preview-token"
		cfg.DownloadToken = "download-token"
		cfg.MapKey = "map-key"
		cfg.Customer = "Acme Corp"
		cfg.SiteUrl = "https://app.example.com/"

		result := shareBootstrapConfig(cfg)

		assert.Empty(t, result.PreviewToken)
		assert.Empty(t, result.DownloadToken)
		assert.Empty(t, result.MapKey)
		assert.Empty(t, result.Customer)
		assert.Equal(t, "https://app.example.com/", result.SiteUrl)
		assert.NotNil(t, result.Settings)
	})
}

func TestShareTokenShared(t *testing.T) {
	t.Run("InvalidToken", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareTokenShared(router)
		r := PerformRequest(app, "GET", "/api/v1/1jxf3jfn2k/ss6sg6bxpogaaba7")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})
	t.Run("RejectsOversizedToken", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareTokenShared(router)
		r := PerformRequest(app, "GET", "/api/v1/"+strings.Repeat("a", 161)+"/as6sg6bxpogaaba8")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})
	t.Run("RejectsUnusableToken", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareTokenShared(router)
		r := PerformRequest(app, "GET", "/api/v1/..../as6sg6bxpogaaba8")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})
	// TODO Why does it panic?
	/*t.Run("ValidTokenAndShare", func(t *testing.T) {
		app, router, _ := NewApiTest()
		ShareTokenShared(router)
		r := PerformRequest(app, "GET", "/api/v1/4jxf3jfn2k/as6sg6bxpogaaba7")
		assert.Equal(t, http.StatusTemporaryRedirect, r.Code)
	})*/
}
