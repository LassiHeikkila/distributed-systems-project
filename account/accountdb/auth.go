package accountdb

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthToken struct {
	gorm.Model

	Value     string `gorm:"unique;not null`
	UserID    string `gorm:"not null"`
	ExpiresAt time.Time
}

// Note: AuthToken schema is nowhere near optimal,
// Value field should probably be the primary key.

func CreateNewTokenForUserID(id string, expiresAt time.Time) (string, error) {
	if dbHandle == nil {
		return "", ErrNoDBConnection
	}
	if id == "" {
		return "", ErrNoIDGiven
	}

	token := GenerateUUID()
	at := AuthToken{
		Value:     token,
		UserID:    id,
		ExpiresAt: expiresAt,
	}

	result := dbHandle.Create(&at)
	if result.Error != nil {
		return "", result.Error
	}

	return token, nil
}

func RevokeToken(t string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	result := dbHandle.Delete(&AuthToken{}, "value = ?", t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RevokeTokensForUser(id string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	_ = dbHandle.Where("user_id = ?", id).Delete(&AuthToken{})
	return nil
}

// AuthenticateToken checks if token `t` exists and is valid.
// It returns the user id corresponding to the token, or an error.
func AuthenticateToken(t string) (string, error) {
	if dbHandle == nil {
		return "", ErrNoDBConnection
	}

	a := AuthToken{}
	result := dbHandle.First(&a, "value = ?", t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", ErrTokenDoesNotExist
	}
	if result.Error != nil {
		return "", result.Error
	}
	if a.ExpiresAt.Before(time.Now()) {
		return "", ErrTokenExpired
	}

	return a.UserID, nil
}
