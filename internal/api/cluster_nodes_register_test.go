package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/service/cluster"
	reg "github.com/photoprism/photoprism/internal/service/cluster/registry"
)

func TestClusterNodesRegister(t *testing.T) {
	t.Run("FeatureDisabled", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeRole = cluster.RoleInstance
		ClusterNodesRegister(router)

		r := PerformRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`)
		assert.Equal(t, http.StatusForbidden, r.Code)
	})

	t.Run("MissingToken", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeRole = cluster.RolePortal
		ClusterNodesRegister(router)

		r := PerformRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`)
		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("DriverConflict", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.SQLite3 {
			t.Skip("test requires SQLite driver in test config")
		}
		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// With SQLite driver in tests, provisioning should fail with conflict.
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`, "t0k3n")
		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Contains(t, r.Body.String(), "portal database must be MySQL/MariaDB")
	})

	t.Run("BadName", func(t *testing.T) {
		app, router, conf := NewApiTest()
		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// Empty nodeName → 400
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":""}`, "t0k3n")
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

	t.Run("RotateSecretPersistsDespiteDBConflict", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.SQLite3 {
			t.Skip("test requires SQLite driver in test config")
		}

		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// Pre-create node in registry so handler goes through existing-node path
		// and rotates the secret before attempting DB ensure.
		regy, err := reg.NewClientRegistryWithConfig(conf)
		assert.NoError(t, err)
		n := &reg.Node{ID: "test-id", Name: "pp-node-01", Role: "instance"}
		assert.NoError(t, regy.Put(n))

		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01","rotateSecret":true}`, "t0k3n")
		assert.Equal(t, http.StatusConflict, r.Code) // DB conflict under SQLite

		// Secret should have rotated and been persisted even though DB ensure failed.
		// Fetch by name (most-recently-updated) to avoid flakiness if another test adds
		// a node with the same name and a different id.
		n2, err := regy.FindByName("pp-node-01")
		assert.NoError(t, err)
		// With client-backed registry, plaintext secret is not persisted; only rotation timestamp is updated.
		assert.NotEmpty(t, n2.SecretRot)
	})

	t.Run("ExistingNodeSiteUrlPersistsEvenOnDBConflict", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.SQLite3 {
			t.Skip("test requires SQLite driver in test config")
		}

		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// Pre-create node in registry so handler goes through existing-node path.
		regy, err := reg.NewClientRegistryWithConfig(conf)
		assert.NoError(t, err)
		n := &reg.Node{Name: "pp-node-02", Role: "instance"}
		assert.NoError(t, regy.Put(n))

		// With SQLite driver in tests, provisioning should fail with 409, but metadata should still persist.
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-02","siteUrl":"https://Photos.Example.COM"}`, "t0k3n")
		assert.Equal(t, http.StatusConflict, r.Code)

		// Ensure normalized/persisted siteUrl.
		n2, err := regy.FindByName("pp-node-02")
		assert.NoError(t, err)
		assert.Equal(t, "https://photos.example.com", n2.SiteUrl)
	})

	t.Run("ValidDriverCreate", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.MySQL {
			t.Skip("test requires MariaDB driver in test config")
		}
		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// With MariaDB driver in tests, provisioning should succeed.
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`, "t0k3n")
		assert.Equal(t, http.StatusCreated, r.Code)
		assert.Contains(t, r.Body.String(), "node")
		if assert.Contains(t, r.Body.String(), "database") {
			cleanupDatabases(r.Body.Bytes(), conf, t)
		}
	})

	t.Run("RotateSecretPersists", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.MySQL {
			t.Skip("test requires MariaDB driver in test config")
		}

		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// Register the node to ensure that the database and registry is there
		rCreate := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01"}`, "t0k3n")
		assert.Equal(t, http.StatusCreated, rCreate.Code)
		assert.Contains(t, rCreate.Body.String(), `"alreadyProvisioned":false`)

		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-01","rotateSecret":true}`, "t0k3n")
		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), `"alreadyProvisioned":true`)

		// Secret should have rotated and been persisted.
		// Fetch by name (most-recently-updated) to avoid flakiness if another test adds
		// a node with the same name and a different id.
		regy, err := reg.NewClientRegistryWithConfig(conf)
		assert.NoError(t, err)
		n2, err := regy.FindByName("pp-node-01")
		assert.NoError(t, err)
		// With client-backed registry, plaintext secret is not persisted; only rotation timestamp is updated.
		assert.NotEmpty(t, n2.SecretRot)

		assert.Contains(t, rCreate.Body.String(), "node")
		if assert.Contains(t, rCreate.Body.String(), "database") {
			cleanupDatabases(rCreate.Body.Bytes(), conf, t)
		}
	})

	t.Run("ExistingNodeSiteUrlPersists", func(t *testing.T) {
		app, router, conf := NewApiTest()
		if conf.Options().DatabaseDriver != config.MySQL {
			t.Skip("test requires MariaDB driver in test config")
		}

		conf.Options().NodeRole = cluster.RolePortal
		conf.Options().JoinToken = "t0k3n"
		ClusterNodesRegister(router)

		// Register the node to ensure that the database and registry is there
		rCreate := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-02"}`, "t0k3n")
		assert.Equal(t, http.StatusCreated, rCreate.Code)
		assert.Contains(t, rCreate.Body.String(), `"alreadyProvisioned":false`)

		// With MariaDB driver in tests, provisioning should Succeed with 200 and metadata should persist.
		r := AuthenticatedRequestWithBody(app, http.MethodPost, "/api/v1/cluster/nodes/register", `{"nodeName":"pp-node-02","siteUrl":"https://Photos.Example.COM"}`, "t0k3n")
		assert.Equal(t, http.StatusOK, r.Code)

		// Ensure normalized/persisted siteUrl.
		regy, err := reg.NewClientRegistryWithConfig(conf)
		assert.NoError(t, err)
		n2, err := regy.FindByName("pp-node-02")
		assert.NoError(t, err)
		assert.Equal(t, "https://photos.example.com", n2.SiteUrl)

		assert.Contains(t, rCreate.Body.String(), "node")
		if assert.Contains(t, rCreate.Body.String(), "database") {
			cleanupDatabases(rCreate.Body.Bytes(), conf, t)
		}
	})
}

func quoteIdent(s string) string { return "`" + strings.ReplaceAll(s, "`", "``") + "`" }

// cleanupDatabases expects a byte array that contains a cluster.RegisterResponse, config.Config and testing.T and drops the database created by the register.
func cleanupDatabases(jb []byte, c *config.Config, t *testing.T) {
	var resp cluster.RegisterResponse
	json.Unmarshal(jb, &resp)
	log.Debugf("Cleanup Database %s, User %s and client_uid %s", resp.Database.Name, resp.Database.User, resp.Node.ID)
	if err := c.Db().Unscoped().Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", quoteIdent(resp.Database.Name))).Error; err != nil {
		assert.Empty(t, err)
		t.Logf("Unable to drop database %s", quoteIdent(resp.Database.Name))
	}
	if err := c.Db().Unscoped().Exec(fmt.Sprintf("DROP USER IF EXISTS %s", quoteIdent(resp.Database.User))).Error; err != nil {
		assert.Empty(t, err)
		t.Logf("Unable to drop user %s", quoteIdent(resp.Database.User))
	}
	if err := c.Db().Unscoped().Exec("DELETE FROM auth_clients WHERE client_uid = ?", resp.Node.ID).Error; err != nil {
		assert.Empty(t, err)
		t.Logf("Unable to remove client_uid %s", quoteIdent(resp.Node.ID))
	}
}
