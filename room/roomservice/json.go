package main

import "github.com/LassiHeikkila/flmnchll/room/roomdb"

type RoomJSON struct {
	ShortID         string      `json:"shortID"`
	FullID          string      `json:"fullID"`
	PeerServerAddr  string      `json:"peerServerAddr"`
	SelectedContent ContentJSON `json:"selectedContent"`
	Members         []UserJSON  `json:"users"`
}

type ContentJSON struct {
	ID string `json:"id"`
}

type UserJSON struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	PeerID string `json:"peerID"`
}

func roomToJSON(r *roomdb.Room) *RoomJSON {
	if r == nil {
		return nil
	}

	members := make([]UserJSON, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(
			members,
			UserJSON{
				ID:     m.ID,
				Name:   m.Name,
				PeerID: m.PeerId,
			},
		)
	}

	return &RoomJSON{
		ShortID:        r.ShortID,
		FullID:         r.ID,
		PeerServerAddr: r.PeerServerAddr,
		SelectedContent: ContentJSON{
			ID: r.SelectedContentId,
		},
		Members: members,
	}
}
