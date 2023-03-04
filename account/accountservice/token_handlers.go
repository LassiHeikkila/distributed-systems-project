package main

import (
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
