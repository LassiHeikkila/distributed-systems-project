package roomdb

import (
	"errors"
	"time"
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

func (u *User) JoinRoom(id string) {
	u.RoomID = id
}
