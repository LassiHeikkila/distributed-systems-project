package accountdb

import (
	"errors"

	"gorm.io/gorm"
)

// keep it simple, just one big table with all the accounts
// keep in mind that GORM somehow translates UpperCamelCase to upper_camel_case

type Account struct {
	gorm.Model

	UserID       string `json:"userID" gorm:"unique;not null;<-:create"`
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	IsAdmin      bool   `json:"-"`
}

// CRUD

func CreateAccount(a Account) (string, error) {
	if dbHandle == nil {
		return "", ErrNoDBConnection
	}

	// if a.UserID is not set, set it to a freshly generated UUID
	if a.UserID == "" {
		a.UserID = GenerateUUID()
	}

	result := dbHandle.Create(&a)
	if result.Error != nil {
		return "", result.Error
	}

	return a.UserID, nil
}

func GetAccount(id string) (*Account, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var a Account
	result := dbHandle.First(&a, "user_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &a, nil
}

func UpdateAccount(a Account) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	// updates modified fields, fails if row doesn't exist in DB yet.
	var original Account

	result := dbHandle.First(&original, "user_id = ?", a.UserID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else if result.Error != nil {
		return result.Error
	}

	dbHandle.Model(&original).Updates(a)

	return nil
}

func DeleteAccount(id string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	result := dbHandle.Delete(&Account{}, "user_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UsernameTaken(n string) (bool, error) {
	if dbHandle == nil {
		return false, ErrNoDBConnection
	}

	result := dbHandle.First(&Account{}, "username = ?", n)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	return false, nil
}

func UserIsAdmin(id string) (bool, error) {
	a, err := GetAccount(id)
	if err != nil {
		return false, err
	}

	return a.IsAdmin, nil
}
