package roomdb

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID                string `gorm:"primaryKey;not null;unique"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	ShortID           string `gorm:"unique"`
	PeerServerAddr    string
	SelectedContentId string
	Owner             User
	Members           []User
}

func CreateRoom(owner User) (*Room, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	peerServer, err := SelectAvailablePeerServer()
	if err != nil {
		return nil, fmt.Errorf("cannot create room, peer server not available: %w", err)
	}
	r := &Room{
		ID:             GenerateUUID(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		ShortID:        GenerateShortID(),
		PeerServerAddr: peerServer.Address(),
		Owner:          owner,
	}

	result := dbHandle.Create(r)
	if result.Error != nil {
		return nil, result.Error
	}

	return r, nil
}

func GetRoom(id string, short bool) (*Room, error) {
	if dbHandle == nil {
		return nil, ErrNoDBConnection
	}

	var r Room
	if short {
		result := dbHandle.Preload("Owner").Preload("Members").First(&r, "short_id = ?", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := dbHandle.First(&r, "id = ?", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return &r, nil
}

func UpdateRoom(r Room) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	// updates modified fields, fails if row doesn't exist in DB yet.
	var original Room

	result := dbHandle.First(&original, "id = ?", r.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrAccountNotFound
	} else if result.Error != nil {
		return result.Error
	}

	dbHandle.Model(&original).Updates(r)

	return nil
}

func DeleteRoom(id string, short bool) error {
	if dbHandle == nil {
		return ErrNoDBConnection
	}

	if short {
		result := dbHandle.Delete(&Room{}, "short_id = ?", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrAccountNotFound
		}
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := dbHandle.Delete(&Room{}, "id = ?", id)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrAccountNotFound
		}
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
