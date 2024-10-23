//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/photoprism/photoprism/internal/ai/classify"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/entity/migrate"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/media"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/photoprism/photoprism/pkg/txt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var log = event.Log

// Log logs the error if any and keeps quiet otherwise.
func Log(model, action string, err error) {
	if err != nil {
		log.Errorf("%s: %s (%s)", model, err, action)
	}
}

// UTC returns the current Coordinated Universal Time (UTC).
func UTC() time.Time {
	return time.Now().UTC()
}

// Set UTC as the default for created and updated timestamps.
func init() {
	gormConfig().NowFunc = func() time.Time {
		return UTC()
	}
}

// Db returns the default *gorm.DB connection.
func Db() *gorm.DB {
	if dbConn == nil {
		return nil
	}

	return dbConn.Db()
}

// UnscopedDb returns an unscoped *gorm.DB connection
// that returns all records including deleted records.
func UnscopedDb() *gorm.DB {
	return Db().Unscoped()
}

// Supported test databases.
const (
	MySQL           = "mysql"
	Postgres        = "postgres"
	SQLite3         = "sqlite"
	SQLiteTestDB    = ".test.db"
	SQLiteMemoryDSN = ":memory:?cache=shared&_foreign_keys=on"
)

var drivers = map[string]func(string) gorm.Dialector{
	MySQL:    mysql.Open,
	Postgres: postgres.Open,
	SQLite3:  sqlite.Open,
}

// dbConn is the global gorm.DB connection provider.
var dbConn Gorm

// Gorm is a gorm.DB connection provider interface.
type Gorm interface {
	Db() *gorm.DB
}

// DbConn is a gorm.DB connection provider.
type DbConn struct {
	Driver string
	Dsn    string

	once sync.Once
	db   *gorm.DB
}

// Db returns the gorm db connection.
func (g *DbConn) Db() *gorm.DB {
	g.once.Do(g.Open)

	if g.db == nil {
		log.Fatal("migrate: database not connected")
	}

	return g.db
}

// Open creates a new gorm db connection.
func (g *DbConn) Open() {
	log.Infof("Opening DB connection with driver %s", g.Driver)
	db, err := gorm.Open(drivers[g.Driver](g.Dsn), gormConfig())

	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			fmt.Printf("gorm.Open(%s, %s) %d\n", g.Driver, g.Dsn, i)
			db, err = gorm.Open(drivers[g.Driver](g.Dsn), gormConfig())

			if db != nil && err == nil {
				break
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		if err != nil || db == nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
	log.Info("DB connection established successfully")

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(4)   // in config_db it uses c.DatabaseConnsIdle(), but we don't have the c here.
	sqlDB.SetMaxOpenConns(256) // in config_db it uses c.DatabaseConns(), but we don't have the c here.

	g.db = db
}

// Close closes the gorm db connection.
func (g *DbConn) Close() {
	if g.db != nil {
		sqlDB, _ := g.db.DB()
		if err := sqlDB.Close(); err != nil {
			log.Fatal(err)
		}

		g.db = nil
	}
}

// IsDialect returns true if the given sql dialect is used.
func IsDialect(name string) bool {
	return name == Db().Dialector.Name()
}

// DbDialect returns the sql dialect name.
func DbDialect() string {
	return Db().Dialector.Name()
}

// SetDbProvider sets the Gorm database connection provider.
func SetDbProvider(conn Gorm) {
	dbConn = conn
}

// HasDbProvider returns true if a db provider exists.
func HasDbProvider() bool {
	return dbConn != nil
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log,
			logger.Config{
				SlowThreshold:             time.Second,  // Slow SQL threshold
				LogLevel:                  logger.Error, // Log level
				IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,        // Don't include params in the SQL log
				Colorful:                  false,        // Disable color
			},
		),
		// Set UTC as the default for created and updated timestamps.
		NowFunc: func() time.Time {
			return UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

func main() {
	var (
		numberOfPhotos int
		driver         string
		dsn            string
		dropdb         bool
		sqlitescript   bool
	)

	log = logrus.StandardLogger()
	log.SetLevel(logrus.TraceLevel)
	event.AuditLog = log

	flag.IntVar(&numberOfPhotos, "numberOfPhotos", 0, "Number of photos to generate")
	flag.StringVar(&driver, "driver", "sqlite", "GORM driver to use.  Choose from sqlite, mysql and postgres")
	flag.StringVar(&dsn, "dsn", "testdb.db", "DSN to access the database")
	flag.BoolVar(&dropdb, "dropdb", false, "Drop/Delete the database")
	flag.BoolVar(&sqlitescript, "sqlitescript", true, "Create an SQLite database from script")
	flag.Parse()

	if numberOfPhotos < 1 {
		flag.PrintDefaults()
		log.Errorf("Number of photos is not enough %d", numberOfPhotos)
		os.Exit(1)
	}

	if _, ok := drivers[driver]; ok == false {
		flag.PrintDefaults()
		log.Errorf("driver %v is not valid", driver)
		os.Exit(1)
	}

	if len(dsn) < 3 {
		flag.PrintDefaults()
		log.Errorf("dsn %v is to short", dsn)
		os.Exit(1)
	}

	// Set default test database driver.
	if driver == "test" || driver == "sqlite" || driver == "" || dsn == "" {
		driver = SQLite3
	}

	// Set default database DSN.
	if driver == SQLite3 {
		if dsn == "" {
			dsn = SQLiteMemoryDSN
		}
	}

	allowDelete := dropdb
	if driver == MySQL && allowDelete {
		basedsn := dsn[0 : strings.Index(dsn, "/")+1]
		basedbname := dsn[strings.Index(dsn, "/")+1 : strings.Index(dsn, "?")]
		log.Infof("Connecting to %v", basedsn)
		database, err := gorm.Open(mysql.Open(basedsn), &gorm.Config{})
		if err != nil {
			log.Errorf("Unable to connect to MariaDB %v", err)
		}
		log.Infof("Dropping database %v if it exists", basedbname)
		if res := database.Exec("DROP DATABASE IF EXISTS " + basedbname + ";"); res.Error != nil {
			log.Errorf("Unable to drop database %v", res.Error)
			os.Exit(1)
		}

		log.Infof("Creating database %v if it doesnt exist", basedbname)
		if res := database.Exec("CREATE DATABASE IF NOT EXISTS " + basedbname + ";"); res.Error != nil {
			log.Errorf("Unable to create database %v", res.Error)
			os.Exit(1)
		}
	}
	if driver == SQLite3 && dsn != SQLiteMemoryDSN && allowDelete {
		filename := dsn
		if strings.Index(dsn, "?") > 0 {
			if strings.Index(dsn, ":") > 0 {
				filename = dsn[strings.Index(dsn, ":")+1 : strings.Index(dsn, "?")]
			} else {
				filename = dsn[0:strings.Index(dsn, "?")]
			}
		}
		log.Infof("Removing file %v", filename)
		os.Remove(filename)
	}

	log.Infof("Connecting to driver %v with dsn %v", driver, dsn)
	// Create gorm.DB connection provider.
	db := &DbConn{
		Driver: driver,
		Dsn:    dsn,
	}
	defer db.Close()

	SetDbProvider(db)

	start := time.Now()

	log.Info("Create PhotoPrism tables if they don't exist")
	// Run migration if the photos table doesn't exist.
	// Otherwise assume that we have a valid structured database.
	photoCounter := int64(0)
	if err := Db().Model(entity.Photo{}).Count(&photoCounter).Error; err != nil {
		// Handle SQLite differently as it does table recreates on initial migrate, so we need to be able to simulate that.
		if driver == SQLite3 && sqlitescript {
			filename := dsn
			if strings.Index(dsn, "?") > 0 {
				if strings.Index(dsn, ":") > 0 {
					filename = dsn[strings.Index(dsn, ":")+1 : strings.Index(dsn, "?")]
				} else {
					filename = dsn[0:strings.Index(dsn, "?")]
				}
			}

			var cmd *exec.Cmd

			bashCmd := fmt.Sprintf("cat ./sqlite3.sql | sqlite3 %s", filename)

			cmd = exec.Command("bash", "-c", bashCmd)

			// Write to stdout or file.
			var f *os.File
			log.Infof("restore: creating database tables from script")
			f = os.Stdout
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = f

			// Log exact command for debugging in trace mode.
			log.Debug(cmd.String())

			// Run restore command.
			if cmdErr := cmd.Run(); cmdErr != nil {
				if errStr := strings.TrimSpace(stderr.String()); errStr != "" {
					log.Error(errStr)
					os.Exit(1)
				}
			}
		} else {
			entity.Entities.Migrate(Db(), migrate.Opt(true, false, nil))
			if err := entity.Entities.WaitForMigration(Db()); err != nil {
				log.Errorf("migrate: %s [%s]", err, time.Since(start))
			}
		}
	} else {
		log.Errorf("The photos table already exists in driver %v dsn %v.\nAborting...", driver, dsn)
		os.Exit(1)
	}

	entity.SetDbProvider(dbConn)

	log.Info("Create default fixtures")

	entity.CreateDefaultFixtures()

	// Load the database with data.

	// Create all the labels and keywords that have specific handling in internal/ai/classify/rules.go
	log.Info("Create labels and keywords")
	keywords := make(map[string]uint)
	labels := make(map[string]uint)
	keywordRandoms := make(map[int]uint)
	labelRandoms := make(map[int]uint)
	keywordPos, labelPos := 0, 0
	for label, rule := range classify.Rules {
		keyword := entity.Keyword{
			Keyword: label,
			Skip:    false,
		}
		Db().Create(&keyword)
		keywords[label] = keyword.ID
		keywordRandoms[keywordPos] = keyword.ID
		keywordPos++
		if rule.Label != "" {
			if _, found := keywords[rule.Label]; found == false {
				keyword = entity.Keyword{
					Keyword: rule.Label,
					Skip:    false,
				}
				Db().Create(&keyword)
				keywords[rule.Label] = keyword.ID
				keywordRandoms[keywordPos] = keyword.ID
				keywordPos++
			}
			for _, category := range rule.Categories {
				if _, found := labels[category]; found == false {
					labelDb := entity.Label{
						LabelSlug:        strings.ToLower(category),
						CustomSlug:       strings.ToLower(category),
						LabelName:        strings.ToLower(category),
						LabelPriority:    0,
						LabelFavorite:    false,
						LabelDescription: "",
						LabelNotes:       "",
						PhotoCount:       0,
						LabelCategories:  []*entity.Label{},
						CreatedAt:        time.Now().UTC(),
						UpdatedAt:        time.Now().UTC(),
						DeletedAt:        gorm.DeletedAt{},
						New:              false,
					}
					Db().Create(&labelDb)
					labels[category] = labelDb.ID
					labelRandoms[labelPos] = labelDb.ID
					labelPos++
				}
			}
			if _, found := labels[rule.Label]; found == false {
				labelDb := entity.Label{
					LabelSlug:        strings.ToLower(rule.Label),
					CustomSlug:       strings.ToLower(rule.Label),
					LabelName:        strings.ToLower(rule.Label),
					LabelPriority:    0,
					LabelFavorite:    false,
					LabelDescription: "",
					LabelNotes:       "",
					PhotoCount:       0,
					LabelCategories:  []*entity.Label{},
					CreatedAt:        time.Now().UTC(),
					UpdatedAt:        time.Now().UTC(),
					DeletedAt:        gorm.DeletedAt{},
					New:              false,
				}
				Db().Create(&labelDb)
				labels[rule.Label] = labelDb.ID
				labelRandoms[labelPos] = labelDb.ID
				labelPos++
				for _, category := range rule.Categories {
					categoryDb := entity.Category{
						LabelID:    labelDb.ID,
						CategoryID: labels[category],
					}
					Db().Create(&categoryDb)
				}
			}
		}
	}

	// Create every possible camera and some lenses.  Yeah the data is garbage but it's test data anyway.
	log.Info("Create cameras and lenses")
	lensList := [6]string{"Wide Angle", "Fisheye", "Ultra Wide Angle", "Macro", "Super Zoom", "F80"}
	cameras := make(map[string]uint)
	lenses := make(map[string]uint)
	cameraRandoms := make(map[int]uint)
	lensRandoms := make(map[int]uint)
	cameraPos, lensPos := 0, 0

	for _, make := range entity.CameraMakes {
		for _, model := range entity.CameraModels {
			camera := entity.NewCamera(make, model)
			if _, found := cameras[camera.CameraSlug]; found == false {
				Db().Create(camera)
				cameras[camera.CameraSlug] = camera.ID
				cameraRandoms[cameraPos] = camera.ID
				cameraPos++
			}
		}
		for _, model := range lensList {
			lens := entity.NewLens(make, model)
			if _, found := lenses[lens.LensSlug]; found == false {
				Db().Create(lens)
				lenses[lens.LensSlug] = lens.ID
				lensRandoms[lensPos] = lens.ID
				lensPos++
			}
		}
	}

	// Load up Countries and Places.
	log.Info("Create countries and places")
	countries := make(map[int]string)
	countryPos := 0
	places := make(map[int]string)
	placePos := 0

	PlaceUID := byte('P')

	file, _ := os.Open("../../pkg/txt/resources/countries.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")

		if len(parts) < 2 {
			continue
		}

		country := entity.NewCountry(strings.ToLower(parts[0]), strings.ToLower(parts[1]))
		counter := int64(0)
		Db().Model(entity.Country{}).Where("id = ?", country.ID).Count(&counter)
		if counter == 0 {
			Db().Create(country)
			countries[countryPos] = strings.ToLower(parts[0])
			countryPos++
		}
	}

	for word, _ := range txt.StopWords {
		placeUID := rnd.GenerateUID(PlaceUID)
		country := countries[rand.IntN(len(countries))]
		place := entity.Place{
			ID:            placeUID,
			PlaceLabel:    word,
			PlaceDistrict: word,
			PlaceCity:     word,
			PlaceState:    word,
			PlaceCountry:  country,
			PlaceKeywords: "",
			PlaceFavorite: false,
			PhotoCount:    0,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		Db().Create(&place)
		places[placePos] = placeUID
		placePos++
	}

	// Create some Subjects
	log.Info("Create subjects")
	subjects := make(map[int]entity.Subject)
	subjectPos := 0

	for i := 1; i <= 100; i++ {
		subject := entity.Subject{
			SubjUID:      rnd.GenerateUID('j'),
			SubjType:     entity.SubjPerson,
			SubjSrc:      entity.SrcImage,
			SubjSlug:     fmt.Sprintf("person-%03d", i),
			SubjName:     fmt.Sprintf("Person %03d", i),
			SubjFavorite: false,
			SubjPrivate:  false,
			SubjExcluded: false,
			FileCount:    0,
			PhotoCount:   0,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			DeletedAt:    gorm.DeletedAt{},
		}
		Db().Create(&subject)
		subjects[subjectPos] = subject
		subjectPos++
	}

	log.Info("Start creating photos")
	for i := 1; i <= numberOfPhotos; i++ {
		if _, frac := math.Modf(float64(i) / 100.0); frac == 0 {
			log.Infof("Generating photo number %v", i)
		}
		month := rand.IntN(11) + 1
		day := rand.IntN(28) + 1
		year := rand.IntN(45) + 1980
		takenAt := time.Date(year, time.Month(month), day, rand.IntN(24), rand.IntN(60), rand.IntN(60), rand.IntN(1000), time.UTC)
		labelCount := rand.IntN(5)

		// Create the cell for the Photo's location
		placeId := places[rand.IntN(len(places))]
		lat := (rand.Float64() * 180.0) - 90.0
		lng := (rand.Float64() * 360.0) - 180.0
		cell := entity.NewCell(lat, lng)
		cell.PlaceID = placeId
		Db().FirstOrCreate(cell)

		// Allocate the keywords for this photo
		keywordCount := rand.IntN(5)
		keywords := []entity.Keyword{}
		keywordStr := ""
		for i := 0; i < keywordCount; i++ {
			keyword := entity.Keyword{}
			Db().Model(entity.Keyword{}).Where("id = ?", keywordRandoms[rand.IntN(len(keywordRandoms))]).First(&keyword)
			keywords = append(keywords, keyword)
			if len(keywordStr) > 0 {
				keywordStr = fmt.Sprintf("%s,%s", keywordStr, keyword.Keyword)
			} else {
				keywordStr = keyword.Keyword
			}
		}

		folder := entity.Folder{}
		if res := Db().Model(entity.Folder{}).Where("path = ?", fmt.Sprintf("%04d", year)).First(&folder); res.RowsAffected == 0 {
			folder = entity.NewFolder("/", fmt.Sprintf("%04d", year), time.Now().UTC())
			folder.Create()
		}
		folder = entity.Folder{}
		if res := Db().Model(entity.Folder{}).Where("path = ?", fmt.Sprintf("%04d/%02d", year, month)).First(&folder); res.RowsAffected == 0 {
			folder = entity.NewFolder("/", fmt.Sprintf("%04d/%02d", year, month), time.Now().UTC())
			folder.Create()
		}

		photo := entity.Photo{
			//	ID
			//
			// UUID
			TakenAt:          takenAt,
			TakenAtLocal:     takenAt,
			TakenSrc:         entity.SrcMeta,
			PhotoUID:         rnd.GenerateUID(entity.PhotoUID),
			PhotoType:        "image",
			TypeSrc:          entity.SrcAuto,
			PhotoTitle:       "Performance Test Load",
			TitleSrc:         entity.SrcImage,
			PhotoDescription: "",
			DescriptionSrc:   entity.SrcAuto,
			PhotoPath:        fmt.Sprintf("%04d/%02d", year, month),
			PhotoName:        fmt.Sprintf("PIC%08d", i),
			OriginalName:     fmt.Sprintf("PIC%08d", i),
			PhotoStack:       0,
			PhotoFavorite:    false,
			PhotoPrivate:     false,
			PhotoScan:        false,
			PhotoPanorama:    false,
			TimeZone:         "America/Mexico_City",
			PlaceID:          placeId,
			PlaceSrc:         entity.SrcMeta,
			CellID:           cell.ID,
			CellAccuracy:     0,
			PhotoAltitude:    5,
			PhotoLat:         lat,
			PhotoLng:         lng,
			PhotoCountry:     countries[rand.IntN(len(countries))],
			PhotoYear:        year,
			PhotoMonth:       month,
			PhotoDay:         day,
			PhotoIso:         400,
			PhotoExposure:    "1/60",
			PhotoFNumber:     8,
			PhotoFocalLength: 2,
			PhotoQuality:     3,
			PhotoFaces:       0,
			PhotoResolution:  0,
			// PhotoDuration    : 0,
			PhotoColor:   12,
			CameraID:     cameraRandoms[rand.IntN(len(cameraRandoms))],
			CameraSerial: "",
			CameraSrc:    "",
			LensID:       lensRandoms[rand.IntN(len(lensRandoms))],
			// Details          :,
			// Camera
			// Lens
			// Cell
			// Place
			Keywords: keywords,
			Albums:   []entity.Album{},
			Files:    []entity.File{},
			Labels:   []entity.PhotoLabel{},
			// CreatedBy
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			EditedAt:    nil,
			PublishedAt: nil,
			CheckedAt:   nil,
			EstimatedAt: nil,
			DeletedAt:   gorm.DeletedAt{},
		}
		Db().Create(&photo)
		for i := 0; i < labelCount; i++ {
			photoLabel := entity.NewPhotoLabel(photo.ID, labelRandoms[rand.IntN(len(labelRandoms))], 0, entity.SrcMeta)
			Db().FirstOrCreate(photoLabel)
		}
		// Create File
		file := entity.File{
			//	ID
			// Photo
			PhotoID:      photo.ID,
			PhotoUID:     photo.PhotoUID,
			PhotoTakenAt: photo.TakenAt,
			// TimeIndex
			// MediaID
			// MediaUTC
			InstanceID:   "",
			FileUID:      rnd.GenerateUID(entity.FileUID),
			FileName:     fmt.Sprintf("%04d/%02d/PIC%08d.jpg", year, month, i),
			FileRoot:     entity.RootSidecar,
			OriginalName: "",
			FileHash:     rnd.GenerateUID(entity.FileUID),
			FileSize:     rand.Int64N(1000000),
			FileCodec:    "",
			FileType:     string(fs.ImageJPEG),
			MediaType:    string(media.Image),
			FileMime:     "image/jpg",
			FilePrimary:  true,
			FileSidecar:  false,
			FileMissing:  false,
			FilePortrait: true,
			FileVideo:    false,
			FileDuration: 0,
			// FileFPS
			// FileFrames
			FileWidth:          1200,
			FileHeight:         1600,
			FileOrientation:    6,
			FileOrientationSrc: entity.SrcMeta,
			FileProjection:     "",
			FileAspectRatio:    0.75,
			// FileHDR            : false,
			// FileWatermark
			// FileColorProfile
			FileMainColor: "magenta",
			FileColors:    "226611CC1",
			FileLuminance: "ABCDEF123",
			FileDiff:      456,
			FileChroma:    15,
			// FileSoftware
			// FileError
			ModTime:   time.Now().Unix(),
			CreatedAt: time.Now().UTC(),
			CreatedIn: 935962,
			UpdatedAt: time.Now().UTC(),
			UpdatedIn: 935962,
			// PublishedAt
			DeletedAt: gorm.DeletedAt{},
			Share:     []entity.FileShare{},
			Sync:      []entity.FileSync{},
			//markers
		}
		Db().Create(&file)

		// Add Markers
		markersToCreate := rand.IntN(5)
		for i := 0; i < markersToCreate; i++ {
			subject := subjects[rand.IntN(len(subjects))]
			marker := entity.Marker{
				MarkerUID:     rnd.GenerateUID('m'),
				FileUID:       file.FileUID,
				MarkerType:    entity.MarkerFace,
				MarkerName:    subject.SubjName,
				MarkerReview:  false,
				MarkerInvalid: false,
				SubjUID:       subject.SubjUID,
				SubjSrc:       subject.SubjSrc,
				X:             rand.Float32() * 1024.0,
				Y:             rand.Float32() * 2048.0,
				W:             rand.Float32() * 10.0,
				H:             rand.Float32() * 20.0,
				Q:             10,
				Size:          100,
				Score:         10,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			}
			Db().Create(&marker)
		}

		// Add to Album
		albumSlug := fmt.Sprintf("my-photos-from-%04d", year)
		album := entity.Album{}
		if res := Db().Model(entity.Album{}).Where("album_slug = ?", albumSlug).First(&album); res.RowsAffected == 0 {
			album = entity.Album{
				AlbumUID:         rnd.GenerateUID(entity.AlbumUID),
				AlbumSlug:        albumSlug,
				AlbumPath:        "",
				AlbumType:        entity.AlbumManual,
				AlbumTitle:       fmt.Sprintf("My Photos From %04d", year),
				AlbumLocation:    "",
				AlbumCategory:    "",
				AlbumCaption:     "",
				AlbumDescription: "A wonderful year",
				AlbumNotes:       "",
				AlbumFilter:      "",
				AlbumOrder:       "oldest",
				AlbumTemplate:    "",
				AlbumCountry:     entity.UnknownID,
				AlbumYear:        year,
				AlbumMonth:       0,
				AlbumDay:         0,
				AlbumFavorite:    false,
				AlbumPrivate:     false,
				CreatedAt:        time.Now().UTC(),
				UpdatedAt:        time.Now().UTC(),
				DeletedAt:        gorm.DeletedAt{},
			}
			Db().Create(&album)
		}
		photoAlbum := entity.PhotoAlbum{
			PhotoUID:  photo.PhotoUID,
			AlbumUID:  album.AlbumUID,
			Order:     0,
			Hidden:    false,
			Missing:   false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		Db().Create(photoAlbum)

		details := entity.Details{
			PhotoID:     photo.ID,
			Keywords:    keywordStr,
			KeywordsSrc: entity.SrcMeta,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		Db().Create(details)
	}

	entity.File{}.RegenerateIndex()
	entity.UpdateCounts()

	log.Infof("Database Creation completed in %s", time.Since(start))
	code := 0
	os.Exit(code)

}
