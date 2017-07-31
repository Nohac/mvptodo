package models

import (
	"net/http"

	"github.com/gocraft/dbr"
	"golang.org/x/net/context"
)

const sessionsTable = "sessions"

type SessionType struct {
	SeesionCookie string
	Request       *http.Request
	Response      http.ResponseWriter
	DB            *dbr.Session
}

func Session(ctx context.Context) *SessionType {
	val := ctx.Value("session")
	if val == nil {
		return nil
	}

	session, ok := val.(*SessionType)
	if ok == false {
		return nil
	}

	return session
}
