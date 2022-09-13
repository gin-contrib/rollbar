# GinRollbar

[![Run Tests](https://github.com/neiybor/ginrollbar/actions/workflows/go.yml/badge.svg)](https://github.com/neiybor/ginrollbar/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/neiybor/ginrollbar/branch/master/graph/badge.svg)](https://codecov.io/gh/neiybor/ginrollbar)
[![Go Report Card](https://goreportcard.com/badge/github.com/neiybor/ginrollbar)](https://goreportcard.com/report/github.com/neiybor/ginrollbar)
[![GoDoc](https://godoc.org/github.com/neiybor/ginrollbar?status.svg)](https://godoc.org/github.com/neiybor/ginrollbar)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Middleware to integrate with [rollbar](https://rollbar.com/) error monitoring. It uses [rollbar-go](https://github.com/rollbar/rollbar-go) SDK that reports errors and logs messages.

## Usage

Download and install it:

```sh
go get github.com/neiybor/ginrollbar
```

Import it in your code:

```go
import "github.com/neiybor/ginrollbar"
```

## Example

```go
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
  r.Use(ginrollbar.Recovery(true, false, ""))

  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```
