package roomdb

import (
	"errors"
	"time"
)

type Room struct {
	ID                string `gorm:"primaryKey;not null;unique"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
	ShortID           string `gorm:"unique"`
	PeerServerAddr    string
	SelectedContentId string
	Members           []User
}

func CreateRoom() (*Room, error) {
	r := &Room{
		ID:             GenerateUUID(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		ShortID:        GenerateShortID(),
		PeerServerAddr: SelectAvailablePeerServer().GetAddress(),
	}

	return r, errors.New("unimplemented")
}
