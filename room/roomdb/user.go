package roomdb

import (
	"errors"
	"log"
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

func GetUserByName(username string) (*User, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var u User
	result := dbHandle.First(&u, "name = ?", username)
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

	log.Printf("adding user %s to room %s\n", u.Name, id)

	u.RoomID = id

	err := dbHandle.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(u).Update("room_id", id).Error; err != nil {
			return err
		}
		if err := tx.First(&Room{}).Where("id = ?", id).Association("Members").Append(u); err != nil {
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

	log.Printf("removing user %s from room %s\n", u.Name, u.RoomID)

	roomID := u.RoomID
	if roomID == "" {
		return nil
	}
	r, _ := GetRoom(roomID, false)

	err := dbHandle.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(u).Update("room_id", "").Error; err != nil {
			log.Println("1")
			return err
		}
		if err := tx.Model(r).Association("Members").Delete(u); err != nil {
			log.Println("2")
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
