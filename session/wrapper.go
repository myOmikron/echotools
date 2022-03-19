package session

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/echotools/middleware"
	"github.com/myOmikron/echotools/utility"
	"reflect"
)

//LoginRequired Helper function to mark endpoint as login only. Requires Session as a middleware.
// Returns ErrSessionMisconfigured if middleware.SessionContext is not present in Context.
func LoginRequired[T middleware.SessionContext](f func(cc T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if SessionContext is available
		field := reflect.ValueOf(c).Elem().FieldByName("SessionContext")
		if field == (reflect.Value{}) {
			c.Logger().Error(middleware.ErrSessionMisconfigured)
			return c.JSON(500, utility.JsonResponse{Error: "Internal server error"})
		}

		// Check if user is authenticated
		t := c.(T)
		if !t.IsAuthenticated() {
			return c.JSON(403, utility.JsonResponse{Error: "Unauthenticated"})
		}

		return f(t)
	}
}
