package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
)

func TestPhotoByID(t *testing.T) {
	t.Run("photo found", func(t *testing.T) {
		result, err := PhotoByID(1000000)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 2790, result.PhotoYear)
	})

	t.Run("no photo found", func(t *testing.T) {
		result, err := PhotoByID(99999)
		assert.Error(t, err, "record not found")
		t.Log(result)
	})
}

func TestPhotoByUID(t *testing.T) {
	t.Run("photo found", func(t *testing.T) {
		result, err := PhotoByUID("ps6sg6be2lvl0y12")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Reunion", result.PhotoTitle)
	})

	t.Run("no photo found", func(t *testing.T) {
		result, err := PhotoByUID("99999")
		assert.Error(t, err, "record not found")
		t.Log(result)
	})
}

func TestPreloadPhotoByUID(t *testing.T) {
	t.Run("photo found", func(t *testing.T) {
		result, err := PhotoPreloadByUID("ps6sg6be2lvl0y12")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Reunion", result.PhotoTitle)
	})

	t.Run("no photo found", func(t *testing.T) {
		result, err := PhotoPreloadByUID("99999")
		assert.Error(t, err, "record not found")
		t.Log(result)
	})
}

func TestMissingPhotos(t *testing.T) {
	result, err := MissingPhotos(15, 0)

	if err != nil {
		t.Fatal(err)
	}

	assert.LessOrEqual(t, 1, len(result))
}

func TestArchivedPhotos(t *testing.T) {
	results, err := ArchivedPhotos(15, 0)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(results))

	if len(results) > 1 {
		result := results[0]
		assert.Equal(t, "image", result.PhotoType)
		assert.Equal(t, "ps6sg6be2lvl0y25", result.PhotoUID)
	}
}

func TestPhotosMetadataUpdate(t *testing.T) {
	interval := entity.MetadataUpdateInterval
	result, err := PhotosMetadataUpdate(10, 0, time.Second, interval)

	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, entity.Photos{}, result)
}

func TestOrphanPhotos(t *testing.T) {
	result, err := OrphanPhotos()

	if err != nil {
		t.Fatal(err)
	}

	assert.IsType(t, entity.Photos{}, result)
}

func TestFixPrimaries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		preview := entity.File{
			ID:              3000002,
			Photo:           entity.PhotoFixtures.Pointer("Photo01"),
			PhotoID:         entity.PhotoFixtures.Pointer("Photo01").ID,
			PhotoUID:        entity.PhotoFixtures.Pointer("Photo01").PhotoUID,
			InstanceID:      "a698ac56-6e7e-42b9-9c3e-a79ec96080ac",
			FileUID:         rnd.GenerateUID(entity.FileUID), // "fs6sg6bw45bn0003",
			FileName:        "2790/02/Photo01.jpg",
			FileRoot:        entity.RootOriginals,
			OriginalName:    "",
			FileHash:        "ocad9168fa6acc5c5c2965ddf6ec465ca42fd812",
			FileSize:        858,
			FileCodec:       "",
			FileType:        "jpg",
			MediaType:       string(fs.ImageJpeg),
			FileMime:        "",
			FilePrimary:     false,
			FileSidecar:     false,
			FileVideo:       false,
			FileMissing:     false,
			FilePortrait:    false,
			FileDuration:    0,
			FileWidth:       100,
			FileHeight:      0,
			FileOrientation: 0,
			FileProjection:  "",
			FileAspectRatio: 1,
			FileMainColor:   "",
			FileColors:      "",
			FileLuminance:   "",
			FileDiff:        0,
			FileChroma:      0,
			FileError:       "",
			Share:           []entity.FileShare{},
			Sync:            []entity.FileSync{},
			ModTime:         time.Date(2019, 3, 6, 2, 6, 51, 0, time.UTC).Unix(),
			CreatedAt:       time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedIn:       12361491,
			UpdatedAt:       time.Date(2020, 3, 28, 14, 6, 0, 0, time.UTC),
			UpdatedIn:       9537701,
			DeletedAt:       nil,
		}

		if err := Db().Create(&preview).Error; err != nil {
			t.Fatal(err)
		}
		var count int64
		if err := Db().Model(entity.File{}).Where("id = ? and file_primary = ?", 3000002, 1).Count(&count).Error; err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(0), count)

		err := FixPrimaries()
		if err != nil {
			t.Fatal(err)
		}

		if err := Db().Model(entity.File{}).Where("id = ? and file_primary = ?", 3000002, 1).Count(&count).Error; err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(1), count)
	})
}

func TestFlagHiddenPhotos(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Set photo quality scores to -1 if files are missing.
		if err := FlagHiddenPhotos(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SuccessWith1000", func(t *testing.T) {
		var checkedTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		// Load 1000 photos that need to be hidden
		for i := 0; i < 1000; i++ {
			newPhoto := entity.Photo{ //JPG, Geo from metadata, indexed
				//ID:               1000049,
				PhotoUID:         rnd.GenerateUID(entity.PhotoUID),
				TakenAt:          time.Date(2020, 11, 11, 9, 7, 18, 0, time.UTC),
				TakenAtLocal:     time.Date(2020, 11, 11, 9, 7, 18, 0, time.UTC),
				TakenSrc:         entity.SrcMeta,
				PhotoType:        "image",
				TypeSrc:          "",
				PhotoTitle:       "desk\"",
				TitleSrc:         entity.SrcManual,
				PhotoCaption:     "",
				CaptionSrc:       "",
				PhotoPath:        "2000\"/02\"",
				PhotoName:        "SuccessWith1000",
				OriginalName:     "",
				PhotoFavorite:    false,
				PhotoPrivate:     false,
				PhotoScan:        false,
				PhotoPanorama:    false,
				TimeZone:         "America/Mexico_City",
				PlaceSrc:         "meta",
				CellAccuracy:     0,
				PhotoAltitude:    3,
				PhotoLat:         48.519234,
				PhotoLng:         9.057997,
				PhotoCountry:     entity.CellFixtures.Pointer("caravan park").Place.CountryCode(),
				PhotoYear:        2020,
				PhotoMonth:       11,
				PhotoDay:         11,
				PhotoIso:         0,
				PhotoExposure:    "",
				PhotoFocalLength: 0,
				PhotoFNumber:     0,
				PhotoQuality:     5,
				PhotoResolution:  0,
				Camera:           entity.CameraFixtures.Pointer("canon-eos-6d"),
				CameraID:         entity.CameraFixtures.Pointer("canon-eos-6d").ID,
				CameraSerial:     "",
				CameraSrc:        "",
				Lens:             entity.LensFixtures.Pointer("lens-f-380"),
				LensID:           entity.LensFixtures.Pointer("lens-f-380").ID,
				Keywords:         []entity.Keyword{},
				Albums:           []entity.Album{},
				Files:            []entity.File{},
				Labels:           []entity.PhotoLabel{},
				CreatedAt:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				EditedAt:         nil,
				CheckedAt:        &checkedTime,
				DeletedAt:        nil,
				PhotoColor:       14,
				PhotoStack:       0,
				PhotoFaces:       0,
			}
			if err := Db().Create(&newPhoto).Error; err != nil {
				t.Fatal(err)
			}
		}
		// Set photo quality scores to -1 if files are missing.
		if err := FlagHiddenPhotos(); err != nil {
			t.Fatal(err)
		}

		var actual int64
		var expected int64 = 1000
		if err := Db().Model(&entity.Photo{}).Where("photo_name = ? AND photo_quality = ?", "SuccessWith1000", -1).Count(&actual).Error; err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, actual)

		if err := UnscopedDb().Where("photo_name = ? AND photo_quality = ?", "SuccessWith1000", -1).Delete(&entity.Photo{}).Error; err != nil {
			t.Fatal(err)
		}
	})
}
