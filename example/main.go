package main

import (
	"github.com/gin-contrib/rollbar"
	"github.com/gin-gonic/gin"

	roll "github.com/rollbar/rollbar-go"
)

func main() {
	roll.SetToken("MY_TOKEN")
	// roll.SetEnvironment("production") // defaults to "development"

	r := gin.Default()
	r.Use(rollbar.Recovery(true))

	r.Run(":8080")
}
