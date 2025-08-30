package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
)

func TestFilesByPath(t *testing.T) {
	t.Run("Holiday", func(t *testing.T) {
		files, err := FilesByPath(10, 0, entity.RootOriginals, "Holiday", false)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(files))
	})
	t.Run("SlashHoliday", func(t *testing.T) {
		files, err := FilesByPath(10, 0, entity.RootOriginals, "/Holiday", false)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(files))
	})
	t.Run("2790Public", func(t *testing.T) {
		files, err := FilesByPath(10, 0, entity.RootOriginals, "2790/07", true)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(files))
	})
	t.Run("HolidayPublic", func(t *testing.T) {
		files, err := FilesByPath(10, 0, entity.RootOriginals, "Holiday", true)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}

		assert.LessOrEqual(t, 1, len(files))
	})
}

func TestExistingFiles(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		files, err := Files(1000, 0, "/", true)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}
		assert.LessOrEqual(t, 5, len(files))
	})
	t.Run("files found - includeMissing false", func(t *testing.T) {
		files, err := Files(1000, 0, "/", false)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}
		assert.LessOrEqual(t, 5, len(files))
	})
	t.Run("search for files path", func(t *testing.T) {
		files, err := Files(1000, 0, "Photos", true)

		t.Logf("files: %+v", files)

		if err != nil {
			t.Fatal(err)
		}

		assert.Empty(t, files)
	})
}

func TestFilesByUID(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		files, err := FilesByUID([]string{"fs6sg6bw45bnlqdw"}, 100, 0)

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(files))
		assert.Equal(t, "2790/07/27900704_070228_D6D51B6C.jpg", files[0].FileName)
	})
	t.Run("no files found", func(t *testing.T) {
		files, err := FilesByUID([]string{"fs6sg6bw45bnlxxx"}, 100, 0)

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 0, len(files))
	})
	//fails on mariadb - Gorm issue, fixed in Gorm 2
	t.Run("Error", func(t *testing.T) {
		files, err := FilesByUID([]string{"fs6sg6bw45bnlxxx"}, -100, 0)

		assert.Error(t, err)
		assert.Equal(t, 0, len(files))
	})
}

func TestFileByPhotoUID(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		file, err := FileByPhotoUID("ps6sg6be2lvl0y11")

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Germany/bridge.jpg", file.FileName)
	})

	t.Run("no files found", func(t *testing.T) {
		file, err := FileByPhotoUID("111")

		assert.Error(t, err, "record not found")
		t.Log(file)
	})
}

func TestVideoByPhotoUID(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		file, err := VideoByPhotoUID("ps6sg6be2lvl0yh0")

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "1990/04/bridge2.mp4", file.FileName)
	})

	t.Run("no files found", func(t *testing.T) {
		file, err := VideoByPhotoUID("111")

		assert.Error(t, err, "record not found")
		t.Log(file)
	})
}

func TestFileByUID(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		file, err := FileByUID("fs6sg6bw45bnlqdw")

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2790/07/27900704_070228_D6D51B6C.jpg", file.FileName)
	})

	t.Run("no files found", func(t *testing.T) {
		file, err := FileByUID("111")

		if err == nil {
			t.Fatal("error expected")
		}

		assert.Error(t, err, "record not found")
		t.Log(file)
	})
}

func TestFileByHash(t *testing.T) {
	t.Run("files found", func(t *testing.T) {
		file, err := FileByHash("2cad9168fa6acc5c5c2965ddf6ec465ca42fd818")

		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "2790/07/27900704_070228_D6D51B6C.jpg", file.FileName)
	})

	t.Run("no files found", func(t *testing.T) {
		file, err := FileByHash("111")

		assert.Error(t, err, "record not found")
		t.Log(file)
	})
}

func TestSetPhotoPrimary(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		assert.Equal(t, false, entity.FileFixturesExampleXMP.FilePrimary)

		err := SetPhotoPrimary("ps6sg6be2lvl0yh7", "fs6sg6bwhhbnlqdn")

		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("no_file_uid", func(t *testing.T) {
		err := SetPhotoPrimary("ps6sg6be2lvl0yh7", "")

		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("no_uid", func(t *testing.T) {
		err := SetPhotoPrimary("", "")

		if err == nil {
			t.Fatal("error expected")
		}
	})
	t.Run("file missing", func(t *testing.T) {
		err := SetPhotoPrimary("ps6sg6be2lvl0y22", "")

		if err == nil {
			t.Fatal("error expected")
		}
		assert.Contains(t, err.Error(), "cannot find primary file")
	})
}

func TestSetFileError(t *testing.T) {
	assert.Equal(t, "", entity.FileFixturesExampleXMP.FileError)

	expected := "errorFromTest"
	fileUID := "fs6sg6bwhhbnlqdn"

	SetFileError(fileUID, expected)

	var actual []string
	if err := Db().Model(entity.File{}).Where("file_uid = ?", fileUID).Pluck("file_error", &actual).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(actual))
	if len(actual) > 0 {
		assert.Equal(t, expected, actual[0])
	}
	// Revert the changes
	if err := Db().Save(&entity.FileFixturesExampleXMP).Error; err != nil {
		t.Fatal(err)
	}
}

func TestRenameFile(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		err := RenameFile("xxx", "", "yyy", "yyy")

		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("Success", func(t *testing.T) {
		assert.Equal(t, "2790/02/Photo01.xmp", entity.FileFixturesExampleXMP.FileName)
		assert.Equal(t, "/", entity.FileFixturesExampleXMP.FileRoot)

		// Ensure that the database matches the fixture
		if err := Db().Save(&entity.FileFixturesExampleXMP).Error; err != nil {
			t.Fatal(err)
		}

		fileRoot := "/"
		fileName := entity.FileFixturesExampleXMP.FileName
		newRoot := "test-root"
		newName := "yyy.jpg"

		var expectedFileIDs []string
		if err := Db().Model(entity.File{}).Where("file_root = ? AND file_name = ?", fileRoot, fileName).Pluck("id", &expectedFileIDs).Error; err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(expectedFileIDs))

		err := RenameFile(fileRoot, fileName, newRoot, newName)

		if err != nil {
			t.Fatal(err)
		}

		var actualFileIDs []string
		if err := Db().Model(entity.File{}).Where("file_root = ? AND file_name = ?", newRoot, newName).Pluck("id", &actualFileIDs).Error; err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, len(actualFileIDs))
		if len(actualFileIDs) > 0 && len(expectedFileIDs) > 0 {
			assert.Equal(t, expectedFileIDs[0], actualFileIDs[0])
		}

		// Revert the changes
		if err := Db().Save(&entity.FileFixturesExampleXMP).Error; err != nil {
			t.Fatal(err)
		}
	})

}

func TestIndexedFiles(t *testing.T) {
	if err := entity.AddDuplicate(
		"Photo18.jpg",
		entity.RootSidecar,
		"3cad9168fa6acc5c5c2965ddf6ec465ca42fd818",
		661858,
		time.Date(2019, 3, 6, 2, 6, 51, 0, time.UTC).Unix(),
	); err != nil {
		t.Fatal(err)
	}

	result, err := IndexedFiles()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("indexed files: %#v", result)
}

func TestOrphanFiles(t *testing.T) {
	files, err := OrphanFiles()

	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, entity.Files{}, files)
}
