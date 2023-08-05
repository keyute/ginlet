# Ginlet

Ginlet is a wrapper around Gin that allows declarative routing and middleware, inspired by Cobra.

## Install

```shell
go get -u github.com/keyute/ginlet
```

## Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keyute/ginlet"
)

func main() {
	rg := ginlet.RestRouterGroup{
		GetRoute: ginlet.Route{
			Handler: func(c *gin.Context) {
				c.String(200, "Hello, world!")
			},
		},
	}
	engine := ginlet.Engine{
		RouterGroups: []ginlet.BaseRoute{rg},
	}
	r, _ := engine.New() // returns *gin.Engine
	_ = r.Run()
}

```

## Contributing

Ginlet is a highly opinionated wrapper that is primarily written for my projects.
However, I am open to making it as a general purpose wrapper for Gin.
If you have any suggestions or improvements, feel free to open an issue or a pull request.