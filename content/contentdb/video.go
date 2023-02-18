package contentdb

import (
	"errors"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbHandle *gorm.DB

	ErrNoDBConnection = errors.New("no connection to database")
	ErrVideoNotFound  = errors.New("video not found")
)

// keep it simple, just one big table with all the videos
// keep in mind that GORM somehow translates UpperCamelCase to upper_camel_case

type Video struct {
	gorm.Model

	FileID         string  `json:"fileID" gorm:"unique;not null;<-:create"`
	OriginalFileID *string `json:"originalFileID"`
	FileHash       string  `json:"fileHash"`

	Name            string    `json:"name" gorm:"not null"`
	License         string    `json:"license"`
	Attribution     string    `json:"attribution"`
	Uploaded        time.Time `json:"uploaded" gorm:"<-:create"`
	Encoding        string    `json:"encoding"`
	DurationSeconds int       `json:"durationSeconds"`
	Resolution      string    `json:"resolution"`
	FileSizeBytes   uint      `json:"fileSizeBytes"`
	Category        string    `json:"category"`
}

func Connect(path string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	dbHandle = db

	return nil
}

func Init() error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	err := dbHandle.AutoMigrate(&Video{})
	if err != nil {
		return err
	}

	return nil
}

func Disconnect() error {
	if dbHandle == nil {
		return nil
	}

	dbHandle = nil
	return nil
}

// CRUD

func AddVideo(v Video) (string, error) {
	if dbHandle == nil {
		return "", ErrNoDBConnection
	}

	// if v.FileID is not set, set it to a freshly generated UUID
	if v.FileID == "" {
		v.FileID = GenerateUUID()
	}

	result := dbHandle.Create(&v)
	if result.Error != nil {
		return "", result.Error
	}

	return v.FileID, nil
}

func GetVideo(id string) (*Video, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var v Video
	result := dbHandle.First(&v, "file_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrVideoNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &v, nil
}

func UpdateVideo(v Video) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	// updates modified fields, fails if row doesn't exist in DB yet.
	var original Video

	result := dbHandle.First(&original, "file_id = ?", v.FileID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrVideoNotFound
	} else if result.Error != nil {
		return result.Error
	}

	dbHandle.Model(&original).Updates(v)

	return nil
}

func DeleteVideo(id string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	result := dbHandle.Delete(&Video{}, "file_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrVideoNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Search

func SearchVideos(searchOptions ...SearchOption) ([]Video, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	getClauses := func(opts []SearchOption) []string {
		clauses := make([]string, len(opts))
		for i := range opts {
			clauses[i] = opts[i]().Clause
		}
		return clauses
	}
	getArguments := func(opts []SearchOption) []interface{} {
		args := make([]interface{}, len(searchOptions))
		for i := range searchOptions {
			args[i] = searchOptions[i]().Argument
		}
		return args
	}
	clause := strings.Join(getClauses(searchOptions), " AND ")
	args := getArguments(searchOptions)

	results := make([]Video, 0)
	result := dbHandle.Where(clause, args...).Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}
