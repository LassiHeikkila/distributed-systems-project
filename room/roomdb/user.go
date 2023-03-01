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
	DeletedAt time.Time
	RoomID    string
	Name      string
}

func CreateUser() (*User, error) {
	return nil, errors.New("unimplemented")
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
