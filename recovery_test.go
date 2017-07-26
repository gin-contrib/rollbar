package rollbar

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(Recovery(false))

	router.GET("/", func(c *gin.Context) {
		baseError := errors.New("test error")
		err := &gin.Error{
			Err:  baseError,
			Type: gin.ErrorTypePublic,
		}
		err.SetMeta("some data")
		c.Error(err)

		panic("occurs panic")
	})

	w := performRequest("GET", "/", router)
	assert.Equal(t, 500, w.Code)
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
