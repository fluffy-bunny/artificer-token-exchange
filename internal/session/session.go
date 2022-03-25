package session

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	sessionName = "_session"
)

func GetSession(c echo.Context) *sessions.Session {
	ss, err := session.Get(sessionName, c)
	if err != nil {
		panic(err)
	}
	return ss
}
func TerminateSession(c echo.Context) {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		panic(err)
	}
	sess.Values = make(map[interface{}]interface{}) // wipe out the session
	sess.Save(c.Request(), c.Response())

}
