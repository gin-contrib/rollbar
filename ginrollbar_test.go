package ginrollbar

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogPanicsToRollbar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	})
	testError := &gin.Error{
		Err:  errors.New("test error"),
		Type: gin.ErrorTypePublic,
	}
	testError.SetMeta("some data")
	router.Use(PanicLogs(false, ""))
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(testError)
		_ = c.Error(testError)
		panic("occurs panic")
	})

	RollbarCritical = func(interfaces ...interface{}) {
		if err, ok := interfaces[0].(error); ok {
			assert.Equal(t, "occurs panic", err.Error())
		} else {
			t.Error("interfaces[0] should be error")
		}
		if request, ok := interfaces[1].(*http.Request); ok {
			assert.Equal(t, "/", request.RequestURI)
		} else {
			t.Error("interfaces[1] should be *http.Request")
		}
		if level, ok := interfaces[2].(int); ok {
			assert.Equal(t, 3, level)
		} else {
			t.Error("interfaces[2] should be int")
		}
		if metaData, ok := interfaces[3].(map[string]interface{}); ok {
			fmt.Printf("%+v", metaData)
			endpoint, _ := metaData["endpoint"].(string)
			assert.Equal(t, "/", endpoint)
		} else {
			t.Error("interfaces[3] should be map[string]interface{}")
		}
	}

	errorCalls := 0
	RollbarError = func(interfaces ...interface{}) {
		errorCalls++
		if err, ok := interfaces[0].(error); ok {
			assert.Equal(t, testError.Err.Error(), err.Error())
		} else {
			t.Error("interfaces[0] should be error")
		}
		if request, ok := interfaces[1].(*http.Request); ok {
			assert.Equal(t, "/", request.RequestURI)
		} else {
			t.Error("interfaces[1] should be *http.Request")
		}
		if metaData, ok := interfaces[2].(map[string]interface{}); ok {
			endpoint, _ := metaData["endpoint"].(string)
			assert.Equal(t, "/", endpoint)
			meta, _ := metaData["meta"].(string)
			assert.Equal(t, "some data", meta)
		} else {
			t.Error("interfaces[2] should be map[string]interface{}")
		}
	}

	w := performRequest("GET", "/", router)

	assert.Equal(t, 500, w.Code, "http status code")
	assert.Equal(t, 2, errorCalls, "Calls to RollbarError")
}

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
