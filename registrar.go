package balaur

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/zenazn/goji/web"
)

type ControllerMethod func(web.C, *http.Request) (string, int)

type RouteRegistrar interface {
	RegisterRoutes(*web.Mux, []Config, map[string]interface{})
}

type MiddlewareRegistrar interface {
	RegisterMiddlewares(*web.Mux, []Config, map[string]interface{})
}

type BasicRouteRegistrar struct {
	ControllerMethodWrapper func(ControllerMethod) web.HandlerFunc
}

func (r *BasicRouteRegistrar) RegisterRoutes(mux *web.Mux, routes []Config, controllers map[string]interface{}) {
	for _, config := range routes {
		controller := config.Get("controller", true)
		method := config.Get("method", true)
		verb := config.Get("verb", true)
		path := config.Get("path", true)
		methodInterface := reflect.ValueOf(controllers[controller]).MethodByName(method).Interface()
		methodValue := methodInterface.(func(web.C, *http.Request) (string, int))
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

type BasicMiddlewareRegistrar struct {
}

func (r *BasicMiddlewareRegistrar) RegisterMiddlewares(mux *web.Mux, middlewareConfigs []Config, middlewares map[string]interface{}) {
	for _, m := range middlewareConfigs {
		key := m.Get("key", true)

		mux.Use(middlewares[key])
	}
}

func NewBasicMiddlewareRegistrar() *BasicMiddlewareRegistrar {
	return &BasicMiddlewareRegistrar{}
}

func NewBasicRouteRegistrar() *BasicRouteRegistrar {
	r := &BasicRouteRegistrar{}

	// basic controller method wrapper
	r.ControllerMethodWrapper = func(method ControllerMethod) web.HandlerFunc {
		fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
			c.Env["Content-Type"] = "text/html"

			response, code := method(c, r)

			switch code {
			case http.StatusOK:
				if _, exists := c.Env["Content-Type"]; exists {
					w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
				}
				fmt.Fprint(w, response)
			case http.StatusSeeOther:
				http.Redirect(w, r, response, code)
			}
		}
		return web.HandlerFunc(fn)
	}

	return r
}
