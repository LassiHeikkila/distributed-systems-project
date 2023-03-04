package main

import (
	"encoding/json"
	"net/http"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
	"github.com/LassiHeikkila/flmnchll/helpers/httputils"
	"github.com/LassiHeikkila/flmnchll/room/roomdb"
	"github.com/gorilla/mux"
)

// Needed API endpoints:
// - creating a new room
//  - does not need to take any arguments
// - user joining a room
//  - get user id from auth token
//  - add user to room
//  - return room id
// - user get room details
// 	- return details if user is member in room
//   - check auth token
// - user get own current room (if any)
// - delete room
//  - only admin or room creator can do it
// - user leaving a room
//  - get user id from auth token
//  - add user to room
// - room having content selected

func CreateRoomHandler(w http.ResponseWriter, req *http.Request) {}

func DeleteRoomHandler(w http.ResponseWriter, req *http.Request) {}

func SetRoomSelectedContentHandler(w http.ResponseWriter, req *http.Request) {}

func JoinRoomHandler(w http.ResponseWriter, req *http.Request) {}

func LeaveRoomHandler(w http.ResponseWriter, req *http.Request) {}

func GetRoomDetailsHandler(w http.ResponseWriter, req *http.Request) {
	// check user is authenticated
	userID, err := accountdb.AuthenticateToken(httputils.GetAuthToken(req))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(forbiddenError))
		return
	}

	// get room id from query parameters
	// account id in URL variables
	roomID := mux.Vars(req)["id"]

	// check that user is in the room
	r, err := roomdb.GetRoom(roomID, true)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(roomWithIdNotFound))
		return
	}

	userIsMember := func() bool {
		for _, u := range r.Members {
			if u.ID == userID {
				return true
			}
		}
		return false
	}

	userIsOwner := func() bool {
		return r.Owner.ID == userID
	}

	if !userIsMember() || !userIsOwner() {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(forbiddenError))
		return
	}

	// return room details
	e := json.NewEncoder(w)

	_ = e.Encode(roomToJSON(r))
}

func GetCurrentRoomHandler(w http.ResponseWriter, req *http.Request) {
	// check user is authenticated
	userID, err := accountdb.AuthenticateToken(httputils.GetAuthToken(req))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(forbiddenError))
		return
	}

	// get user object
	u, err := roomdb.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(userWithIdNotFound))
		return
	}

	// get the room based on room id
	roomID := u.RoomID
	if roomID == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(roomWithIdNotFound))
		return
	}

	r, err := roomdb.GetRoom(roomID, false)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(roomWithIdNotFound))
		return
	}

	// return room json
	e := json.NewEncoder(w)

	_ = e.Encode(roomToJSON(r))
}
