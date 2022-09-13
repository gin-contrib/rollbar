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
func Recovery(onlyCrashes, printStack bool) gin.HandlerFunc {
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
				for _, item := range c.Errors {
					rollbar.Error(item.Err, c.Request, map[string]interface{}{
						"meta":     fmt.Sprint(item.Meta),
						"endpoint": c.Request.RequestURI,
					})
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

// Helper functions to convert pkg/errors stack traces to be compatible with rollbar
// Cribbed from pkg/errors

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func FramePc(f errors.Frame) uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func FrameFile(f errors.Frame) string {
	fn := runtime.FuncForPC(FramePc(f))
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(FramePc(f))
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func FrameLine(f errors.Frame) int {
	fn := runtime.FuncForPC(FramePc(f))
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(FramePc(f))
	return line
}

// name returns the name of this function, if known.
func FrameName(f errors.Frame) string {
	fn := runtime.FuncForPC(FramePc(f))
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}
