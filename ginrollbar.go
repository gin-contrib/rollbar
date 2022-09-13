package ginrollbar

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
)

// Recovery middleware for rollbar error monitoring
func Recovery(onlyCrashes, printStack bool, requestIdCtxKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				if printStack {
					debug.PrintStack()
				}

				rollbar.Critical(errors.New(fmt.Sprint(rval)), c.Request, getCallers(3), map[string]interface{}{
					"endpoint": c.Request.RequestURI,
				})

				c.AbortWithStatus(http.StatusInternalServerError)
			}

			if !onlyCrashes {
				extraData := make(map[string]interface{})
				extraData["endpoint"] = c.Request.RequestURI
				if requestIdCtxKey != "" {
					extraData["request_id"] = c.Writer.Header().Get(requestIdCtxKey)
				}
				for _, item := range c.Errors {
					extraData["meta"] = fmt.Sprint(item.Meta)
					rollbar.Error(item.Err, c.Request, extraData)
				}
			}
		}()

		c.Next()
	}
}

func getCallers(skip int) (pc []uintptr) {
	pc = make([]uintptr, 1000)
	i := runtime.Callers(skip+1, pc)
	return pc[0:i]
}
