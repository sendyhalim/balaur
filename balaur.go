package balaur

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/fatih/color"
	"github.com/zenazn/goji/web"
)

type ControllerMethod func(web.C, *http.Request) (string, int)

func NewApp(dir string) *App {
	app := &App{
		AppConfig:        NewConfig(dir + "/config.toml"),
		RouteConfig:      NewConfig(dir + "/routes.toml"),
		MiddlewareConfig: NewConfig(dir + "/middlewares.toml"),
		Controllers:      map[string]interface{}{},
		Middlewares:      map[string]interface{}{},
		Mux:              web.New(),
	}

	app.Name = app.AppConfig.Get("name", false)

	// default controllerMethodWrapper
	app.ControllerMethodWrapper = func(method ControllerMethod) web.HandlerFunc {
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

	return app
}

type App struct {
	Name                    string
	Priority                int
	AppConfig               Config
	RouteConfig             Config
	MiddlewareConfig        Config
	Controllers             map[string]interface{}
	Middlewares             map[string]interface{}
	ControllerMethodWrapper func(method ControllerMethod) web.HandlerFunc
	Mux                     *web.Mux
}

func (a *App) Boot() {
	color.Green("Booting app %s(priority %d)\n", a.Name, a.Priority)
	a.registerRoutes()
	a.registerMiddlewares()
	AddApp(a)
}

func (a *App) registerRoutes() {
	for _, r := range a.RouteConfig.GetChildren("route", false) {
		controller := r.Get("controller", true)
		method := r.Get("method", true)
		verb := r.Get("verb", true)
		path := r.Get("path", true)
		methodInterface := reflect.ValueOf(a.Controllers[controller]).MethodByName(method).Interface()
		handler := a.ControllerMethodWrapper(methodInterface.(func(web.C, *http.Request) (string, int)))

		switch verb {
		case "GET":
			a.Mux.Get(path, handler)
		case "POST":
			a.Mux.Post(path, handler)
		case "PUT":
			a.Mux.Put(path, handler)
		case "DELETE":
			a.Mux.Delete(path, handler)
		}
	}
}

func (a *App) registerMiddlewares() {
	for _, m := range a.MiddlewareConfig.GetChildren("middleware", false) {
		key := m.Get("key", true)

		a.Mux.Use(a.Middlewares[key])
	}
}
