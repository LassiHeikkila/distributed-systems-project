package main

import (
	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

func tokenMatchesUserId(token string, id string) bool {
	tokenId, err := accountdb.AuthenticateToken(token)
	if err != nil {
		return false
	}

	return tokenId == id
}
