package roomdb

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	RoomID    string
	PeerId    string
	Name      string
}

func CreateUser(id string, name string) (*User, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	u := User{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		RoomID:    "",
		PeerId:    GenerateUUID(),
		Name:      name,
	}

	return &u, nil
}

func GetUser(id string) (*User, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var u User
	result := dbHandle.First(&u, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &u, nil
}

func DeleteUser(id string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	result := dbHandle.Delete(&User{}, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *User) JoinRoom(id string) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	u.RoomID = id

	err := dbHandle.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", u.ID).Update("room_id", id).Error; err != nil {
			return err
		}
		if err := tx.Model(&Room{}).Where("id = ?", id).Association("Members").Append(u); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *User) LeaveRoom() error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	roomID := u.RoomID
	if roomID == "" {
		return nil
	}

	err := dbHandle.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", u.ID).Update("room_id", "").Error; err != nil {
			return err
		}
		if err := tx.Model(&Room{}).Where("id = ?", roomID).Association("Members").Delete(u); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
