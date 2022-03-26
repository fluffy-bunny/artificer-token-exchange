package session

import (
	"github.com/labstack/echo/v4"

	"echo-starter/internal/session"
)

// EnsureDevelopmentSession is a middleware that ensures that the session is
// wiped out when the app restarts
func EnsureDevelopmentSession(appInstanceID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.GetSession(c)
			appInstanceValue, ok := sess.Values["_appInstanceID"]
			if !ok {
				sess.Values["_appInstanceID"] = appInstanceID
				sess.Save(c.Request(), c.Response())
			} else {
				if appInstanceValue != appInstanceID {
					sess.Values = make(map[interface{}]interface{}) // wipe out the session
					sess.Values["_appInstanceID"] = appInstanceID
					sess.Save(c.Request(), c.Response())
				}
			}
			return next(c)
		}
	}
}
