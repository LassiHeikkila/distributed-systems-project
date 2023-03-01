package main

import "net/http"

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
// - delete room
//  - only admin or room creator can do it
// - user leaving a room
//  - get user id from auth token
//  - add user to room
// - room having content selected

func CreateRoomHandler(w http.ResponseWriter, req *http.Request)      {}
func DeleteRoom(w http.ResponseWriter, req *http.Request)             {}
func SetRoomSelectedContent(w http.ResponseWriter, req *http.Request) {}

func JoinRoom(w http.ResponseWriter, req *http.Request)  {}
func LeaveRoom(w http.ResponseWriter, req *http.Request) {}

func GetRoomDetails(w http.ResponseWriter, req *http.Request) {}
