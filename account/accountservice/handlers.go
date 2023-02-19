package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

const (
	maxSignupBodyLen = 1024
)

func SignupHandler(w http.ResponseWriter, req *http.Request) {
	// POST, read body
	// Response will be JSON containing possible error message or id of newly created account

	// Sign up form should contain:
	// - username (unique)
	// - password

	// Don't want or need email, phone, personal info

	// Form should be sent as a JSON:
	type signupForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b, _ := io.ReadAll(io.LimitReader(req.Body, maxSignupBodyLen))
	defer req.Body.Close()

	var f signupForm
	if err := json.Unmarshal(b, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(failedToParseRequest))
		return
	}

	taken, err := accountdb.UsernameTaken(f.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(databaseError))
		return
	}
	if taken {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(usernameTaken))
		return
	}

	pwHash := HashPassword(f.Password)

	id, err := accountdb.CreateAccount(
		accountdb.Account{
			Username:     f.Username,
			PasswordHash: pwHash,
		},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf(accountCreationSuccessFmt, id)))
}

func AccountInfoHandler(w http.ResponseWriter, req *http.Request) {
	// GET, account id in URL variables
	id := mux.Vars(req)["id"]

	a, err := accountdb.GetAccount(id)
	if errors.Is(err, accountdb.ErrAccountNotFound) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(accountWithIdNotFound))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(databaseError))
		return
	}

	e := json.NewEncoder(w)
	_ = e.Encode(a)
}

func AccountInfoUpdateHandler(w http.ResponseWriter, req *http.Request) {
	// PUT
}

func AccountInfoDeleteHandler(w http.ResponseWriter, req *http.Request) {
	// DELETE
}

func TokenAuthenticateHandler(w http.ResponseWriter, req *http.Request) {
	// GET, check token from header

	// Auth middleware handles the check,
	// so if we're here, we're authenticated
	_, _ = w.Write([]byte(genericOK))
}

func TokenDeauthenticateHandler(w http.ResponseWriter, req *http.Request) {
	// POST, read body
	// body contains tokens to deauthenticate, one per line
	// also accept wildcard "*" to deauthenticate every token for the user
}
