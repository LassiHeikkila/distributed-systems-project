package main

import (
	"bufio"
	"errors"
	"io"

	"github.com/LassiHeikkila/flmnchll/account/accountdb"
)

func tokenMatchesUserId(token string, id string) bool {
	tokenId, err := accountdb.AuthenticateToken(token)
	if err != nil {
		return false
	}

	if tokenId == id {
		return true
	}

	// check if token belongs to admin, in that case allow it
	if isAdmin, _ := accountdb.UserIsAdmin(tokenId); isAdmin {
		return true
	}

	return false
}

func invalidateTokens(r io.Reader, userId string) error {
	s := bufio.NewScanner(r)

	for s.Scan() {
		switch s.Text() {
		case "*":
			// invalidate all tokens for user with id "userId"
			// stop processing the rest
			return accountdb.RevokeTokensForUser(userId)
		default:
			t := s.Text()
			// invalidate the token in s.Text() *if* it matches with userId
			if !tokenMatchesUserId(t, userId) {
				return errors.New(unauthorizedError)
			}
			err := accountdb.RevokeToken(t)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
