package main

const (
	failedToParseRequest  = `{"ok":false,"msg":"failed to parse request"}`
	databaseError         = `{"ok":false,"msg":"database error"}`
	usernameTaken         = `{"ok":false,"msg":"username taken"}`
	accountWithIdNotFound = `{"ok":false,"msg":"account with given id not found"}`
	unauthorizedError     = `{"ok":false,"msg":"you must provide authentication bearer token with Authorization header"}`
	forbiddenError        = `{"ok":false,"msg":"access forbidden"}`

	genericOK                 = `{"ok":true}`
	accountCreationSuccessFmt = `{"ok":true,"msg":"account created successfully","id":"%s"}`
)
