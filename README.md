#Balaur
An app container for [goji](https://github.com/zenazn/goji) app

[![Build Status](https://travis-ci.org/sendyhalim/balaur.svg)](https://travis-ci.org/sendyhalim/balaur)


#Usage
For example I have a web app and it has 2 parts, user and blog, with `balaur` it would be like this:

**Structure**
```
├── main.go
├── blog
│   ├── app.go
│   ├── controller.go
│   ├── middleware.go
│   ├── config.toml
│   ├── middlewares.toml
│   └── routes.toml
└── user
    ├── app.go
    ├── controller.go
    ├── middleware.go
    ├── config.toml
    ├── middlewares.toml
    └── routes.toml
```
user `app.go` should handle creation of user app and `boot` the app (the same applies to `blog`)

here's what the basic app should look like
```go
package user
import (
	"github.com/sendyhalim/balaur"
)

func init() {
	var ctrl = &UserController{}
	// 1st param is the directory path relative to your root app
	// 2nd param is the configs (app, route, and middleware must be supplied)
	// 3rd param is RouteRegistrar interface
	// 4th param is MiddlewareRegistrar interface 
	// by default basic registrar will be created for you if you pass nil but of course you can create your own registrar 
	// as long as it conforms the interface
	var app = balaur.NewApp("user", map[string]string{
		"app":        "config.toml",
		"route":      "routes.toml",
		"middleware": "middlewares.toml",
	}, nil, nil)
    
    // register the controller, it will be mapped by route registrar
	app.Controllers["user"] = ctrl
	// register middlewares, the middleware will be registered for the app's parent route (set in routes.toml)
	// app.Controller["authMiddleware"] = authMiddleware
	// then in middlewares.toml:
	// [[middleware]]
	// key = "authMiddleware" 
	app.Boot()
}
```

**App config (config.toml)**
```
name = "User App"
```

**routes config (routes.toml)**
```
parent = "/*"

[[route]]
path       = "/users"
verb       = "GET"
controller = "user"
method     = "Index"

[[route]]
path       = "/users/:id"
verb       = "GET"
controller = "user"
method     = "Get"
```

**User `controller.go`**
```go
package user

import (
	"net/http"
	"github.com/sendyhalim/balaur/context"
)

type UserController struct {
}

func (a *UserController) Index(c context.ContextInterface, r *http.Request) (string, int) {
	return "User index!", http.StatusOK
}

func (a *UserController) Get(c context.ContextInterface, r *http.Request) (string, int) {
	return "User id: " + c.GetParam("id"), http.StatusOK
}

```

then in `main.go`
```go
package main

import (
	_ "myproject/user"
	_ "myproject/blog"
	"github.com/sendyhalim/balaur"
	"github.com/zenazn/goji"
)

func main() {
    // balaur.GetApp(appName)
	goji.Handle("/*", balaur.GetApp("User App").Mux)
	goji.Handle("/*", balaur.GetApp("Blog App").Mux)
	goji.Serve()
}
```

# TODOs
* Add readme example how to use different registrars and controllers
* API documentation

# Notes
for now `balaur` only use `toml` as config, but it will support more config types in the future (such as json)
