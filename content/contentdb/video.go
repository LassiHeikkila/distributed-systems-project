package contentdb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/glebarez/sqlite"

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
	ContentID string `json:"contentID" gorm:"primaryKey;unique;not null;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Files []VideoFile `gorm:"foreignKey:content_id"`

	Name            string `json:"name" gorm:"not null"`
	License         string `json:"license"`
	Attribution     string `json:"attribution"`
	DurationSeconds int    `json:"durationSeconds"`
	Category        string `json:"category"`
}

type VideoFile struct {
	FileID    string `json:"fileID" gorm:"primaryKey;unique;not null;<-:create"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ContentID     string    `json:"contentID"`
	Uploaded      time.Time `json:"uploaded" gorm:"<-:create"`
	Encoding      string    `json:"encoding"`
	Resolution    string    `json:"resolution"`
	FileSizeBytes uint      `json:"fileSizeBytes"`
	Hash          string    `json:"hash"`
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

	if err := dbHandle.AutoMigrate(&Video{}); err != nil {
		return fmt.Errorf("error automigrating Video table: %w", err)
	}
	if err := dbHandle.AutoMigrate(&VideoFile{}); err != nil {
		return fmt.Errorf("error automigrating VideoFile table: %w", err)
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

	// if v.ContentID is not set, set it to a freshly generated UUID
	if v.ContentID == "" {
		v.ContentID = GenerateUUID()
	}

	result := dbHandle.Create(&v)
	if result.Error != nil {
		return "", result.Error
	}

	return v.ContentID, nil
}

func AddVideoFile(vf VideoFile) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	err := dbHandle.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vf).Error; err != nil {
			return err
		}
		var v Video
		if err := tx.First(&v).Where("content_id = ?", vf.ContentID).Association("Files").Append(&vf); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to append video file to video: %w", err)
	}

	return nil
}

func GetVideo(id string) (*Video, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var v Video
	result := dbHandle.Preload("Files").First(&v, "content_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrVideoNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &v, nil
}

func GetVideoFile(id string) (*VideoFile, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var vf VideoFile
	result := dbHandle.First(&vf, "file_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrVideoNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &vf, nil
}

func UpdateVideo(v Video) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	// updates modified fields, fails if row doesn't exist in DB yet.
	var original Video

	result := dbHandle.First(&original, "content_id = ?", v.ContentID)
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

	result := dbHandle.Delete(&Video{}, "content_id = ?", id)
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
	result := dbHandle.Preload("Files").Where(clause, args...).Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}
