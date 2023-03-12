package accountdb

import (
	"errors"
	"fmt"

	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbHandle *gorm.DB

	ErrNoDBConnection    = errors.New("no connection to database")
	ErrAccountNotFound   = errors.New("account not found")
	ErrNoIDGiven         = errors.New("empty id parameter")
	ErrTokenDoesNotExist = errors.New("token does not exist")
	ErrTokenExpired      = errors.New("token has expired")
)

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

	if err := dbHandle.AutoMigrate(&Account{}); err != nil {
		return fmt.Errorf("error automigrating Account: %w", err)
	}
	if err := dbHandle.AutoMigrate(&AuthToken{}); err != nil {
		return fmt.Errorf("error automigrating AuthToken: %w", err)
	}

	// TODO: scan auth tokens and remove any that have expired

	return nil
}

func Disconnect() error {
	if dbHandle == nil {
		return nil
	}

	dbHandle = nil
	return nil
}
