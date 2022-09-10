package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/neiybor/ginrollbar"

	"github.com/rollbar/rollbar-go"
)

func main() {
	rollbar.SetToken("MY_TOKEN")
	// roll.SetEnvironment("production") // defaults to "development"

	r := gin.Default()
	r.Use(ginrollbar.Recovery(true, false))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
