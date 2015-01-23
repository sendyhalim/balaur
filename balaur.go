package balaur

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/fatih/color"
	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
)

var neededConfigs = []string{"app", "route", "middleware"}

type ControllerMethod func(web.C, *http.Request) (string, int)

func NewApp(dir string, configs map[string]string) *App {
	app := &App{
		appConfig:   map[string]Config{},
		Controllers: map[string]interface{}{},
		Middlewares: map[string]interface{}{},
		Mux:         web.New(),
	}

	// load all given configs
	checkNeededConfigs(configs)
	for k, p := range configs {
		app.appConfig[k] = NewConfig(fmt.Sprintf("%s/%s", dir, p))
	}

	app.Name = app.Config("app").Get("name", false)

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
func checkNeededConfigs(configs map[string]string) {
	for _, k := range neededConfigs {
		_, ok := configs[k]
		if !ok {
			glog.Fatalf("Please set config for %s")
		}
	}
}

type App struct {
	Name                    string
	Priority                int
	Controllers             map[string]interface{}
	Middlewares             map[string]interface{}
	ControllerMethodWrapper func(method ControllerMethod) web.HandlerFunc
	Mux                     *web.Mux

	appConfig map[string]Config
}

func (a *App) Boot() {
	color.Green("Booting app %s(priority %d)\n", a.Name, a.Priority)
	a.registerRoutes()
	a.registerMiddlewares()
	AddApp(a)
}

func (a *App) Config(section string) Config {
	conf, ok := a.appConfig[section]

	if !ok {
		glog.Fatalf("Config (%s) is not set", section)
	}

	return conf
}

func (a *App) registerRoutes() {
	conf := a.Config("route")

	for _, r := range conf.GetChildren("route", false) {
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
	conf := a.Config("middleware")

	for _, m := range conf.GetChildren("middleware", false) {
		key := m.Get("key", true)

		a.Mux.Use(a.Middlewares[key])
	}
}
