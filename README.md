#Balaur
An app container for [goji](https://github.com/zenazn/goji) app

[![Build Status](https://travis-ci.org/sendyhalim/balaur.svg)](https://travis-ci.org/sendyhalim/balaur)
[![GoDoc](https://godoc.org/github.com/sendyhalim/balaur/web?status.svg)](https://godoc.org/github.com/sendyhalim/balaur)
[![views](https://sourcegraph.com/api/repos/github.com/sendyhalim/balaur/.counters/views.png)](https://sourcegraph.com/github.com/sendyhalim/balaur)
[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/sendyhalim/balaur/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

#Table of contents
- [Balaur](#balaur)
- [Usage](#usage)
- [Custom](#custom)
    - [RouteRegistrar](#route-registrar)
    - [MiddlewareRegistrar](#middleware-registrar)
- [TODOs](#todos)
- [Notes](#notes)


#Usage
**NOTE** 
You don't have to follow this structure, you are free to define your own structure, this structure is just my convention

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

#Custom

##Route Registrar
You can inject your own route registrar by creating custom registrar and inject it(3rd param) when you use `balaur.NewApp()`.
```go
type RouteRegistrar interface {
	// 1st param is goji Mux
	// 2nd param is routes config (balaur already convert routes config as Config interface)
	// 3rd param is the mapping of app controllers
	RegisterRoutes(*web.Mux, []Config, map[string]interface{})
}
```

By creating your own custom route registrar, you can gain full flexibility on controller methods.
Here's `balaur.BasicRouteRegistrar` code ([full](https://github.com/sendyhalim/balaur/blob/master/registrar.go))
```go
type ControllerMethod func(context.ContextInterface, *http.Request) (string, int)

type BasicRouteRegistrar struct {
	// ControllerMethodWrapper is assigned when NewBasicRouteRegistrar() is called
	ControllerMethodWrapper func(ControllerMethod) web.HandlerFunc
}

func (r *BasicRouteRegistrar) RegisterRoutes(mux *web.Mux, routes []Config, controllers map[string]interface{}) {
	for _, config := range routes {
		controller := config.Get("controller", true)
		method := config.Get("method", true)
		verb := config.Get("verb", true)
		path := config.Get("path", true)
		methodInterface := reflect.ValueOf(controllers[controller]).MethodByName(method).Interface()
		methodValue := methodInterface.(func(context.ContextInterface, *http.Request) (string, int))
		handler := r.ControllerMethodWrapper(methodValue)

		switch verb {
		case "GET":
			mux.Get(path, handler)
		case "POST":
			mux.Post(path, handler)
		case "PUT":
			mux.Put(path, handler)
		case "DELETE":
			mux.Delete(path, handler)
		}
	}
}
```
As you can see by creating custom registrar, you can also modify controller methods type. If you only want to use different context (by default goji context is used), just do this
```go
var routeRegistrar *BasicRouteRegistrar = balaur.NewBasicRouteRegistrar()
routeRegistrar.ControllerMethodWrapper = func(method balaur.ControllerMethod) web.HandlerFunc {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		// as long as it conforms context.ContextInterface, then you can inject any 
		// custom context to  BasicRouteRegistrar's  ControllerMethod
		response, code := method(&GorillaContext, r)
		// do something else..
	}
	return web.HandlerFunc(fn)	
}

// create app and inject routeRegistrar manually
var app = balaur.NewApp("user", map[string]string{
		"app":        "config.toml",
		"route":      "routes.toml",
		"middleware": "middlewares.toml",
}, routeRegistrar, nil)
```

##Middleware Registrar
You can inject your own middleware registrar by creating custom registrar and inject it (4th param) when you use `balaur.NewApp()`.
```go
type MiddlewareRegistrar interface {
	// 1st param is goji Mux
	// 2nd param is middlewares config (balaur already convert middlewares config as Config interface)
	// 3rd param is the mapping of app controllers
	RegisterMiddlewares(*web.Mux, []Config, map[string]interface{})
}
```

`BasicMiddlewareRegistrar` is really simple, it just maps the middleware config based on its index

```go
type BasicMiddlewareRegistrar struct {}

func (r *BasicMiddlewareRegistrar) RegisterMiddlewares(mux *web.Mux, middlewareConfigs []Config, middlewares map[string]interface{}) {
	for _, m := range middlewareConfigs {
		key := m.Get("key", true)

		mux.Use(middlewares[key])
	}
}
```

# TODOs
* Better API documentation
* Scaffolding tools (automatic app generation e.g `balaur mkapp blog` then create the basic templates inside the given dir)

# Notes
for now `balaur` only use `toml` as config, but it will support more config types in the future (such as json)

# Support
I will always improve and support this project as I'm using `balaur` for most of my projects.
