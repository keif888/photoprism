package testextras

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/photoprism/photoprism/pkg/fs"
)

var storagePath string
var basePath string
var originalPath string

// SetupStorage creates a unique storage folder and assigns to PHOTOPRISM_STORAGE_PATH so test packages are independant
func SetupStorage() (err error) {
	var tmpPath string
	if cwd, err := os.Getwd(); err == nil {
		tmpPaths := strings.SplitAfter(cwd, "photoprism/photoprism")
		if len(tmpPaths) == 2 {
			tmpPath = filepath.Join(tmpPaths[0], fs.StorageDir)
		}
	} else {
		log.Warningf("testextras: Getwd error %s", err.Error())
	}

	if basePath, err = os.MkdirTemp(tmpPath, "testextras-*"); err != nil {
		log.Errorf("testextras: MkdirTemp error %s", err.Error())
		return err
	}

	storagePath = filepath.Join(basePath, fs.StorageDir)
	if err = os.MkdirAll(storagePath, os.ModePerm); err != nil {
		log.Errorf("testextras: MkdirTemp error %s", err.Error())
		return err
	}
	log.Debugf("testextras: created %+v", storagePath)

	originalPath = os.Getenv("PHOTOPRISM_STORAGE_PATH")
	return os.Setenv("PHOTOPRISM_STORAGE_PATH", storagePath)
}

// CleanupStorage removes the unique storage folder
func CleanupStorage() (err error) {
	if storagePath == "" {
		return nil
	}

	err = os.Setenv("PHOTOPRISM_STORAGE_PATH", originalPath)

	if errra := os.RemoveAll(basePath); err != nil {
		log.Errorf("testextras: %s", errra.Error())
		return errra
	}
	log.Debugf("testextras: cleaned up %+v", basePath)
	return err
}
