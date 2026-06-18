package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileMap_Get(t *testing.T) {
	t.Run("GetExistingFile", func(t *testing.T) {
		r := FileFixtures.Get("exampleFileName.jpg")
		assert.Equal(t, "fs6sg6bw45bnlqdw", r.FileUID)
		assert.Equal(t, "2790/07/27900704_070228_D6D51B6C.jpg", r.FileName)
		assert.IsType(t, File{}, r)
	})
	t.Run("GetNotExistingFile", func(t *testing.T) {
		r := FileFixtures.Get("TestName")
		assert.Equal(t, "TestName", r.FileName)
		assert.IsType(t, File{}, r)
	})
}

func TestFileMap_Pointer(t *testing.T) {
	t.Run("GetExistingFile", func(t *testing.T) {
		r := FileFixtures.Pointer("exampleFileName.jpg")
		assert.Equal(t, "fs6sg6bw45bnlqdw", r.FileUID)
		assert.Equal(t, "2790/07/27900704_070228_D6D51B6C.jpg", r.FileName)
		assert.IsType(t, &File{}, r)
	})
	t.Run("GetNotExistingFile", func(t *testing.T) {
		r := FileFixtures.Pointer("TestName")
		assert.Equal(t, "TestName", r.FileName)
		assert.IsType(t, &File{}, r)
	})
}

func TestFileFixtureLoad(t *testing.T) {
	t.Run("ZeroValuesFileDiff", func(t *testing.T) {
		e := FileFixtures.Get("Photo06.png")
		a := &File{}

		if assert.NoError(t, Db().First(&a, e.ID).Error) {
			assert.Equal(t, e.ID, a.ID)
			if assert.NotNil(t, e.FileRoot) && assert.NotNil(t, a.FileRoot) {
				assert.Equal(t, *e.FileRoot, *a.FileRoot, "FileRoot")
			}
			assert.Equal(t, e.FilePages, a.FilePages, "FilePages")
			assert.Equal(t, e.FileOrientationSrc, a.FileOrientationSrc, "FileOrientationSrc")
			if assert.NotNil(t, e.FileDiff) && assert.NotNil(t, a.FileDiff) {
				assert.Equal(t, *e.FileDiff, *a.FileDiff, "FileDiff")
			}
			if assert.NotNil(t, e.FileChroma) && assert.NotNil(t, a.FileChroma) {
				assert.Equal(t, *e.FileChroma, *a.FileChroma, "FileChroma")
			}
		}
	})
	t.Run("ZeroValuesFileChroma", func(t *testing.T) {
		e := FileFixtures.Get("exampleXmpFile.xmp")
		a := &File{}

		if assert.NoError(t, Db().First(&a, e.ID).Error) {
			assert.Equal(t, e.ID, a.ID)
			if assert.NotNil(t, e.FileRoot) && assert.NotNil(t, a.FileRoot) {
				assert.Equal(t, *e.FileRoot, *a.FileRoot, "FileRoot")
			}
			assert.Equal(t, e.FilePages, a.FilePages, "FilePages")
			assert.Equal(t, e.FileOrientationSrc, a.FileOrientationSrc, "FileOrientationSrc")
			if assert.NotNil(t, e.FileDiff) && assert.NotNil(t, a.FileDiff) {
				assert.Equal(t, *e.FileDiff, *a.FileDiff, "FileDiff")
			}
			if assert.NotNil(t, e.FileChroma) && assert.NotNil(t, a.FileChroma) {
				assert.Equal(t, *e.FileChroma, *a.FileChroma, "FileChroma")
			}
		}
	})
}
