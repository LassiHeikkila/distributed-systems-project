package roomdb

import (
	"errors"
	"fmt"
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
	}

	return r, errors.New("unimplemented")
}
