package main

const (
	failedToParseRequest = `{"ok":false,"msg":"failed to parse request"}`
	databaseError        = `{"ok":false,"msg":"database error"}`
	usernameTaken        = `{"ok":false,"msg":"username taken"}`
	roomWithIdNotFound   = `{"ok":false,"msg":"room with given id not found"}`
	userWithIdNotFound   = `{"ok":false,"msg":"user with given id not found"}`
	//	unauthorizedError     = `{"ok":false,"msg":"you must provide authentication bearer token with Authorization header"}`
	forbiddenError        = `{"ok":false,"msg":"access forbidden"}`
	unimplementedError    = `{"ok":false,"msg":"unimplemented feature"}`
	peerServerUnavailable = `{"ok":false","msg":"PeerJS server unavailable at this time, please try again later"}`

	genericOK       = `{"ok":true}`
	genericErrorFmt = `{"ok":false,"msg":"%s"}`
)

func NewOKResponseWithDetails(d map[string]any) map[string]any {
	if d == nil {
		d = make(map[string]any)
	}
	d["ok"] = true
	return d
}
