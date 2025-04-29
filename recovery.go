package rollbar

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rollbar/rollbar-go"
)

// Recovery middleware for rollbar error monitoring
func Recovery(onlyCrashes bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rollbar.Critical(
					errors.New(fmt.Sprint(rval)),
					3,
					map[string]string{"endpoint": c.Request.RequestURI},
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}

			if !onlyCrashes {
				for _, item := range c.Errors {
					rollbar.Error(item.Err, map[string]string{
						"meta":     fmt.Sprint(item.Meta),
						"endpoint": c.Request.RequestURI,
					})
				}
			}
		}()

		c.Next()
	}
}
