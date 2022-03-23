package session

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) *sessions.Session {
	ss, err := session.Get("session", c)
	if err != nil {
		panic(err)
	}
	return ss
}
