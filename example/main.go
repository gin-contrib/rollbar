package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/neiybor/ginrollbar/v2"

	"github.com/rollbar/rollbar-go"
)

func main() {
	rollbar.SetToken("MY_TOKEN")
	// roll.SetEnvironment("production") // defaults to "development"

	r := gin.Default()
	r.Use(ginrollbar.PanicLogs(false, ""))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
