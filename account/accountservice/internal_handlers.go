package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

func ServiceHandlerValidateToken(w http.ResponseWriter, req *http.Request) {
	authToken := mux.Vars(req)["token"]

	userID, err := accountdb.AuthenticateToken(authToken)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf(genericErrorFmt, "token could not be validated")))
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf(
		genericOKWithKVStringFmtFmt,
		"userID",
		userID,
	)))
}

func ServiceHandlerAccountLookup(w http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["id"]

	a, err := accountdb.GetAccount(userID)
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
