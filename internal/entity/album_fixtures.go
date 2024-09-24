package entity

import (
	"time"

	"gorm.io/gorm"
)

type AlbumMap map[string]Album

func (m AlbumMap) Get(name string) Album {
	if result, ok := m[name]; ok {
		return result
	}

	return *NewAlbum(name, AlbumManual)
}

func (m AlbumMap) Pointer(name string) *Album {
	if result, ok := m[name]; ok {
		return &result
	}

	return NewAlbum(name, AlbumManual)
}

var AlbumFixtures = AlbumMap{
	"christmas2030": {
		ID:               1000000,
		AlbumUID:         "as6sg6bxpogaaba7",
		AlbumSlug:        "christmas-2030",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Christmas 2030",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "Wonderful Christmas",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "oldest",
		AlbumTemplate:    "",
		AlbumCountry:     UnknownID,
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"holiday-2030": {
		ID:               1000001,
		AlbumUID:         "as6sg6bxpogaaba8",
		AlbumSlug:        "holiday-2030",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Holiday 2030",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "Wonderful Christmas Holiday",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "newest",
		AlbumTemplate:    "",
		AlbumCountry:     UnknownID,
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    true,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"berlin-2019": {
		ID:               1000002,
		AlbumUID:         "as6sg6bxpogaaba9",
		AlbumSlug:        "berlin-2019",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Berlin 2019",
		AlbumLocation:    "Berlin",
		AlbumCategory:    "City",
		AlbumCaption:     "",
		AlbumDescription: "We love Berlin 🌿!",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "oldest",
		AlbumTemplate:    "",
		AlbumCountry:     "",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"april-1990": {
		ID:               1000003,
		AlbumUID:         "as6sg6bipogaaba1",
		AlbumSlug:        "april-1990",
		AlbumPath:        "1990/04",
		AlbumType:        AlbumFolder,
		AlbumTitle:       "April 1990",
		AlbumLocation:    "",
		AlbumCategory:    "Friends",
		AlbumCaption:     "",
		AlbumDescription: "Spring is the time of year when many things change.",
		AlbumNotes:       "",
		AlbumFilter:      "path:\"1990/04\" public:true",
		AlbumOrder:       "added",
		AlbumTemplate:    "",
		AlbumCountry:     "ca",
		AlbumYear:        1990,
		AlbumMonth:       4,
		AlbumDay:         11,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"import": {
		ID:               1000004,
		AlbumUID:         "as6sg6bitoga0004",
		AlbumSlug:        "import",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Import Album",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "ca",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"emptyMoment": {
		ID:               1000005,
		AlbumUID:         "as6sg6bitoga0005",
		AlbumSlug:        "empty-moment",
		AlbumPath:        "",
		AlbumType:        AlbumMoment,
		AlbumTitle:       "Empty Moment",
		AlbumLocation:    "Favorite Park",
		AlbumCategory:    "Fun",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "public:true country:at year:2016",
		AlbumOrder:       "oldest",
		AlbumTemplate:    "",
		AlbumCountry:     "at",
		AlbumYear:        2016,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"2016-04": {
		ID:               1000006,
		AlbumUID:         "as6sg6bipogaabj8",
		AlbumSlug:        "2016-04",
		AlbumPath:        "2016/04",
		AlbumType:        AlbumFolder,
		AlbumTitle:       "April 2016",
		AlbumLocation:    "",
		AlbumCategory:    "Fun",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "path:\"2016/04\" public:true",
		AlbumOrder:       "added",
		AlbumTemplate:    "",
		AlbumCountry:     UnknownID,
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"september-2021": {
		ID:               1000007,
		AlbumUID:         "as6sg6bipogaabj9",
		AlbumSlug:        "september-2021",
		AlbumPath:        "",
		AlbumType:        AlbumMonth,
		AlbumTitle:       "September 2021",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "public:true year:2021 month:9",
		AlbumOrder:       "newest",
		AlbumTemplate:    "",
		AlbumCountry:     UnknownID,
		AlbumYear:        2021,
		AlbumMonth:       9,
		AlbumDay:         0,
		AlbumFavorite:    false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"california-usa": {
		ID:               1000008,
		AlbumUID:         "as6sg6bipogaab11",
		AlbumSlug:        "california-usa",
		AlbumPath:        "",
		AlbumType:        AlbumState,
		AlbumTitle:       "California / USA",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "public:true country:us state:California",
		AlbumOrder:       "newest",
		AlbumTemplate:    "",
		AlbumCountry:     "us",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"california-duplicate-1": {
		ID:               1000009,
		AlbumUID:         "as6sg6bipotaab12",
		AlbumSlug:        "california-usa",
		AlbumPath:        "",
		AlbumType:        AlbumState,
		AlbumTitle:       "California / USA",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "country:us state:California",
		AlbumOrder:       "newest",
		AlbumTemplate:    "",
		AlbumCountry:     "us",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"california-duplicate-2": {
		ID:               1000010,
		AlbumUID:         "as6sg6bipotaab19",
		AlbumSlug:        "california",
		AlbumPath:        "",
		AlbumType:        AlbumState,
		AlbumTitle:       "California",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "public:true country:us state:California",
		AlbumOrder:       "newest",
		AlbumTemplate:    "",
		AlbumCountry:     "us",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		CreatedAt:        time.Date(2019, 7, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"&ilikefood": {
		ID:               1000011,
		AlbumUID:         "as6sg6bipotaab18",
		AlbumSlug:        "&ilikefood",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "&IlikeFood",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"i-love-%-dog": {
		ID:               1000012,
		AlbumUID:         "as6sg6bipotaab20",
		AlbumSlug:        "i-love-%-dog",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "I love % dog",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"%gold": {
		ID:               1000013,
		AlbumUID:         "as6sg6bipotaab21",
		AlbumSlug:        "%gold",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "%gold",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"sale%": {
		ID:               1000014,
		AlbumUID:         "as6sg6bipotaab22",
		AlbumSlug:        "sale%",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "sale%",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"pets&dogs": {
		ID:               1000015,
		AlbumUID:         "as6sg6bipotaab23",
		AlbumSlug:        "pest&dogs",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Pets & Dogs",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"light&": {
		ID:               1000016,
		AlbumUID:         "as6sg6bipotaab24",
		AlbumSlug:        "light&",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Light&",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"'family": {
		ID:               1000017,
		AlbumUID:         "as6sg6bipotaab25",
		AlbumSlug:        "'family",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "'Family",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"father's-day": {
		ID:               1000018,
		AlbumUID:         "as6sg6bipotaab26",
		AlbumSlug:        "father's-day",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Father's Day",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"ice-cream'": {
		ID:               1000019,
		AlbumUID:         "as6sg6bipotaab27",
		AlbumSlug:        "ice-cream'",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Ice Cream'",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"*forrest": {
		ID:               1000020,
		AlbumUID:         "as6sg6bipotaab28",
		AlbumSlug:        "*forrest",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "*Forrest",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"my*kids": {
		ID:               1000021,
		AlbumUID:         "as6sg6bipotaab29",
		AlbumSlug:        "my*kids",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "My*Kids",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"yoga***": {
		ID:               1000022,
		AlbumUID:         "as6sg6bipotaab30",
		AlbumSlug:        "yoga***",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Yoga***",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"|banana": {
		ID:               1000023,
		AlbumUID:         "as6sg6bipotaab31",
		AlbumSlug:        "|banana",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "|Banana",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"blue|": {
		ID:               1000024,
		AlbumUID:         "as6sg6bipotaab33",
		AlbumSlug:        "blue|",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Blue|",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"345-shirt": {
		ID:               1000025,
		AlbumUID:         "as6sg6bipotaab34",
		AlbumSlug:        "345-shirt",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "345 Shirt",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"color-555-blue": {
		ID:               1000026,
		AlbumUID:         "as6sg6bipotaab35",
		AlbumSlug:        "color-555-blue",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Color555 Blue",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"route-66": {
		ID:               1000027,
		AlbumUID:         "as6sg6bipotaab36",
		AlbumSlug:        "route-66",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Route 66",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"red|green": {
		ID:               1000028,
		AlbumUID:         "as6sg6bipotaab32",
		AlbumSlug:        "red|green",
		AlbumPath:        "",
		AlbumType:        AlbumManual,
		AlbumTitle:       "Red|Green",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumFilter:      "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"germany": {
		ID:               1000029,
		AlbumUID:         "as6sg6bipotaah64",
		AlbumSlug:        "germany",
		AlbumPath:        "",
		AlbumType:        AlbumMoment,
		AlbumTitle:       "Germany",
		AlbumFilter:      "public:true country:de",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
	"mexico": {
		ID:               1000030,
		AlbumUID:         "as6sg6bipotaaj10",
		AlbumSlug:        "mexico",
		AlbumPath:        "",
		AlbumType:        AlbumMoment,
		AlbumTitle:       "Nature",
		AlbumFilter:      "public:true country:mx",
		AlbumLocation:    "",
		AlbumCategory:    "",
		AlbumCaption:     "",
		AlbumDescription: "",
		AlbumNotes:       "",
		AlbumOrder:       "name",
		AlbumTemplate:    "",
		AlbumCountry:     "zz",
		AlbumYear:        0,
		AlbumMonth:       0,
		AlbumDay:         0,
		AlbumFavorite:    false,
		AlbumPrivate:     false,
		CreatedAt:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:        gorm.DeletedAt{},
	},
}

// CreateAlbumFixtures inserts known entities into the database for testing.
func CreateAlbumFixtures() {
	for _, entity := range AlbumFixtures {
		Db().Save(&entity)
	}
}
